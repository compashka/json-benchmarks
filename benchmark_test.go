package test

import (
	jsonv1 "encoding/json"
	"path"
	"testing"

	"github.com/bytedance/sonic"
	jsoniter "github.com/json-iterator/go"
	"github.com/mailru/easyjson"
)

func BenchmarkMarshal(b *testing.B) {
	for _, data := range testsdata {
		data := data
		b.Run(path.Join(data.name, "encoding/json"), func(b *testing.B) {
			fnc := compatTab[data.name]
			val := fnc()
			err := jsonv1.Unmarshal(data.cont, val)
			if err != nil {
				b.Fatal(err)
			}
			b.ReportAllocs()
			b.SetBytes(int64(len(data.cont)))
			b.ResetTimer()
			for idx := 0; idx < b.N; idx++ {
				_, err := jsonv1.Marshal(val)
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		b.Run(path.Join(data.name, "json-iterator/go"), func(b *testing.B) {
			fnc := compatTab[data.name]
			val := fnc()
			err := jsoniter.Unmarshal(data.cont, val)
			if err != nil {
				b.Fatal(err)
			}
			b.ReportAllocs()
			b.SetBytes(int64(len(data.cont)))
			b.ResetTimer()
			for idx := 0; idx < b.N; idx++ {
				_, err := jsoniter.Marshal(val)
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		b.Run(path.Join(data.name, "bytedance/sonic"), func(b *testing.B) {
			fnc := compatTab[data.name]
			val := fnc()
			err := sonic.Unmarshal(data.cont, val)
			if err != nil {
				b.Fatal(err)
			}
			b.ReportAllocs()
			b.SetBytes(int64(len(data.cont)))
			b.ResetTimer()
			for idx := 0; idx < b.N; idx++ {
				_, err := sonic.Marshal(val)
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		b.Run(path.Join(data.name, "mailru/easyjson"), func(b *testing.B) {
			fnc := easyJSONTab[data.name]
			val := fnc()
			err := easyjson.Unmarshal(data.cont, val)
			if err != nil {
				b.Fatal(err)
			}
			b.ReportAllocs()
			b.SetBytes(int64(len(data.cont)))
			b.ResetTimer()
			for idx := 0; idx < b.N; idx++ {
				_, err := easyjson.Marshal(val)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	for _, data := range testsdata {
		data := data
		b.Run(path.Join(data.name, "encoding/json"), func(b *testing.B) {
			fnc := compatTab[data.name]
			b.ReportAllocs()
			b.SetBytes(int64(len(data.cont)))
			b.ResetTimer()
			for idx := 0; idx < b.N; idx++ {
				err := jsonv1.Unmarshal(data.cont, fnc())
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		b.Run(path.Join(data.name, "json-iterator/go"), func(b *testing.B) {
			fnc := compatTab[data.name]
			b.ReportAllocs()
			b.SetBytes(int64(len(data.cont)))
			b.ResetTimer()
			for idx := 0; idx < b.N; idx++ {
				err := jsoniter.Unmarshal(data.cont, fnc())
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		b.Run(path.Join(data.name, "bytedance/sonic"), func(b *testing.B) {
			fnc := compatTab[data.name]
			b.ReportAllocs()
			b.SetBytes(int64(len(data.cont)))
			b.ResetTimer()
			for idx := 0; idx < b.N; idx++ {
				err := sonic.Unmarshal(data.cont, fnc())
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		b.Run(path.Join(data.name, "mailru/easyjson"), func(b *testing.B) {
			fnc := easyJSONTab[data.name]
			b.ReportAllocs()
			b.SetBytes(int64(len(data.cont)))
			b.ResetTimer()
			for idx := 0; idx < b.N; idx++ {
				err := easyjson.Unmarshal(data.cont, fnc())
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
