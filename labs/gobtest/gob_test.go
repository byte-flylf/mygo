// A test case to compare encode/gob & encoding/json performance
package gobtest

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
)

type St struct {
	I1, I2, I3, I4 int
	S1, S2, S3, S4 string
	M1             map[string]int
	M2             map[string]string
	L1             []int
	L2             []string
}

func printBR(name string, br testing.BenchmarkResult) {
	fmt.Println(name, ": Total used:",
		br.T,
		", Per operation:",
		int64(br.T)/int64(br.N),
		"ns/ops", br.N)
}

func TestBench(t *testing.T) {
	benchGobEncodeList()
	benchJsonEncodeList()

	benchGobEncodeListString()
	benchJsonEncodeListString()

	benchGobEncodeMap()
	benchJsonEncodeMap()
	benchGobEncodeMapInt()

	benchGobDecodeList()
	benchJsonDecodeList()

	benchGobDecodeMap()
	benchJsonDecodeMap()
}

const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func keygen(i int, l int) string {
	var res string

	clen := len(characters)
	rand.Seed(int64(i))
	for idx := 0; idx < l; idx++ {
		res += string(characters[rand.Intn(clen)])
	}
	return res
}

func gobEncode(data interface{}) {
	var (
		enc *gob.Encoder
		buf bytes.Buffer
		err error
	)
	enc = gob.NewEncoder(&buf)
	err = enc.Encode(data)
	if err != nil {
		panic(err.Error())
	}
}

func gobEncodeList(b *testing.B, l int) {
	s := make([]int, l)
	for idx, _ := range s {
		s[idx] = rand.Int()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gobEncode(s)
	}
}

func gobEncodeListNi(b *testing.B, l int) {
	s := make([]int, l)
	for idx, _ := range s {
		s[idx] = rand.Int()
	}

	var (
		enc *gob.Encoder
		buf bytes.Buffer
		err error
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		enc = gob.NewEncoder(&buf)
		err = enc.Encode(s)
		if err != nil {
			panic(err.Error())
		}
	}
}

func gobEncodeListString(b *testing.B, l int, slen int) {
	s := make([]string, l)
	for idx, _ := range s {
		s[idx] = keygen(idx, slen)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gobEncode(s)
	}
}

func gobEncodeListStringNi(b *testing.B, l int, slen int) {
	s := make([]string, l)
	for idx, _ := range s {
		s[idx] = keygen(idx, slen)
	}

	var (
		enc *gob.Encoder
		buf bytes.Buffer
		err error
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		enc = gob.NewEncoder(&buf)
		err = enc.Encode(s)
		if err != nil {
			panic(err.Error())
		}
	}
}

func gobEncodeMap(b *testing.B, l int, kl int) {
	s := make(map[string]int)
	for idx := 0; idx < l; idx++ {
		key := keygen(idx, kl)
		s[key] = rand.Int()
	}
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		gobEncode(s)
	}
}

func gobEncodeMapNi(b *testing.B, l int, kl int) {
	s := make(map[string]int)
	for idx := 0; idx < l; idx++ {
		key := keygen(idx, kl)
		s[key] = rand.Int()
	}

	var (
		enc *gob.Encoder
		buf bytes.Buffer
		err error
	)

	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		enc = gob.NewEncoder(&buf)
		err = enc.Encode(s)
		if err != nil {
			panic(err.Error())
		}
	}
}

func gobEncodeMapInt(b *testing.B, l int) {
	s := make(map[int]int)
	for idx := 0; idx < l; idx++ {
		s[idx] = rand.Int()
	}
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		gobEncode(s)
	}
}

func gobEncodeMapIntNi(b *testing.B, l int) {
	s := make(map[int]int)
	for idx := 0; idx < l; idx++ {
		s[idx] = rand.Int()
	}
	var (
		enc *gob.Encoder
		buf bytes.Buffer
		err error
	)
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		enc = gob.NewEncoder(&buf)
		err = enc.Encode(s)
		if err != nil {
			panic(err.Error())
		}
		//gobEncode(s)
	}
}

