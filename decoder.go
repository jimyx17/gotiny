package gotiny

import (
	"errors"
	"reflect"
	"unsafe"
)

type Decoder struct {
	buf     []byte //buf
	index   int    //Next byte index
	boolPos byte   // Next bool pos (buf[boolPos])
	boolBit byte   // Next bool bit in buf boolpos

	engines []decEng // Decoders
	length  int      // n of Decoders
}

func Unmarshal(buf []byte, is ...interface{}) (ret int, err error) {
	d, err := NewDecoderWithPtr(is...)
	if err != nil {
		return 0, errors.New("could not unmarshal this")
	}
	return d.Decode(buf, is...)
}

func NewDecoderWithPtr(is ...interface{}) (dec *Decoder, err error) {
	defer func() {
		if r := recover(); r != nil {
			dec = nil
			err = errors.New("could not decode this")
		}
	}()

	l := len(is)
	engines := make([]decEng, l)
	for i := 0; i < l; i++ {
		rt := reflect.TypeOf(is[i])
		if rt.Kind() != reflect.Ptr {
			panic("must a pointer type!")
		}
		engines[i] = getDecEngine(rt.Elem())
	}
	return &Decoder{
		length:  l,
		engines: engines,
	}, nil
}

func NewDecoder(is ...interface{}) (dec *Decoder, err error) {

	defer func() {
		if r := recover(); r != nil {
			dec = nil
			err = errors.New("could not build decoder")
		}
	}()

	l := len(is)
	engines := make([]decEng, l)
	for i := 0; i < l; i++ {
		engines[i] = getDecEngine(reflect.TypeOf(is[i]))
	}
	return &Decoder{
		length:  l,
		engines: engines,
	}, nil
}

func NewDecoderWithType(ts ...reflect.Type) (dec *Decoder, err error) {
	defer func() {
		if r := recover(); r != nil {
			dec = nil
			err = errors.New("could not build decoder")
		}
	}()

	l := len(ts)
	des := make([]decEng, l)
	for i := 0; i < l; i++ {
		des[i] = getDecEngine(ts[i])
	}
	return &Decoder{
		length:  l,
		engines: des,
	}, nil
}

func (d *Decoder) reset() int {
	index := d.index
	d.index = 0
	d.boolPos = 0
	d.boolBit = 0
	return index
}

// is is pointer of variable
func (d *Decoder) Decode(buf []byte, is ...interface{}) (o int, err error) {

	defer func() {
		if r := recover(); r != nil {
			o = 0
			err = errors.New("could not decode")
		}
	}()

	d.buf = buf
	engines := d.engines
	for i := 0; i < len(engines) && i < len(is); i++ {
		engines[i](d, (*[2]unsafe.Pointer)(unsafe.Pointer(&is[i]))[1])
	}
	return d.reset(), nil
}

// ps is a unsafe.Pointer of the variable
func (d *Decoder) DecodePtr(buf []byte, ps ...unsafe.Pointer) (o int, err error) {

	defer func() {
		if r := recover(); r != nil {
			o = 0
			err = errors.New("could not decode")
		}
	}()

	d.buf = buf
	engines := d.engines
	for i := 0; i < len(engines) && i < len(ps); i++ {
		engines[i](d, ps[i])
	}
	return d.reset(), nil
}

func (d *Decoder) DecodeValue(buf []byte, vs ...reflect.Value) (o int, err error) {

	defer func() {
		if r := recover(); r != nil {
			o = 0
			err = errors.New("could not decode")
		}
	}()

	d.buf = buf
	engines := d.engines
	for i := 0; i < len(engines) && i < len(vs); i++ {
		engines[i](d, unsafe.Pointer(vs[i].UnsafeAddr()))
	}
	return d.reset(), nil
}
