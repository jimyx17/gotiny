package main

import (
	"fmt"
	"reflect"

	"github.com/jimyx17/gotiny"
)

func main() {
	src1, src2 := "hello", []byte(" world!")
	ret1, ret2 := "", []byte{3, 4, 5}
	d, _ := gotiny.Marshal(src1, &src2)
	gotiny.Unmarshal(d, &ret1, &ret2)
	fmt.Println(ret1 + string(ret2)) // print "hello world!"

	enc, _ := gotiny.NewEncoder(src1, src2)
	dec, _ := gotiny.NewDecoder(ret1, ret2)

	ret1, ret2 = "", []byte{3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 4, 5, 6, 7, 44, 7, 5, 6, 4, 7}
	o, _ := enc.EncodeValue(reflect.ValueOf(src1), reflect.ValueOf(src2))
	dec.DecodeValue(o,
		reflect.ValueOf(&ret1).Elem(), reflect.ValueOf(&ret2).Elem())
	fmt.Println(ret1 + string(ret2)) // print "hello world!"
}