func benchGobEncodeListString() {
	printBR("Gob Encode []string(10)(2000)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeListString(b, 2000, 10)
		}))
	printBR("Gob Encode []string(20)(2000)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeListString(b, 2000, 20)
		}))
	printBR("Gob Encode-ni []string(10)(2000)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeListStringNi(b, 2000, 10)
		}))
	printBR("Gob Encode-ni []string(20)(2000)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeListStringNi(b, 2000, 20)
		}))

	printBR("Gob Encode []string(10)(1000)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeListString(b, 1000, 10)
		}))
	printBR("Gob Encode []string(20)(1000)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeListString(b, 1000, 20)
		}))
	printBR("Gob Encode-ni []string(10)(1000)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeListStringNi(b, 1000, 10)
		}))
	printBR("Gob Encode-ni []string(20)(1000)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeListStringNi(b, 1000, 20)
		}))

	printBR("Gob Encode []string(10)(100)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeListString(b, 100, 10)
		}))
	printBR("Gob Encode []string(20)(100)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeListString(b, 100, 20)
		}))
	printBR("Gob Encode-ni []string(10)(100)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeListStringNi(b, 100, 10)
		}))
	printBR("Gob Encode-ni []string(20)(100)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeListStringNi(b, 100, 20)
		}))

	printBR("Gob Encode []string(10)(10)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeListString(b, 10, 10)
		}))
	printBR("Gob Encode []string(20)(10)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeListString(b, 10, 20)
		}))
	printBR("Gob Encode-ni []string(10)(10)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeListStringNi(b, 10, 10)
		}))
	printBR("Gob Encode-ni []string(20)(10)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeListStringNi(b, 10, 20)
		}))
}

func benchGobEncodeList() {
	printBR("Gob Encode []int(2000)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeList(b, 2000)
		}))
	printBR("Gob Encode-ni []int(2000)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeListNi(b, 2000)
		}))

	printBR("Gob Encode []int(1000)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeList(b, 1000)
		}))
	printBR("Gob Encode-ni []int(1000)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeListNi(b, 1000)
		}))

	printBR("Gob Encode []int(100)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeList(b, 100)
		}))
	printBR("Gob Encode-ni []int(100)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeListNi(b, 100)
		}))

	printBR("Gob Encode []int(10)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeList(b, 10)
		}))
	printBR("Gob Encode-ni []int(10)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeListNi(b, 10)
		}))
}

func benchGobEncodeMap() {
	printBR("Gob Encode map[string(10)]int(2000)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeMap(b, 2000, 10)
		}))
	printBR("Gob Encode map[string(20)]int(2000)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeMap(b, 2000, 20)
		}))
	printBR("Gob Encode-ni map[string(10)]int(2000)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeMapNi(b, 2000, 10)
		}))
	printBR("Gob Encode-ni map[string(20)]int(2000)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeMapNi(b, 2000, 20)
		}))

	printBR("Gob Encode map[string(10)]int(1000)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeMap(b, 1000, 10)
		}))
	printBR("Gob Encode map[string(20)]int(1000)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeMap(b, 1000, 20)
		}))
	printBR("Gob Encode-ni map[string(10)]int(1000)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeMapNi(b, 1000, 10)
		}))
	printBR("Gob Encode-ni map[string(20)]int(1000)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeMapNi(b, 1000, 20)
		}))

	printBR("Gob Encode map[string(10)]int(100)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeMap(b, 100, 10)
		}))
	printBR("Gob Encode map[string(20)]int(100)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeMap(b, 100, 20)
		}))
	printBR("Gob Encode-ni map[string(10)]int(100)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeMapNi(b, 100, 10)
		}))
	printBR("Gob Encode-ni map[string(20)]int(100)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeMapNi(b, 100, 20)
		}))

	printBR("Gob Encode map[string(10)]int(10)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeMap(b, 10, 10)
		}))
	printBR("Gob Encode map[string(20)]int(10)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeMap(b, 10, 20)
		}))
	printBR("Gob Encode-ni map[string(10)]int(10)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeMapNi(b, 10, 10)
		}))
	printBR("Gob Encode-ni map[string(20)]int(10)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeMapNi(b, 10, 20)
		}))
}

