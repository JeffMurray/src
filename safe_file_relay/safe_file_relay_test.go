package safe_file_relay

import (
	"testing"
	"io/ioutil"
	"os"
	"io"
	"fmt"
	"path"
	"time"
)
func Test_sfr(t *testing.T) {
	dir, dir_err := os.Getwd()
	if dir_err != nil {
		panic(dir_err.Error())
	}
	os.MkdirAll(path.Join(dir, "testing/tmp"), os.ModeDir)
	os.MkdirAll(path.Join(dir, "testing/in"), os.ModeDir)
	os.MkdirAll(path.Join(dir, "testing/out"), os.ModeDir)
	dat := fmt.Sprintf("Put other files in the testing/in dir to further tests @ %v",time.Now())
	ioutil.WriteFile(path.Join(path.Join(dir, "testing/in"), "testing.txt"), []byte(dat), 0644)
	//sfr := NewSafeFileRelay(path.Join(dir, "testing/tmp"))
	fis, _ := ioutil.ReadDir(path.Join(dir, "testing/in"))
	done := make(chan bool)
	num_f := 0
	for _, fi := range fis {
		if !fi.IsDir() {
			num_f++
			fmt.Println(fi.Name())
			sfr := NewSafeFileRelay(path.Join(dir, "testing/tmp"))
			go writeFile(fi.Name(), sfr, done, dir)
			go readFile(fi.Name(), sfr, done, dir)
		}
	}
	for i := 0; i < num_f; i++ {
		<- done
		<- done
	}
}

func writeFile(name string, sfr *SafeFileRelay, done chan bool, dir string) {
	_, file_name := path.Split(name)
	dest_dir := path.Join(path.Join(dir, "testing/out"),file_name)
	the_file, the_file_err := os.Create(dest_dir)
	if the_file_err != nil {
		//fmt.Println("writeFile:the_file_err %v",the_file_err)
		done <- true
		return
	}	
	buf := make([]byte, 50000)
	for i:=0;;i++{
		//fmt.Println("pre sfr.Read")
		num, read_err := sfr.Read(buf[:])
		//fmt.Println("post sfr.Read %s", string(buf[:num]))
		if read_err == nil {
			//fmt.Printf("writing %v %v\n",i,read_err)
			the_file.Write(buf[:num])
		} else {
			fmt.Println("the_file_err %v",read_err)
			done <- true
			break
		}
	}
}

func readFile(name string, sfr *SafeFileRelay, done chan bool, dir string) {
	the_file, the_file_err := os.Open(path.Join(path.Join(dir, "testing/in"),name))
	if the_file_err != nil {
		fmt.Println("the_file_err %v",the_file_err)
		done <- true
		return
	}
	buf := make([]byte, 50000)
	pos := int64(0)
	for {
		num, file_err := the_file.Read(buf)
		if file_err == nil {
			//fmt.Println("pre sfr.WriteAt")
			//fmt.Println(".")
			cnt, _ := sfr.WriteAt(buf[:num],pos)
			//fmt.Println("post sfr.WriteAt %s", string(buf[:cnt]))
			if cnt != num {
				fmt.Println("cnt != num")
				break
			}
			pos = pos + int64(cnt)
		} else if file_err == io.EOF {
			if num > 0 {
				cnt, _ := sfr.WriteAt(buf[:num],pos)
				if cnt != num {
					fmt.Println("cnt != num2")
					break
					
				}				
			}
			sfr.CloseWriteAt()
			break
		} else {
			panic(file_err.Error())
		}
	}
	done <- true
}