package main

import (
    "fmt"
	"math/rand"
	"time"
)
//generates a 6 char unique id that is handy for linking errors to lines of code
//go install newid ---> will make newid.exe which spits out a 6 digit random string
//Handy for use in errs package. 
func newid() string {
	return make_key(6)
}
const lbytes = "abcdefghijklmnopqrstuvwxyz0123456789"
var berand = rand.New(rand.NewSource(time.Now().UnixNano()))
func make_key(n int) (rval string) {
	b := make([]byte, n)
    b[0]=lbytes[berand.Intn(len(lbytes)-10)] //always start with a letter
	for i:=1;i<n;i++ {
		b[i] = lbytes[berand.Intn(len(lbytes))]
	}
    return string(b)
}
func main() {
	fmt.Println("")
	fmt.Println(newid())
	fmt.Println("")
}