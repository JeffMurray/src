package errs

import (
	"time"
	"encoding/json"
	"fmt"
)

//package errs formalizes the seperation of system information 
//and user information when handeling errors durng code execution.
//use NewClientError for expected client messages like "Last name is required."
//And NewSysError for system errors.  SysError.Error() reports a state free error
//location token.  Use SysError.String() to get the state info internally.
type SysError struct {
	Text string
	Token string
	FileName string
	Args string
	Time time.Time
	Trace []SourceInfo
}
type SourceInfo struct{
	Token string
	FileName string	
	Args string
}
type ClientError struct {
	Token string
	message string
	TraceInfo []ClientError
}
func (ce *ClientError) Error() string {
	rval := ce.message
	for i:= 0; i < len(ce.TraceInfo); i++ {
		rval = rval + fmt.Sprintf("-> %s : %s", ce.TraceInfo[i].Token, ce.TraceInfo[i].message, )
	}
	return rval
}
func IfSysError(err error, token, file_name, args string) *SysError {
	if err == nil {
		return nil
	}
	return NewSysError(err.Error(), token, file_name, args) 
}
func NewClientError(token string, message string) *ClientError {
	return &ClientError{token, message, nil}
}
//do a few mundan things that remotely might return errors and are not conditional on ech other
//then chek for one at the same time
func CheckErrors(errs_to_check ...*ClientError) *ClientError {
	for i := 0; i < len(errs_to_check); i++ {
		if errs_to_check[i] != nil {
			return errs_to_check[i]
		}
	}
	return nil
}
//use to trace errors up the stack
func (ce *ClientError)Traced(token string, message string) *ClientError {
	ce.TraceInfo = append(ce.TraceInfo, ClientError{token, message, nil})
	return ce
}
func (ce *ClientError) String() string {
	return ce.Token + " : " + ce.Error()
}
//an easy way to trace errors up the stack and locate them in code using token
func (se *SysError) Traced(token string, file_name string, args string ) *SysError {
	se.Trace = append(se.Trace,SourceInfo{token, file_name, args})
	return se
}
func TraceIfError(sys_error *SysError, token string, file_name string, args string) *SysError{
	if sys_error != nil {
		return sys_error.Traced(token, file_name, args)
	}
	return nil
}
func NewSysError(text, token, file_name, args string) *SysError {
	return &SysError{text, token, file_name, 
					args, time.Now(), make([]SourceInfo,0)}
}
//A state-free client code traceable error
func (se *SysError) ToClientError(token, messsage string) *ClientError {
	ce := NewClientError(token, messsage)
	ce.Traced(se.Token,"Internal system error.")
	for i := 0; i < len(se.Trace); i++ {
		ce.Traced(se.Trace[i].Token, "Internal system error.")
	}
	return ce
}
//the default behavior is not to give out state info
//but still be able to locate the line of code (grep is great)
func (se *SysError) Error() string {
	rval := "System error. Trace info: " + se.Token
	for i := 0; i < len(se.Trace); i++ {
		rval += fmt.Sprintf(", %s", se.Trace[i].Token)
	}
	return rval
}
//from a json encoded byte array
func SysErrorFromBytes(key_bytes []byte) (*SysError, error) {
	var sys_err SysError
	err := json.Unmarshal(key_bytes, &sys_err)
	return &sys_err, err
}
//to a packed json encoded byte array
func (se *SysError) ToBytes() []byte {
	bytes,_ := json.Marshal(*se)
	return bytes
}
//to a more readable json encoded byte array.
func (se *SysError) ToPrettyBytes() []byte {
	bytes, _ := json.MarshalIndent(*se, "", "    ")
	return bytes
}