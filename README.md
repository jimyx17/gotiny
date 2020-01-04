## <font color="#FF4500" >gotiny [Alpha]。</font>

# gotiny   [![Build status][travis-img]][travis-url] [![License][license-img]][license-url] [![GoDoc][doc-img]][doc-url] [![Go Report Card](https://goreportcard.com/badge/github.com/jimyx17/gotiny)](https://goreportcard.com/report/github.com/jimyx17/gotiny)
The base idea is to generate encoders/decoders in advance so the use of reflect lib is reduced to the minimum
## hello word 
    package main
    import (
   	    "fmt"
   	    "github.com/jimyx17/gotiny"
    )
    
    func main() {
   	    src1, src2 := "hello", []byte(" world!")
   	    ret1, ret2 := "", []byte{}
   	    gotiny.Unmarshal(gotiny.Marshal(&src1, &src2), &ret1, &ret2)
   	    fmt.Println(ret1 + string(ret2)) // print "hello world!"
    }

## Features
- It's fast
- No memory allocations
- Support all golang data types except func and chan
- Serialize ALL fields even those not exported (customizable via go tag)
- Only strictly the same encoder / decoder would marshal/unmarshal
- null type would be serialized
- Really small size of serialized data
## Cycle values won't work TODO 
	type a *a
	var b a
	b = &b

## install
```bash
$ go get -u github.com/jimyx17/gotiny
```

## Encoding protocol
### Bool
The bool type occupies one bit, the true value is encoded as 1, and the false value is encoded as 0. When the bool type is encountered for the first time, a byte is applied, and the value is programmed into the lowest bit. When it is encountered the second time, the byte is programmed into the second lowest bit. When the ninth encounter, the bool value is applied for another byte, programmed into the lowest bit. And so on.
### Int
- uint/int8 typed as a byte into the next byte of the string。
- uint16,uint32,uint64,uint,uintptr use [Varints](https://developers.google.com/protocol-buffers/docs/encoding#varints) encoding
- int16,int32,int64,int use ZigZag [Varints](https://developers.google.com/protocol-buffers/docs/encoding#varints) encoding

### Float
float32/float64 use [gob](https://golang.org/pkg/encoding/gob/) encoding
### Complex
- complex64 is casted to uint64 for encoding
- complex128 uses two float64 for encoding real / imaginary parts

### String
first encodes the length as uint64，then encodes the byte array itself
### Pointers
If nil, it ends with false bool encoded. Else, it encodes a true and then dereference the pointer and encodes it.
### Array & Slice
First convert length to uint64, then use each element own encoding method.
### Maps
Like in arrays, first, encode the length as uint64, then encode key with it's own encoder, then a value, and so on
### Struct
Encode all members of the struct (including non exported ones). The struct will be strictly reduced


### License
MIT

[travis-img]: https://travis-ci.org/jimyx17/gotiny.svg?branch=master
[travis-url]: https://travis-ci.org/jimyx17/gotiny
[license-img]: http://img.shields.io/badge/license-MIT-green.svg?style=flat-square
[license-url]: http://opensource.org/licenses/MIT
[doc-img]: http://img.shields.io/badge/GoDoc-reference-blue.svg?style=flat-square
[doc-url]: https://godoc.org/github.com/jimyx17/gotiny


### Jimyx17 fork

The idea will remain the same, the only changes that are going to be introduced are:

- Errors won't panic (this might imply performance penalties)
- Will try to find a solucion for cycling values
- Will try to translate chinese into english... without understanding a single word of chinese and english not being my mother tongue