package errs

import (
	"testing"
	"fmt"
	"time"
)

func TestSysError(t *testing.T) {
	tm := time.Now()
	sys_err := NewSysError("Wow, that hurt", "zzkesi", "errors_test.go", "x")
	if sys_err.Text !=  "Wow, that hurt" {
		t.Error("Text mismatch")
	}
	if sys_err.Token !=  "zzkesi" {
		t.Error("Token mismatch")
	}
	if sys_err.FileName !=  "errors_test.go" {
		t.Error("FileName mismatch")
	}	
	if sys_err.Args !=  "x" {
		t.Error("Args mismatch")
	}		
	dur := sys_err.Time.Sub(tm)
	if dur.Seconds() != 0 {
		t.Error("Seconds != 0")
	}
	sys_err_copy,err := SysErrorFromBytes(sys_err.ToBytes())
	if err != nil {
		t.Error(err.Error)
	}
	if sys_err_copy.Text != sys_err.Text {
		t.Error("sys_err_copy.Text != sys_err.Text")
	}
	if sys_err_copy.Token != sys_err.Token {
		t.Error("sys_err_copy.Token != sys_err.Token")
	}
	if sys_err_copy.FileName != sys_err.FileName {
		t.Error("sys_err_copy.FileName != sys_err.FileName")
	}
	if sys_err_copy.Args != sys_err.Args {
		t.Error("sys_err_copy.Args != sys_err.Args")
	}
	traced := sys_err_copy.Traced("s2som7","errors_test.go","abc")
	if traced.Trace[0].Token != "s2som7" {
		t.Error("traced.Trace[0].Token != s2som7")	
	}
	fmt.Println(string(traced.ToPrettyBytes()))
}