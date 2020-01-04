package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"runtime/pprof"
	"time"

	"github.com/jimyx17/gotiny"
	"github.com/niubaoshu/goutils"
)

type str struct {
	A map[int]map[int]string
	B []bool
	c int
}

type ET0 struct {
	s str
	F map[int]map[int]string
}

var (
	//_   = rand.Intn(1)
	now = time.Now()
	a   = "234234"
	i   = map[int]map[int]string{
		1: map[int]string{
			1: a,
		},
	}
	strs = `抵制西方的司法独立，有什么错？有人说马克思主义还是西方的，有本事别用啊。这都是犯了形而上学的错误，任何理论、思想都必须和中国国情相结合，和当前实际相结合。全部照搬照抄的教条主义王明已经试过一次，结果怎么样？歪解周强讲话，不是蠢就是别有用心，蠢的可以教育，别有用心就该打倒`
	st   = str{A: i, B: []bool{true, false, false, false, false, true, true, false, true, false, true}, c: 234234}
	//st     = str{c: 234234}
	et0      = ET0{s: st, F: i}
	stp      = &st
	stpp     = &stp
	nilslice []byte
	slice    = []byte{1, 2, 3}
	mapt     = map[int]int{0: 1, 1: 2, 2: 3, 3: 4}
	nilmap   map[int][]byte
	nilptr   *map[int][]string
	inta          = 2
	ptrint   *int = &inta
	nilint   *int
	vs       = []interface{}{
		ptrint,
		strs,
		`Xi Jinping leaves Beijing for state visit to Swiss Federation
		Attend World Economic Forum Annual Meeting 2017 and visit international organizations in Switzerland
		Xinhua News Agency, Beijing, January 15th. On the morning of January 15th, President Xi Jinping left Beijing on a special plane. He was invited to the Swiss Federal Council chaired by Leuthard and paid a state visit to Switzerland. Executive Chairman Schwab invited to attend the 2017 World Economic Forum Annual Meeting in Davos; visited the UN headquarters in Geneva at the invitation of UN Secretary-General Guterres, Director-General of the World Health Organization Chen Feng Fuzhen and President of the International Olympic Committee Bach , World Health Organization, International Olympic Committee.
		Accompanying Xi Jinping are: Peng Liyuan, wife of President Xi Jinping, Wang Huning, member of the Political Bureau of the CPC Central Committee and director of the Central Policy Research Office, Li Zhanshu, member of the Political Bureau of the CPC Central Committee, secretary of the Central Secretariat, director of the Central Office, and state councilor Yang Jiechi. Back to Tencent Home >>`,
		true,
		false,
		int(123456),
		int8(123),
		int16(-12345),
		int32(123456),
		int64(-1234567),
		int64(1<<63 - 1),
		int64(rand.Int63()),
		uint(123),
		uint8(123),
		uint16(12345),
		uint32(123456),
		uint64(1234567),
		uint64(1<<64 - 1),
		uint64(rand.Uint32() * rand.Uint32()),
		uintptr(12345678),
		float32(1.2345),
		float64(1.2345678),
		complex64(1.2345 + 2.3456i),
		complex128(1.2345678 + 2.3456789i),
		string("hello, Japan"),
		string("9b899bec35bc6bb8"),
		inta,
		[][][][3][][3]int{{{{{{2, 3}}}}}},
		map[int]map[int]map[int]map[int]map[int]map[int]map[int]map[int]int{1: {1: {1: {1: {1: {1: {1: {1: 2}}}}}}}},
		map[int]map[int]int{1: {2: 3}},
		[][]bool{},
		[]byte("hello，Chinese"),
		[][]byte{[]byte("hello"), []byte("world")},
		[4]string{"2324", "23423", "Buffy", "《：LSESERsef pink ask me 2D cow"},
		map[int]string{1: "h", 2: "h", 3: "nihao"},
		map[string]map[int]string{"werwer": {1: "Shout"}, "Char": {2: "world"}},
		a,
		i,
		&i,
		st,
		stp,
		stpp,
		struct{}{},
		[][][]struct{}{},
		struct {
			a, C int
		}{1, 2},
		et0,
		[100]int{},
		now,
		ptrint,
		nilmap,
		nilslice,
		nilptr,
		nilint,
		slice,
		mapt,
	}
	e, _ = gotiny.NewEncoder(vs...)
	d, _ = gotiny.NewDecoder(vs...)

	spvals = make([]interface{}, len(vs))
	rpvals = make([]interface{}, len(vs))
	c      = goutils.NewComparer()

	buf = make([]byte, 0, 2048)
)

func init() {

	for i := 0; i < len(vs); i++ {
		typ := reflect.TypeOf(vs[i])
		temp := reflect.New(typ)
		temp.Elem().Set(reflect.ValueOf(vs[i]))
		spvals[i] = temp.Interface()

		if i == len(vs)-2 {
			a := make([]byte, 15)
			rpvals[i] = &a
		} else if i == len(vs)-1 {
			//a := map[int]int{111: 233, 6: 7}
			a := map[int]int{}
			rpvals[i] = &a
		} else {
			rpvals[i] = reflect.New(typ).Interface()
		}
	}
	e.AppendTo(buf[:0])
}

func main() {
	f, err := os.Create("cpuprofile.pprof")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	for i := 0; i < 1000; i++ {
		for i := 0; i < 1000; i++ {
			e.AppendTo(buf[:0])
			t, _ := e.Encode(spvals...)
			d.Decode(t, rpvals...)
			for i, result := range rpvals {
				r := reflect.ValueOf(result).Elem().Interface()
				if Assert(vs[i], r) != nil {
					fmt.Println(err)
				}
			}
		}
	}
}

func Assert(x, y interface{}) error {
	if !c.DeepEqual(x, y) {
		return fmt.Errorf("\n exp type =  %T; value = %#v;\n got type = %T; value = %#v ", x, x, y, y)
	}
	return nil
}
