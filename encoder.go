package gotiny

import (
	"errors"
	"reflect"
	"unsafe"
)

const MAXOBJREFS = 1024

type Encoder struct {
	buf     []byte // out encode buffer
	off     int
	boolPos int  // Next bool pos (buf[boolPos])
	boolBit byte //N ext bool bit in buf boolpos
	objPos uint16
	ptr    [MAXOBJREFS]uint64

	engines []encEng
	length  int
}

func Marshal(is ...interface{}) (out []byte, err error) {

	e, err := NewEncoderWithPtr(is...)
	if err != nil {
		return nil, errors.New("could not marshal this")
	}
	return e.Encode(is...), nil
}

// Creates a new from ps (given that ps is a pointer)
func NewEncoderWithPtr(ps ...interface{}) (e *Encoder, err error) {

	l := len(ps)
	engines := make([]encEng, l)
	for i := 0; i < l; i++ {
		rt := reflect.TypeOf(ps[i])
		if rt.Kind() != reflect.Ptr {
			return nil, errors.New("must a pointer type!")
		}
		engines[i] = getEncEngine(rt.Elem())
	}
	return &Encoder{
		length:  l,
		engines: engines,
	}, nil
}

// Creates a new from is
func NewEncoder(is ...interface{}) *Encoder {

	l := len(is)
	engines := make([]encEng, l)
	for i := 0; i < l; i++ {
		engines[i] = getEncEngine(reflect.TypeOf(is[i]))
	}
	return &Encoder{
		length:  l,
		engines: engines,
	}
}

func NewEncoderWithType(ts ...reflect.Type) *Encoder {
	l := len(ts)
	engines := make([]encEng, l)
	for i := 0; i < l; i++ {
		engines[i] = getEncEngine(ts[i])
	}
	return &Encoder{
		length:  l,
		engines: engines,
	}
}

// Encoder object in bytes (input value must be a pointer)
func (e *Encoder) Encode(is ...interface{}) []byte {
	engines := e.engines
	for i := 0; i < len(engines) && i < len(is); i++ {
		engines[i](e, (*[2]unsafe.Pointer)(unsafe.Pointer(&is[i]))[1])
		e.objPos = 0
	}
	return e.reset()
}

// Encoder object in bytes (input value must be a pointer of type unsafe.Pointer)
func (e *Encoder) EncodePtr(ps ...unsafe.Pointer) []byte {

	engines := e.engines
	for i := 0; i < len(engines) && i < len(ps); i++ {
		engines[i](e, ps[i])
		e.objPos = 0
	}
	return e.reset()
}

// Encode value vs
func (e *Encoder) EncodeValue(vs ...reflect.Value) []byte {

	engines := e.engines
	for i := 0; i < len(engines) && i < len(vs); i++ {
		engines[i](e, getUnsafePointer(&vs[i]))
		e.objPos = 0
	}
	return e.reset()
}

// Sets output buffer for encoder
func (e *Encoder) AppendTo(buf []byte) {
	e.off = len(buf)
	e.buf = buf
}

func (e *Encoder) reset() []byte {
	buf := e.buf
	e.buf = buf[:e.off]
	e.boolBit = 0
	e.boolPos = 0
	// e.ptr = bst.Node{}
	e.objPos = 0
	return buf
}
