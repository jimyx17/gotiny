package gotiny

import (
	"errors"
	"reflect"
	"sync"
	"time"
	"unsafe"
)

type encEng func(*Encoder, unsafe.Pointer) // Encoders

var (
	rt2encEng = map[reflect.Type]encEng{
		reflect.TypeOf((*bool)(nil)).Elem():           encBool,
		reflect.TypeOf((*int)(nil)).Elem():            encInt,
		reflect.TypeOf((*int8)(nil)).Elem():           encInt8,
		reflect.TypeOf((*int16)(nil)).Elem():          encInt16,
		reflect.TypeOf((*int32)(nil)).Elem():          encInt32,
		reflect.TypeOf((*int64)(nil)).Elem():          encInt64,
		reflect.TypeOf((*uint)(nil)).Elem():           encUint,
		reflect.TypeOf((*uint8)(nil)).Elem():          encUint8,
		reflect.TypeOf((*uint16)(nil)).Elem():         encUint16,
		reflect.TypeOf((*uint32)(nil)).Elem():         encUint32,
		reflect.TypeOf((*uint64)(nil)).Elem():         encUint64,
		reflect.TypeOf((*uintptr)(nil)).Elem():        encUintptr,
		reflect.TypeOf((*unsafe.Pointer)(nil)).Elem(): encPointer,
		reflect.TypeOf((*float32)(nil)).Elem():        encFloat32,
		reflect.TypeOf((*float64)(nil)).Elem():        encFloat64,
		reflect.TypeOf((*complex64)(nil)).Elem():      encComplex64,
		reflect.TypeOf((*complex128)(nil)).Elem():     encComplex128,
		reflect.TypeOf((*[]byte)(nil)).Elem():         encBytes,
		reflect.TypeOf((*string)(nil)).Elem():         encString,
		reflect.TypeOf((*time.Time)(nil)).Elem():      encTime,
		reflect.TypeOf((*struct{})(nil)).Elem():       encIgnore,
		reflect.TypeOf(nil):                           encIgnore,
	}

	encEngines = [...]encEng{
		reflect.Invalid:       encIgnore,
		reflect.Bool:          encBool,
		reflect.Int:           encInt,
		reflect.Int8:          encInt8,
		reflect.Int16:         encInt16,
		reflect.Int32:         encInt32,
		reflect.Int64:         encInt64,
		reflect.Uint:          encUint,
		reflect.Uint8:         encUint8,
		reflect.Uint16:        encUint16,
		reflect.Uint32:        encUint32,
		reflect.Uint64:        encUint64,
		reflect.Uintptr:       encUintptr,
		reflect.UnsafePointer: encPointer,
		reflect.Float32:       encFloat32,
		reflect.Float64:       encFloat64,
		reflect.Complex64:     encComplex64,
		reflect.Complex128:    encComplex128,
		reflect.String:        encString,
	}

	encLock sync.RWMutex
)

func UnusedUnixNanoEncodeTimeType() {
	delete(rt2encEng, reflect.TypeOf((*time.Time)(nil)).Elem())
	delete(rt2decEng, reflect.TypeOf((*time.Time)(nil)).Elem())
}

func getEncEngine(rt reflect.Type) (engine encEng, err error) {
	encLock.RLock()
	engine = rt2encEng[rt]
	encLock.RUnlock()
	if engine != nil {
		return
	}
	encLock.Lock()
	err = buildEncEngine(rt, &engine)
	encLock.Unlock()
	return
}

