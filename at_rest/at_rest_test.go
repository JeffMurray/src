package at_rest

import (
	"testing"
	"io/ioutil"
	"fmt"
	"time"
)
func Test_at_rest(t *testing.T) {


	dat := fmt.Sprintf("This is some test to encrypt @ %v",time.Now())
	ioutil.WriteFile("testing.txt", []byte(dat), 0644)

	dat2, err := EncryptedAtRest("testing.txt")
	
	if err != nil {
	
		t.Error(err.Error())
	}
	
	if string(dat) != string(dat2) {
		t.Error(fmt.Sprintf("%s != %s",string(dat), string(dat2)))
	}
	
	
}