func benchGobEncodeMapInt() {
	printBR("Gob Encode map[int]int(2000)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeMapInt(b, 2000)
		}))
	printBR("Gob Encode-ni map[int]int(2000)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeMapIntNi(b, 2000)
		}))

	printBR("Gob Encode map[int]int(1000)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeMapInt(b, 1000)
		}))
	printBR("Gob Encode-ni map[int]int(1000)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeMapIntNi(b, 1000)
		}))

	printBR("Gob Encode map[int]int(100)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeMapInt(b, 100)
		}))
	printBR("Gob Encode-ni map[int]int(100)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeMapIntNi(b, 100)
		}))

	printBR("Gob Encode map[int]int(10)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeMapInt(b, 10)
		}))
	printBR("Gob Encode-ni map[int]int(10)",
		testing.Benchmark(func(b *testing.B) {
			gobEncodeMapIntNi(b, 10)
		}))
}

func gobDecodeList(b *testing.B, l int) {
	var (
		enc *gob.Encoder
		buf bytes.Buffer
		err error
	)
	s := make([]int, l)
	for i, _ := range s {
		s[i] = rand.Int()
	}

	enc = gob.NewEncoder(&buf)
	err = enc.Encode(s)
	if err != nil {
		panic(err.Error())
	}

	var dest []int
	bs := buf.Bytes()
	b.ResetTimer()
	// decode
	for i := 0; i < b.N; i++ {
		nbuffer := bytes.NewBuffer(bs)
		dec := gob.NewDecoder(nbuffer)
		err := dec.Decode(&dest)
		if err != nil {
			panic(err.Error())
		}
	}
}

func gobDecodeMap(b *testing.B, l int, kl int) {
	var (
		enc *gob.Encoder
		buf bytes.Buffer
		err error
	)
	s := make(map[string]int)
	for i := 0; i < l; i++ {
		key := keygen(i, kl)
		s[key] = rand.Int()
	}

	enc = gob.NewEncoder(&buf)
	err = enc.Encode(s)
	if err != nil {
		panic(err.Error())
	}

	var dest map[string]int

	bs := buf.Bytes()
	b.ResetTimer()
	// decode
	for i := 0; i < b.N; i++ {
		nbuffer := bytes.NewBuffer(bs)
		dec := gob.NewDecoder(nbuffer)
		err := dec.Decode(&dest)
		if err != nil {
			panic(err.Error())
		}
	}
}

func benchGobDecodeList() {
	printBR("Gob Decode []int(2000)",
		testing.Benchmark(func(b *testing.B) {
			gobDecodeList(b, 2000)
		}))
	printBR("Gob Decode []int(1000)",
		testing.Benchmark(func(b *testing.B) {
			gobDecodeList(b, 1000)
		}))
	printBR("Gob Decode []int(100)",
		testing.Benchmark(func(b *testing.B) {
			gobDecodeList(b, 100)
		}))
	printBR("Gob Decode []int(10)",
		testing.Benchmark(func(b *testing.B) {
			gobDecodeList(b, 10)
		}))
}

func benchGobDecodeMap() {
	printBR("Gob Decode map[string(10)]int(2000)",
		testing.Benchmark(func(b *testing.B) {
			gobDecodeMap(b, 2000, 10)
		}))
	printBR("Gob Decode map[string(20)]int(2000)",
		testing.Benchmark(func(b *testing.B) {
			gobDecodeMap(b, 2000, 20)
		}))

	printBR("Gob Decode map[string(10)]int(1000)",
		testing.Benchmark(func(b *testing.B) {
			gobDecodeMap(b, 1000, 10)
		}))
	printBR("Gob Decode map[string(20)]int(1000)",
		testing.Benchmark(func(b *testing.B) {
			gobDecodeMap(b, 1000, 20)
		}))

	printBR("Gob Decode map[string(10)]int(100)",
		testing.Benchmark(func(b *testing.B) {
			gobDecodeMap(b, 100, 10)
		}))
	printBR("Gob Decode map[string(20)]int(100)",
		testing.Benchmark(func(b *testing.B) {
			gobDecodeMap(b, 100, 20)
		}))

	printBR("Gob Decode map[string(10)]int(10)",
		testing.Benchmark(func(b *testing.B) {
			gobDecodeMap(b, 10, 10)
		}))
	printBR("Gob Decode map[string(20)]int(10)",
		testing.Benchmark(func(b *testing.B) {
			gobDecodeMap(b, 10, 20)
		}))
}

func jsonEncode(data interface{}) {
	_, err := json.Marshal(data)
	if err != nil {
		panic(err.Error())
	}
}