func buildEncEngine(rt reflect.Type, engPtr *encEng) (err error) {
	engine := rt2encEng[rt]
	if engine != nil {
		*engPtr = engine
		return
	}

	if engine, _ = implementOtherSerializer(rt); engine != nil {
		rt2encEng[rt] = engine
		*engPtr = engine
		return
	}

	kind := rt.Kind()
	var eEng encEng
	switch kind {
	case reflect.Ptr:
		defer buildEncEngine(rt.Elem(), &eEng)
		engine = func(e *Encoder, p unsafe.Pointer) {
			isNotNil := !isNil(p)
			e.encIsNotNil(isNotNil)
			if isNotNil {
				ref, build := getReference(e, *(*unsafe.Pointer)(p))
				e.encBool(build)
				e.encUint16(ref)
				if build {
					eEng(e, *(*unsafe.Pointer)(p))
				}
			}
		}

	case reflect.Array:
		et, l := rt.Elem(), rt.Len()
		defer buildEncEngine(et, &eEng)
		size := et.Size()
		engine = func(e *Encoder, p unsafe.Pointer) {
			for i := 0; i < l; i++ {
				eEng(e, unsafe.Pointer(uintptr(p)+uintptr(i)*size))
			}
		}

	case reflect.Slice:
		et := rt.Elem()
		size := et.Size()
		defer buildEncEngine(et, &eEng)
		engine = func(e *Encoder, p unsafe.Pointer) {
			isNotNil := !isNil(p)
			e.encIsNotNil(isNotNil)
			if isNotNil {
				header := (*reflect.SliceHeader)(p)
				l := header.Len
				e.encLength(l)
				for i := 0; i < l; i++ {
					eEng(e, unsafe.Pointer(header.Data+uintptr(i)*size))
				}
			}
		}
	case reflect.Map:
		var kEng encEng
		defer buildEncEngine(rt.Key(), &kEng)
		defer buildEncEngine(rt.Elem(), &eEng)
		engine = func(e *Encoder, p unsafe.Pointer) {
			isNotNil := !isNil(p)
			e.encIsNotNil(isNotNil)
			if isNotNil {
				v := reflect.NewAt(rt, p).Elem()
				e.encLength(v.Len())
				keys := v.MapKeys()
				for i := 0; i < len(keys); i++ {
					val := v.MapIndex(keys[i])
					kEng(e, getUnsafePointer(&keys[i]))
					eEng(e, getUnsafePointer(&val))
				}
			}
		}
	case reflect.Struct:
		fields, offs := getFieldType(rt, 0)
		nf := len(fields)
		fEngines := make([]encEng, nf)
		defer func() {
			for i := 0; i < nf; i++ {
				buildEncEngine(fields[i], &fEngines[i])
			}
		}()
		engine = func(e *Encoder, p unsafe.Pointer) {
			for i := 0; i < len(fEngines) && i < len(offs); i++ {
				fEngines[i](e, unsafe.Pointer(uintptr(p)+offs[i]))
			}
		}
	case reflect.Interface:
		if rt.NumMethod() > 0 {
			engine = func(e *Encoder, p unsafe.Pointer) {
				isNotNil := !isNil(p)
				e.encIsNotNil(isNotNil)
				if isNotNil {
					var tmp encEng
					var name string
					v := reflect.ValueOf(*(*interface {
						M()
					})(p))
					et := v.Type()
					name, err = getNameOfType(et)
					if err != nil {
						return
					}
					e.encString(name)
					tmp, err = getEncEngine(et)
					if err != nil {
						return
					}
					tmp(e, getUnsafePointer(&v))
				}
			}
		} else {
			engine = func(e *Encoder, p unsafe.Pointer) {
				isNotNil := !isNil(p)
				e.encIsNotNil(isNotNil)
				if isNotNil {
					var tmp encEng
					var name string
					v := reflect.ValueOf(*(*interface{})(p))
					et := v.Type()
					name, err = getNameOfType(et)
					e.encString(name)

					tmp, err = getEncEngine(et)
					if err != nil {
						return
					}
					tmp(e, getUnsafePointer(&v))
				}
			}
		}
	case reflect.Chan, reflect.Func:
		err = errors.New("not support " + rt.String() + " type")
		return
	default:
		engine = encEngines[kind]
	}
	rt2encEng[rt] = engine
	*engPtr = engine
	return
}
