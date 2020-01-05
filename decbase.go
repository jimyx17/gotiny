package gotiny

import (
	"errors"
	"time"
	"unsafe"
)

func (d *Decoder) decBool(v *bool) error {

	if d.boolBit == 0 {

		if len(d.buf) <= d.index {
			return errors.New("error decoding bool")
		}

		d.boolBit = 1
		d.boolPos = d.buf[d.index]
		d.index++
	}
	*v = d.boolPos&d.boolBit != 0
	d.boolBit <<= 1
	return nil
}

func (d *Decoder) decUint64(v *uint64) error {

	if len(d.buf) <= d.index {
		return errors.New("error decoding uint64")
	}
	buf, i := d.buf, d.index
	x := uint64(buf[i])
	if x < 0x80 {
		d.index++
		*v = x
		return nil
	}

	if len(d.buf) <= d.index+1 {
		return errors.New("error decoding uint64")
	}
	x1 := buf[i+1]
	x += uint64(x1) << 7
	if x1 < 0x80 {
		d.index += 2
		*v = x - 1<<7
		return nil
	}

	if len(d.buf) <= d.index+2 {
		return errors.New("error decoding uint64")
	}
	x2 := buf[i+2]
	x += uint64(x2) << 14
	if x2 < 0x80 {
		d.index += 3
		*v = x - (1<<7 + 1<<14)
		return nil
	}

	if len(d.buf) <= d.index+3 {
		return errors.New("error decoding uint64")
	}
	x3 := buf[i+3]
	x += uint64(x3) << 21
	if x3 < 0x80 {
		d.index += 4
		*v = x - (1<<7 + 1<<14 + 1<<21)
		return nil
	}

	if len(d.buf) <= d.index+4 {
		return errors.New("error decoding uint64")
	}
	x4 := buf[i+4]
	x += uint64(x4) << 28
	if x4 < 0x80 {
		d.index += 5
		*v = x - (1<<7 + 1<<14 + 1<<21 + 1<<28)
		return nil
	}

	if len(d.buf) <= d.index+5 {
		return errors.New("error decoding uint64")
	}
	x5 := buf[i+5]
	x += uint64(x5) << 35
	if x5 < 0x80 {
		d.index += 6
		*v = x - (1<<7 + 1<<14 + 1<<21 + 1<<28 + 1<<35)
		return nil
	}

	if len(d.buf) <= d.index+6 {
		return errors.New("error decoding uint64")
	}
	x6 := buf[i+6]
	x += uint64(x6) << 42
	if x6 < 0x80 {
		d.index += 7
		*v = x - (1<<7 + 1<<14 + 1<<21 + 1<<28 + 1<<35 + 1<<42)
		return nil
	}

	if len(d.buf) <= d.index+7 {
		return errors.New("error decoding uint64")
	}
	x7 := buf[i+7]
	x += uint64(x7) << 49
	if x7 < 0x80 {
		d.index += 8
		*v = x - (1<<7 + 1<<14 + 1<<21 + 1<<28 + 1<<35 + 1<<42 + 1<<49)
		return nil
	}

	if len(d.buf) <= d.index+8 {
		return errors.New("error decoding uint64")
	}
	d.index += 9
	*v = x + uint64(buf[i+8])<<56 - (1<<7 + 1<<14 + 1<<21 + 1<<28 + 1<<35 + 1<<42 + 1<<49 + 1<<56)
	return nil
}

func (d *Decoder) decUint16(v *uint16) error {
	if len(d.buf) <= d.index {
		return errors.New("error decoding uint16")
	}
	buf, i := d.buf, d.index
	x := uint16(buf[i])
	if x < 0x80 {
		d.index++
		*v = x
		return nil
	}

	if len(d.buf) <= d.index+1 {
		return errors.New("error decoding uint16")
	}
	x1 := buf[i+1]
	x += uint16(x1) << 7
	if x1 < 0x80 {
		d.index += 2
		*v = x - 1<<7
		return nil
	}

	if len(d.buf) <= d.index+2 {
		return errors.New("error decoding uint16")
	}
	d.index += 3
	*v = x + uint16(buf[i+2])<<14 - (1<<7 + 1<<14)
	return nil
}

func (d *Decoder) decUint32(v *uint32) error {

	if len(d.buf) <= d.index {
		return errors.New("error decoding uint32")
	}
	buf, i := d.buf, d.index
	x := uint32(buf[i])
	if x < 0x80 {
		d.index++
		*v = x
		return nil
	}

	if len(d.buf) <= d.index+1 {
		return errors.New("error decoding uint32")
	}
	x1 := buf[i+1]
	x += uint32(x1) << 7
	if x1 < 0x80 {
		d.index += 2
		*v = x - 1<<7
		return nil
	}

	if len(d.buf) <= d.index+2 {
		return errors.New("error decoding uint32")
	}
	x2 := buf[i+2]
	x += uint32(x2) << 14
	if x2 < 0x80 {
		d.index += 3
		*v = x - (1<<7 + 1<<14)
		return nil
	}

	if len(d.buf) <= d.index+3 {
		return errors.New("error decoding uint32")
	}
	x3 := buf[i+3]
	x += uint32(x3) << 21
	if x3 < 0x80 {
		d.index += 4
		*v = x - (1<<7 + 1<<14 + 1<<21)
		return nil
	}

	if len(d.buf) <= d.index+4 {
		return errors.New("error decoding uint32")
	}
	x4 := buf[i+4]
	x += uint32(x4) << 28
	d.index += 5
	*v = x - (1<<7 + 1<<14 + 1<<21 + 1<<28)
	return nil
}