func jsonEncodeList(b *testing.B, l int) {
	s := make([]int, l)
	for idx := 0; idx < l; idx++ {
		s[idx] = rand.Int()
	}
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		jsonEncode(s)
	}
}
func jsonEncodeListNi(b *testing.B, l int) {
	s := make([]int, l)
	for idx := 0; idx < l; idx++ {
		s[idx] = rand.Int()
	}
	var err error
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err = json.Marshal(s)
		if err != nil {
			panic(err.Error())
		}
	}
}

func jsonEncodeListString(b *testing.B, l int, slen int) {
	s := make([]string, l)
	for idx := 0; idx < l; idx++ {
		s[idx] = keygen(idx, slen)
	}
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		jsonEncode(s)
	}
}

func jsonEncodeListStringNi(b *testing.B, l int, slen int) {
	s := make([]string, l)
	for idx := 0; idx < l; idx++ {
		s[idx] = keygen(idx, slen)
	}
	var err error
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err = json.Marshal(s)
		if err != nil {
			panic(err.Error())
		}
	}
}

func jsonEncodeMap(b *testing.B, l int, kl int) {
	s := make(map[string]int)
	for idx := 0; idx < l; idx++ {
		key := keygen(idx, kl)
		s[key] = rand.Int()
	}
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		jsonEncode(s)
	}
}

func jsonEncodeMapNi(b *testing.B, l int, kl int) {
	s := make(map[string]int)
	for idx := 0; idx < l; idx++ {
		key := keygen(idx, kl)
		s[key] = rand.Int()
	}
	var err error
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err = json.Marshal(s)
		if err != nil {
			panic(err.Error())
		}
	}
}

func benchJsonEncodeList() {
	printBR("JSON Encode []int(2000)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeList(b, 2000)
		}))
	printBR("JSON Encode-ni []int(2000)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeListNi(b, 2000)
		}))

	printBR("JSON Encode []int(1000)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeList(b, 1000)
		}))
	printBR("JSON Encode-ni []int(1000)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeListNi(b, 1000)
		}))

	printBR("JSON Encode []int(100)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeList(b, 100)
		}))
	printBR("JSON Encode-ni []int(100)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeListNi(b, 100)
		}))

	printBR("JSON Encode []int(10)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeList(b, 10)
		}))
	printBR("JSON Encode-ni []int(10)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeListNi(b, 10)
		}))
}

func benchJsonEncodeListString() {
	printBR("JSON Encode []string(10)(2000)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeListString(b, 2000, 10)
		}))
	printBR("JSON Encode []string(20)(2000)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeListString(b, 2000, 20)
		}))
	printBR("JSON Encode-ni []string(10)(2000)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeListStringNi(b, 2000, 10)
		}))
	printBR("JSON Encode-ni []string(20)(2000)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeListStringNi(b, 2000, 20)
		}))

	printBR("JSON Encode []string(10)(1000)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeListString(b, 1000, 10)
		}))
	printBR("JSON Encode []string(20)(1000)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeListString(b, 1000, 20)
		}))
	printBR("JSON Encode-ni []string(10)(1000)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeListStringNi(b, 1000, 10)
		}))
	printBR("JSON Encode-ni []string(20)(1000)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeListStringNi(b, 1000, 20)
		}))

	printBR("JSON Encode []string(10)(100)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeListString(b, 100, 10)
		}))
	printBR("JSON Encode []string(20)(100)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeListString(b, 100, 20)
		}))
	printBR("JSON Encode-ni []string(10)(100)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeListStringNi(b, 100, 10)
		}))
	printBR("JSON Encode-ni []string(20)(100)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeListStringNi(b, 100, 20)
		}))

	printBR("JSON Encode []string(10)(10)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeListString(b, 10, 10)
		}))
	printBR("JSON Encode []string(20)(10)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeListString(b, 10, 20)
		}))
	printBR("JSON Encode-ni []string(10)(10)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeListStringNi(b, 10, 10)
		}))
	printBR("JSON Encode-ni []string(20)(10)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeListStringNi(b, 10, 20)
		}))
}

