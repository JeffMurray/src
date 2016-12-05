package errs

import (
	"testing"
	"fmt"
	"time"
)

func TestClnErr(t *testing.T) {
	tm := time.Now()
	ce := NewClnErr("x9l8ew", "Some user bound message.")
	if ce.Token != "x9l8ew" {
		t.Error("ce.Token != x9l8ew")
	}
	if ce.Message != "Some user bound message." {
		t.Error("ce.Message != Some user bound message.")
	}
	dur := ce.Time.Sub(tm)
	if dur.Seconds() != 0 {
		t.Error("ce.Time.Seconds != 0")
	}
	ce.Traced("zen6h5","Rethrow 1").Traced("gfkdu3","Rethrow 2")
	fmt.Println(ce.Error())
	ce, _ = ClnErrFromBytes([]byte(ce.Error()))
	if ce.Token != "x9l8ew" {
		t.Error("ce.Token != x9l8ew")
	}
	if ce.Message != "Some user bound message." {
		t.Error("ce.Message != Some user bound message.")
	}
	se := ce.ToSysErr("This shit's fucked up", "e5x5ox", "up==down")
	if se.Token != "e5x5ox" {
		t.Error("se.Token != e5x5ox")
	}
	if se.Trace[0].Token != ce.Token {
		t.Error("se.Trace[0]Token != ce.Token")
	}
	for i, trc := range ce.Trace {
		if se.Trace[i+1].Token != trc.Token {
			t.Error(fmt.Sprintf("%s != %s",se.Trace[i+1].Token, trc.Token ))
		}
	}
	if TraceClnErrIfErr(nil,"","") != nil {
		t.Error("value should be nil")
	}
	if ce = TraceClnErrIfErr(ce,"twwk4k","yo yo yo"); ce== nil {
		t.Error("value should not be nil")
	}
	if ce.Trace[len(ce.Trace)-1].Token != "twwk4k" {
		t.Error("Token != twwk4k")
	}
}
func TestSysErr(t *testing.T) {
	tm := time.Now()
	sys_err := NewSysErr("Wow, that hurt", "zzkesi", "x")
	if sys_err.Text !=  "Wow, that hurt" {
		t.Error("Text mismatch")
	}
	if sys_err.Token !=  "zzkesi" {
		t.Error("Token mismatch")
	}
	if sys_err.Args !=  "x" {
		t.Error("Args mismatch")
	}		
	dur := sys_err.Time.Sub(tm)
	if dur.Seconds() != 0 {
		t.Error("Seconds != 0")
	}
	sys_err_copy,err := SysErrFromBytes(sys_err.ToBytes())
	if err != nil {
		t.Error(err.Error)
	}
	if sys_err_copy.Text != sys_err.Text {
		t.Error("sys_err_copy.Text != sys_err.Text")
	}
	if sys_err_copy.Token != sys_err.Token {
		t.Error("sys_err_copy.Token != sys_err.Token")
	}
	if sys_err_copy.Args != sys_err.Args {
		t.Error("sys_err_copy.Args != sys_err.Args")
	}
	traced := sys_err_copy.Traced("s2som7","abc")
	if traced.Trace[0].Token != "s2som7" {
		t.Error("traced.Trace[0].Token != s2som7")	
	}
	fmt.Println(string(traced.ToPrettyBytes()))
	serr := TraceSysErrIfErr(NewSysErr("text", "b7dade","a=b"),"timbbm", "")
	if serr == nil {
		t.Error("val == nil")
	} else {
		if serr.Trace[0].Token != "timbbm" {
			t.Error("serr.Trace[0].Token != timbbm")
		}
	}
	cerr := serr.ToClnErr("b2t8y8", "Although you are not getting your answer, remember the sun is still shining.")
	if cerr.Token != "b2t8y8" {
		t.Error("cerr.Token != b2t8y8")
	}
	if cerr.Trace[0].Token != serr.Token {
		t.Error("se.Trace[0]Token != ce.Token")
	}
	if cerr.Trace[0].Message != "Internal system error." {
		t.Error("Internal system error wording error, lets not give away the farm.")
	}	
}