func (d *Decoder) decLength(v *int) error {

	var t uint32
	if err := d.decUint32(&t); err != nil {
		return err
	}
	*v = int(t)
	return nil
}

func (d *Decoder) decIsNotNil(v *bool) error {
	return d.decBool(v)
}

func decIgnore(*Decoder, unsafe.Pointer) error {
	return nil
}

func decBool(d *Decoder, p unsafe.Pointer) error {
	return d.decBool((*bool)(p))
}

func decInt(d *Decoder, p unsafe.Pointer) error {
	if err := d.decUint64((*uint64)(p)); err != nil {
		return err
	}

	*(*int)(p) = int(uint64ToInt64(*(*uint64)(p)))
	return nil
}

func decInt8(d *Decoder, p unsafe.Pointer) error {
	if len(d.buf) <= d.index+1 {
		return errors.New("error decoding int8")
	}

	*(*int8)(p) = int8(d.buf[d.index])
	d.index++
	return nil
}

func decInt16(d *Decoder, p unsafe.Pointer) error {
	if err := d.decUint16((*uint16)(p)); err != nil {
		return err
	}

	*(*int16)(p) = uint16ToInt16(*(*uint16)(p))
	return nil
}

func decInt32(d *Decoder, p unsafe.Pointer) error {
	if err := d.decUint32((*uint32)(p)); err != nil {
		return err
	}
	*(*int32)(p) = uint32ToInt32(*(*uint32)(p))
	return nil
}

func decInt64(d *Decoder, p unsafe.Pointer) error {
	if err := d.decUint64((*uint64)(p)); err != nil {
		return err
	}

	*(*int64)(p) = uint64ToInt64(*(*uint64)(p))
	return nil
}

func decUint(d *Decoder, p unsafe.Pointer) error {
	if err := d.decUint64((*uint64)(p)); err != nil {
		return err
	}

	*(*uint)(p) = uint(*(*uint64)(p))
	return nil
}

func decUint8(d *Decoder, p unsafe.Pointer) error {
	if len(d.buf) <= d.index+1 {
		return errors.New("error decoding uint8")
	}

	*(*uint8)(p) = d.buf[d.index]
	d.index++
	return nil
}

func decUint16(d *Decoder, p unsafe.Pointer) error {
	return d.decUint16((*uint16)(p))
}

func decUint32(d *Decoder, p unsafe.Pointer) error {
	return d.decUint32((*uint32)(p))
}

func decUint64(d *Decoder, p unsafe.Pointer) error {
	return d.decUint64((*uint64)(p))
}
func decUintptr(d *Decoder, p unsafe.Pointer) error {
	if err := d.decUint64((*uint64)(p)); err != nil {
		return err
	}

	*(*uintptr)(p) = uintptr(*(*uint64)(p))
	return nil
}

func decPointer(d *Decoder, p unsafe.Pointer) error {
	if err := d.decUint64((*uint64)(p)); err != nil {
		return err
	}

	*(*uintptr)(p) = uintptr(*(*uint64)(p))
	return nil
}

func decFloat32(d *Decoder, p unsafe.Pointer) error {
	if err := d.decUint32((*uint32)(p)); err != nil {
		return err
	}

	*(*float32)(p) = uint32ToFloat32(*(*uint32)(p))
	return nil
}

func decFloat64(d *Decoder, p unsafe.Pointer) error {
	if err := d.decUint64((*uint64)(p)); err != nil {
		return err
	}

	*(*float64)(p) = uint64ToFloat64(*(*uint64)(p))
	return nil
}

func decTime(d *Decoder, p unsafe.Pointer) error {

	if err := d.decUint64((*uint64)(p)); err != nil {
		return err
	}

	*(*time.Time)(p) = time.Unix(0, int64(*(*uint64)(p)))
	return nil
}

func decComplex64(d *Decoder, p unsafe.Pointer) error {
	return d.decUint64((*uint64)(p))
}

func decComplex128(d *Decoder, p unsafe.Pointer) error {
	if err := d.decUint64((*uint64)(p)); err != nil {
		return err
	}

	return d.decUint64((*uint64)(unsafe.Pointer(uintptr(p) + ptr1Size)))
}

func decString(d *Decoder, p unsafe.Pointer) error {
	var ut uint32
	if err := d.decUint32(&ut); err != nil {
		return err
	}

	l, val := int(ut), (*string)(p)

	if len(d.buf) < d.index+l {
		return errors.New("error decoding string")
	}
	*val = string(d.buf[d.index : d.index+l])
	d.index += l
	return nil
}

func decBytes(d *Decoder, p unsafe.Pointer) error {
	var notnil bool

	if err := d.decIsNotNil(&notnil); err != nil {
		return err
	}

	bytes := (*[]byte)(p)
	if notnil {
		var ut uint32
		if err := d.decUint32(&ut); err != nil {
			return err
		}

		l := int(ut)
		if len(d.buf) < d.index+l {
			return errors.New("error decoding bytes")
		}
		*bytes = d.buf[d.index : d.index+l]
		d.index += l
	} else if !isNil(p) {
		*bytes = nil
	}
	return nil
}