func benchJsonEncodeMap() {
	printBR("JSON Encode map[string(10)]int(2000)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeMap(b, 2000, 10)
		}))
	printBR("JSON Encode map[string(20)]int(2000)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeMap(b, 2000, 20)
		}))
	printBR("JSON Encode-ni map[string(10)]int(2000)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeMapNi(b, 2000, 10)
		}))
	printBR("JSON Encode-ni map[string(20)]int(2000)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeMapNi(b, 2000, 20)
		}))

	printBR("JSON Encode-ni map[string(10)]int(1000)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeMapNi(b, 1000, 10)
		}))
	printBR("JSON Encode-ni map[string(20)]int(1000)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeMapNi(b, 1000, 20)
		}))

	printBR("JSON Encode-ni map[string(10)]int(100)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeMapNi(b, 100, 10)
		}))
	printBR("JSON Encode-ni map[string(20)]int(100)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeMapNi(b, 100, 20)
		}))

	printBR("JSON Encode map[string(10)]int(10)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeMap(b, 10, 10)
		}))
	printBR("JSON Encode map[string(20)]int(10)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeMap(b, 10, 20)
		}))
	printBR("JSON Encode-ni map[string(10)]int(10)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeMapNi(b, 10, 10)
		}))
	printBR("JSON Encode-ni map[string(20)]int(10)",
		testing.Benchmark(func(b *testing.B) {
			jsonEncodeMapNi(b, 10, 20)
		}))
}

func jsonDecodeListInt(buf []byte) {
	var dest []int

	err := json.Unmarshal(buf, &dest)
	if err != nil {
		panic(err.Error())
	}
}

func jsonDecodeList(b *testing.B, l int) {
	s := make([]int, l)
	for idx := 0; idx < l; idx++ {
		s[idx] = rand.Int()
	}
	buf, err := json.Marshal(s)
	if err != nil {
		panic(err.Error())
	}

	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		jsonDecodeListInt(buf)
	}
}

func jsonDecodeMap(b *testing.B, l int, kl int) {
	s := make(map[string]int)
	for idx := 0; idx < l; idx++ {
		key := keygen(idx, kl)
		s[key] = rand.Int()
	}
	buf, err := json.Marshal(s)
	if err != nil {
		panic(err.Error())
	}

	var dest map[string]int
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		err = json.Unmarshal(buf, &dest)
		if err != nil {
			panic(err.Error())
		}
	}
}

func benchJsonDecodeList() {
	printBR("JSON Decode []int(2000)",
		testing.Benchmark(func(b *testing.B) {
			jsonDecodeList(b, 2000)
		}))
	printBR("JSON Decode []int(1000)",
		testing.Benchmark(func(b *testing.B) {
			jsonDecodeList(b, 1000)
		}))
	printBR("JSON Decode []int(100)",
		testing.Benchmark(func(b *testing.B) {
			jsonDecodeList(b, 100)
		}))
	printBR("JSON Decode []int(10)",
		testing.Benchmark(func(b *testing.B) {
			jsonDecodeList(b, 10)
		}))
}

func benchJsonDecodeMap() {
	printBR("JSON Decode map[string(10)]int(2000)",
		testing.Benchmark(func(b *testing.B) {
			jsonDecodeMap(b, 2000, 10)
		}))
	printBR("JSON Decode map[string(20)]int(2000)",
		testing.Benchmark(func(b *testing.B) {
			jsonDecodeMap(b, 2000, 20)
		}))

	printBR("JSON Decode map[string(10)]int(1000)",
		testing.Benchmark(func(b *testing.B) {
			jsonDecodeMap(b, 1000, 10)
		}))
	printBR("JSON Decode map[string(20)]int(1000)",
		testing.Benchmark(func(b *testing.B) {
			jsonDecodeMap(b, 1000, 20)
		}))

	printBR("JSON Decode map[string(10)]int(100)",
		testing.Benchmark(func(b *testing.B) {
			jsonDecodeMap(b, 100, 10)
		}))
	printBR("JSON Decode map[string(20)]int(100)",
		testing.Benchmark(func(b *testing.B) {
			jsonDecodeMap(b, 100, 20)
		}))

	printBR("JSON Decode map[string(10)]int(10)",
		testing.Benchmark(func(b *testing.B) {
			jsonDecodeMap(b, 10, 10)
		}))
	printBR("JSON Decode map[string(20)]int(10)",
		testing.Benchmark(func(b *testing.B) {
			jsonDecodeMap(b, 10, 20)
		}))
}
