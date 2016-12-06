package safe_file_relay

import (
	"errs"
	"errors"
	"os"
	"bytes"
	"at_rest"
	"log"
	"io"
	"io/ioutil"
)

type SafeFileRelay struct {
	write_at_chan chan *write_req
	reader_chan chan *read_req
	closer_chan chan bool
	closed bool
}
type write_req struct {
	buf []byte
	pos int64
	rsp_chan chan write_rsp
}
type read_req struct {
	b []byte
	rsp_chan chan read_rsp
}
type read_rsp struct {
	num_read int
	read_err *errs.ClnErr
}
type write_rsp struct {
	num_write int
	write_err *errs.ClnErr
}
type closer_req struct {
	rsp_chan chan bool
}
type chunk struct {
	pos int64
	size int
	crypt bool
}
type chunk_central struct {
	chunk_map map[int64]*chunk
	chunk_cache map[int64][]byte
	spill_map map[int64]*chunk
	spill_file *os.File
	spill_pos int64
	spill_file_dir string
}
func NewChunkCentral(spill_file_dir string) *chunk_central {
	return &chunk_central{
		map[int64]*chunk{},
		map[int64][]byte{},
		map[int64]*chunk{},
		nil,
		0,
		spill_file_dir,
	}
}
func (cc *chunk_central) GetSpillFile() (*os.File, *errs.SysErr) {
	if cc.spill_file == nil {
		outf, outf_err := ioutil.TempFile(cc.spill_file_dir, "sfr")
		if outf_err != nil {
			return nil, errs.NewSysErr(outf_err.Error(), "ssib80", "")
		}
		cc.spill_file = outf		
	}
	return cc.spill_file, nil
}
//this is called once at the very end
func (cc *chunk_central) CleanupSpillFile(){
	if cc.spill_file != nil {
		fn := cc.spill_file.Name()
		cc.spill_file.Close()
		cc.spill_file=nil
		os.Remove(fn)
	}
}
//check out the test source for an example of using this relay
func NewSafeFileRelay(tmp_files_dir string) (*SafeFileRelay) {
	write_at_chan, reader_chan, closer_chan := start(tmp_files_dir)
	return &SafeFileRelay {
		write_at_chan,
		reader_chan,
		closer_chan,
		false,
	}	
}
func (sfr *SafeFileRelay) WriteAt( b []byte, pos int64 ) (int, error) {
	rsp_chan := make(chan write_rsp)
	sfr.write_at_chan <- &write_req{b[:], pos, rsp_chan}
	w_rsp := <- rsp_chan
	if w_rsp.write_err == nil {
		return w_rsp.num_write, nil
	}
	return w_rsp.num_write, errors.New(w_rsp.write_err.Error())
}
func (sfr *SafeFileRelay) Read(b []byte) (int, error) {
	rsp_chan := make(chan read_rsp)
	sfr.reader_chan <- &read_req{b[:],rsp_chan}
	r_rsp := <- rsp_chan
	if r_rsp.read_err != nil {
		if r_rsp.read_err.IsEOF() {
			return r_rsp.num_read, io.EOF
		}
		return r_rsp.num_read, errors.New(r_rsp.read_err.Error())
	}
	return r_rsp.num_read, nil
}
//this closes WriteAts.  The go func is still serving the Reader at this point, 
//so it will not exit until the reader hits io.EOF.
func (sfr *SafeFileRelay) CloseWriteAt() {
	if !sfr.closed {
		sfr.closed = true
		sfr.closer_chan <- true
	}
}
func start(tmp_files_dir string) (chan *write_req, chan *read_req, chan bool) {
	writes := make(chan *write_req)
	closer := make(chan bool)
	reader := make(chan *read_req)
	go func() {
		key := at_rest.NewEncryptionKey()
		var pending_read_req *read_req = nil
		read_pos := int64(0)
		closed := false
		cc := NewChunkCentral(tmp_files_dir)
		for {
			select {
				case write_at := <- writes:
					if closed {
						cc.CleanupSpillFile()
						log.Fatalf("Write called after close.")
					}
					wr := write_rsp{len(write_at.buf), nil}
					if len(cc.chunk_cache) <= 10 { 
						//this cache of 10 buffers greatly reduces encrypted overflow 
						//to disk if the WriteAter and Reader are operating concurently
						cc.chunk_map[write_at.pos]=&chunk{write_at.pos, wr.num_write, false}
						csh := make([]byte, wr.num_write)
						copy(csh,write_at.buf[:])
						cc.chunk_cache[write_at.pos] = csh
					} else {
						enc_chunk, enc_file_err := encryptChunk(cc,key,write_at.buf[:])
						if enc_file_err != nil {
							wr.write_err = errs.NewClnErr("bupnhg", "Internal error.")
							cc.CleanupSpillFile()
							log.Fatalf(enc_file_err.Error())
						}
						cc.chunk_map[write_at.pos] = &chunk{write_at.pos, wr.num_write, true}
						cc.spill_map[write_at.pos] = enc_chunk
						cc.spill_pos = cc.spill_pos + int64(enc_chunk.size)
					}
					write_at.rsp_chan <- wr
					//it is possible for the reading to get ahead of the writing, 
					//pending_read_req deals with that
					if pending_read_req != nil {
						read_size, chunk_is_null ,read_err := 
							readNextEtc(read_pos, pending_read_req, key, cc, 
								tmp_files_dir, closed)
						read_pos = read_pos + int64(read_size)
						if read_err != nil && read_err.IsEOF() {
							return
						}
						if !chunk_is_null {
							pending_read_req = nil
						}
					}
				case read_rq := <- reader:
					read_size, chunk_is_null, read_err := 
						readNextEtc(read_pos, read_rq, key, cc, 
							tmp_files_dir, closed)
					read_pos = read_pos + int64(read_size)
					if read_err != nil && read_err.IsEOF() {
						return
					}
					if chunk_is_null {
						pending_read_req = read_rq
					}
				case <- closer:
					//This closes writes. The whole thing wont exit until after writes 
					//are closed and the reader reaches EOF
					closed = true
					//println("WRITING CLOSED")
					if pending_read_req != nil {
						read_size, chunk_is_null ,read_err := 
							readNextEtc(read_pos, pending_read_req, key, cc, 
								tmp_files_dir, closed)
						read_pos = read_pos + int64(read_size)
						if read_err != nil && read_err.IsEOF() {
							return
						}
						if !chunk_is_null {
							pending_read_req = nil
						}
					}
			}
		}
	}()
	return writes, reader, closer
}
func readNextEtc(reader_pos int64, req *read_req, key *[32]byte, 
	cc *chunk_central, tmp_files_dir string, closed bool) (int,bool,*errs.ClnErr) {
	if closed && reader_pos == totalSize(cc) {
		cc.CleanupSpillFile()
		req.rsp_chan <- read_rsp{0,errs.NewEOFClnErr()}
		return 0, true, errs.NewEOFClnErr()
	}
	chunk := chunkAt(cc,reader_pos)
	if chunk != nil {
		rr := new(read_rsp)
		if !chunk.crypt {
			cache := cc.chunk_cache[chunk.pos]
			rr.num_read = readCache(chunk, cache[:], reader_pos, req)
			reader_pos = reader_pos + int64(rr.num_read)
			if reader_pos == (chunk.pos + int64(chunk.size)) {
				//println("deleting cache item")
				delete(cc.chunk_cache,chunk.pos)
			}						
		} else {
			//println("reading enc chunk")
			rr.num_read, rr.read_err = readCryptoChunk(chunk, reader_pos, req, key, cc)
			reader_pos = reader_pos + int64(rr.num_read)
			if reader_pos == (chunk.pos + int64(chunk.size)) {
				//println("deleting spill map item")
				delete(cc.spill_map,chunk.pos)
			}						
		}				
		if rr.read_err != nil {
			cc.CleanupSpillFile()
			log.Fatalf(rr.read_err.Error())
		}
		req.rsp_chan <- *rr
		return rr.num_read, false, nil
	}
	return 0, true, nil
}
func readCryptoChunk(chunk *chunk, reader_pos int64,  
		req *read_req, key *[32]byte, cc *chunk_central) (int, *errs.ClnErr) {
	enc_chunk := cc.spill_map[chunk.pos]
	enc_data := make([]byte,enc_chunk.size)
	_, enc_data_err := cc.spill_file.ReadAt(enc_data, enc_chunk.pos)
	if enc_data_err != nil {
		log.Fatalf("%v : %v", chunk.pos, enc_data_err)
	}
	data, data_err := at_rest.Decrypt(enc_data[:], key)
	if enc_data_err != nil {
		log.Fatalf("%s : %v", chunk.pos, data_err)
	}
	buf := bytes.NewBuffer(data[reader_pos-chunk.pos:])
	//writing directly and safely to the buffer of the 
	//original read request thanks to the synchronizing 
	//power of channels
	read_size, read_err := buf.Read(req.b[:]) //<------
	if read_err != nil {
		if read_err == io.EOF {
			return read_size, errs.NewEOFClnErr()
		}
		return read_size, errs.NewClnErr("x9xki3","internal read error")
	}
	return read_size, nil
}
func readCache(chunk *chunk, cache []byte, reader_pos int64, req *read_req) int {
	buf_pos := reader_pos - chunk.pos
	return copy(req.b[:], cache[buf_pos:])
}
func totalSize( cc *chunk_central) int64 {
	for next_pos := int64(0);; {
		chunk := cc.chunk_map[next_pos]
		if chunk == nil {
			return next_pos
		}
		next_pos = (chunk.pos + int64(chunk.size))
	}
}
func chunkAt( cc *chunk_central, pos int64 )*chunk {
	//chunk_map makes for interesting buffer location math
	//Note to future self: The interface cannot gurantee
	//each buffer size is the same.  This method is compatible 
	//with that.
	for next_pos := int64(0);; {
		chunk := cc.chunk_map[next_pos]
		if chunk == nil {
			return nil
		}
		if pos >= chunk.pos && pos <= (chunk.pos + int64(chunk.size) - 1) {
			return chunk
		}
		next_pos = (chunk.pos + int64(chunk.size))
	}
}
func encryptChunk( cc *chunk_central, key *[32]byte, chnk []byte ) (*chunk, *errs.SysErr) {
	enc_bytes, enc_bytes_err := at_rest.Encrypt(chnk[:], key)
	if enc_bytes_err != nil {
		return nil, errs.NewSysErr(enc_bytes_err.Error(), "odwpmc", "encrypting")
	}
	sf, sf_err := cc.GetSpillFile()
	if sf_err != nil {
		return nil, sf_err
	}
	write_size, write_err := sf.WriteAt(enc_bytes[:],cc.spill_pos)
	if write_err != nil {
		return nil, errs.NewSysErr(write_err.Error(), "wch4gp", "")
	} else if write_size != len(enc_bytes[:]) {
		return nil, errs.NewSysErr("unable to complete write","rhzedt","")
	}
	return &chunk{cc.spill_pos, write_size, true}, nil
}