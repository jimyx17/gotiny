package gotiny

import (
	"errors"
	"reflect"
	"unsafe"
)

type Encoder struct {
	buf     []byte // out encode buffer
	off     int
	boolPos int  // Next bool pos (buf[boolPos])
	boolBit byte //N ext bool bit in buf boolpos

	engines []encEng
	length  int
}

func Marshal(is ...interface{}) (out []byte, err error) {

	e, err := NewEncoderWithPtr(is...)
	if err != nil {
		return nil, errors.New("could not marshal this")
	}
	return e.Encode(is...)
}

// Creates a new from ps (given that ps is a pointer)
func NewEncoderWithPtr(ps ...interface{}) (e *Encoder, err error) {

	defer func() {
		if r := recover(); r != nil {
			e = nil
			err = errors.New("could not build encoder")
		}
	}()

	l := len(ps)
	engines := make([]encEng, l)
	for i := 0; i < l; i++ {
		rt := reflect.TypeOf(ps[i])
		if rt.Kind() != reflect.Ptr {
			panic("must a pointer type!")
		}
		engines[i] = getEncEngine(rt.Elem())
	}
	return &Encoder{
		length:  l,
		engines: engines,
	}, nil
}

// Creates a new from is
func NewEncoder(is ...interface{}) (enc *Encoder, err error) {

	defer func() {
		if r := recover(); r != nil {
			enc = nil
			err = errors.New("could not build encoder")
		}
	}()

	l := len(is)
	engines := make([]encEng, l)
	for i := 0; i < l; i++ {
		engines[i] = getEncEngine(reflect.TypeOf(is[i]))
	}
	return &Encoder{
		length:  l,
		engines: engines,
	}, nil
}

func NewEncoderWithType(ts ...reflect.Type) (enc *Encoder, err error) {

	defer func() {
		if r := recover(); r != nil {
			enc = nil
			err = errors.New("could not build encoder")
		}
	}()

	l := len(ts)
	engines := make([]encEng, l)
	for i := 0; i < l; i++ {
		engines[i] = getEncEngine(ts[i])
	}
	return &Encoder{
		length:  l,
		engines: engines,
	}, nil
}

// Encoder object in bytes (input value must be a pointer)
func (e *Encoder) Encode(is ...interface{}) (o []byte, err error) {

	defer func() {
		if r := recover(); r != nil {
			o = nil
			err = errors.New("could not encode")
		}
	}()

	engines := e.engines
	for i := 0; i < len(engines) && i < len(is); i++ {
		engines[i](e, (*[2]unsafe.Pointer)(unsafe.Pointer(&is[i]))[1])
	}
	return e.reset(), nil
}

// Encoder object in bytes (input value must be a pointer of type unsafe.Pointer)
func (e *Encoder) EncodePtr(ps ...unsafe.Pointer) (o []byte, err error) {

	defer func() {
		if r := recover(); r != nil {
			o = nil
			err = errors.New("could not encode")
		}
	}()

	engines := e.engines
	for i := 0; i < len(engines) && i < len(ps); i++ {
		engines[i](e, ps[i])
	}
	return e.reset(), nil
}

// Encode value vs
func (e *Encoder) EncodeValue(vs ...reflect.Value) (o []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			o = nil
			err = errors.New("could not encode")
		}
	}()

	engines := e.engines
	for i := 0; i < len(engines) && i < len(vs); i++ {
		engines[i](e, getUnsafePointer(&vs[i]))
	}
	return e.reset(), nil
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
	return buf
}
