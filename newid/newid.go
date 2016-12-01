package main

import (
    "fmt"
	"math/rand"
	"time"
)

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