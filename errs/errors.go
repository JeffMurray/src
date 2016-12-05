package errs

import (
	"time"
	"encoding/json"
	//"reflect"
	//"fmt"
)

//package errs formalizes the seperation of system information 
//and user information when handeling errors durng code execution.
//use NewClnErr for expected client messages like "Last name is required."
//And NewSysErr for system errors.  SysErr.Error() reports a state free error
//location token.
type ClnErr struct {
	Token string
	Message string
	Time time.Time
	Trace[]*ClnErr
}
type SysErr struct {
	Text string
	Token string
	Args string
	Time time.Time
	Trace []*SrcInfo
}
type SrcInfo struct{
	Token string
	Args string
}
func NewEOFErr() *ClnErr {
	return NewClnErr("t0t84g","EOF")
}
func (ce *ClnErr) IsEOF() bool {
	return ce.Token == "t0t84g"
}
func NewClnErr(token, message string) *ClnErr {
	return &ClnErr{token, message, time.Now(), []*ClnErr{}}
}
func (ce *ClnErr) Error() string {
	return string(ce.ToPrettyBytes())
}
func (ce *ClnErr) ToSysErr(text, token, args string) *SysErr {
	se := NewSysErr(text, token, args)
	se.Traced(ce.Token,ce.Message)
	for _, ti := range ce.Trace {
		se.Traced(ti.Token,ti.Message)
	}
	return se
}
//use to trace errors up the stack as errors get re-thrown
func (ce *ClnErr)Traced(token, message string) *ClnErr {
	ce.Trace = append(ce.Trace, NewClnErr(token, message))
	return ce
}
//helper func to add trace infor in the case it is an error
func TraceClnErrIfErr(err *ClnErr, token, message string) *ClnErr {
	if err != nil {
		return err.Traced(token, message)
	}
	return nil
}
//from a json encoded byte array
func ClnErrFromBytes(json_bytes []byte) (*ClnErr, error) {
	var ce ClnErr
	err := json.Unmarshal(json_bytes, &ce)
	return &ce, err
}
//to a packed json encoded byte array
func (ce *ClnErr) ToBytes() []byte {
	bytes,_ := json.Marshal(*ce)
	return bytes
}
//to a more readable json encoded byte array.
func (ce *ClnErr) ToPrettyBytes() []byte {
	bytes, _ := json.MarshalIndent(*ce, "", "    ")
	return bytes
}
//do a few mundan things that remotely might return errors and are not conditional on ech other
//then chek for one at the same time
func CheckErrors(errs_to_check ...*ClnErr) *ClnErr {
	for _, err := range errs_to_check {
		if err != nil {
			return err
		}
	}
	return nil
}
//an easy way to trace errors up the stack and locate them in code using token
func (se *SysErr) Traced(token string, args string ) *SysErr {
	se.Trace = append(se.Trace,&SrcInfo{token, args})
	return se
}
//helper func to add trace infor in the case it is an arror
func TraceSysErrIfErr(err *SysErr, token string, args string) *SysErr{
	if err != nil {
		return err.Traced(token,args)
	}
	return nil
}
func NewSysErr(text, token, args string) *SysErr {
	//args might be big, just get a snip
	return &SysErr{text, token,args, time.Now(), []*SrcInfo{}}
}
//A state-free client code traceable error
func (se *SysErr) ToClnErr(token, messsage string) *ClnErr {
	ce := NewClnErr(token, messsage)
	ce.Traced(se.Token,"Internal system error.")
	for  _, t := range se.Trace  {
		ce.Traced(t.Token, "Internal system error.")
	}
	return ce
}
//the default behavior is not to give out state info
//but still be able to locate the line of code (grep is great)
func (se *SysErr) Error() string {
	return se.ToClnErr("w0vfq8","Internal system error.").Error()
}
//get all the jucy stuff
func (se *SysErr) ErrorWithStateInfo() string {
	return string(se.ToPrettyBytes())
}
//from a json encoded byte array
func SysErrFromBytes(key_bytes []byte) (*SysErr, error) {
	var sys_err SysErr
	err := json.Unmarshal(key_bytes, &sys_err)
	return &sys_err, err
}
//to a packed json encoded byte array
func (se *SysErr) ToBytes() []byte {
	bytes,_ := json.Marshal(*se)
	return bytes
}
//to a more readable json encoded byte array.
func (se *SysErr) ToPrettyBytes() []byte {
	bytes, _ := json.MarshalIndent(*se, "", "    ")
	return bytes
}