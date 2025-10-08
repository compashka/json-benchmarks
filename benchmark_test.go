package main

import (
	"path"
	"test/internal/gzip"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/go-json-experiment/json"
	jsoniter "github.com/json-iterator/go"
	"github.com/mailru/easyjson"
)

func BenchmarkMarshal(b *testing.B) {
	for _, d := range testdata {
		for _, libName := range libs {
			b.Run(path.Join(d.name, libName), func(b *testing.B) {
				val := d.new()
				rawData := d.rawDataGetter()
				err := unmarshalWith(libName, rawData, val)
				if err != nil {
					b.Fatal(err)
				}

				b.ReportAllocs()
				b.SetBytes(int64(len(rawData)))
				b.ResetTimer()
				for idx := 0; idx < b.N; idx++ {
					_, err := marshalWith(libName, val)
					if err != nil {
						b.Fatal(err)
					}
				}
			})
		}
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	for _, d := range testdata {
		for _, libName := range libs {
			b.Run(path.Join(d.name, libName), func(b *testing.B) {
				val := d.new()
				rawData := d.rawDataGetter()

				b.ReportAllocs()
				b.SetBytes(int64(len(rawData)))
				b.ResetTimer()
				for idx := 0; idx < b.N; idx++ {
					err := unmarshalWith(libName, rawData, val)
					if err != nil {
						b.Fatal(err)
					}
				}
			})
		}
	}
}

func marshalWith(lib string, v any) ([]byte, error) {
	switch lib {
	case "encoding/json":
		return json.Marshal(v)
	case "json-iterator/go":
		return jsoniter.Marshal(v)
	case "bytedance/sonic":
		return sonic.Marshal(v)
	case "mailru/easyjson":
		return easyjson.Marshal(v.(easyjson.Marshaler))
	default:
		panic("unknown library: " + lib)
	}
}

func unmarshalWith(lib string, data []byte, v any) error {
	switch lib {
	case "encoding/json":
		return json.Unmarshal(data, v)
	case "json-iterator/go":
		return jsoniter.Unmarshal(data, v)
	case "bytedance/sonic":
		return sonic.Unmarshal(data, v)
	case "mailru/easyjson":
		return easyjson.Unmarshal(data, v.(easyjson.Unmarshaler))
	default:
		panic("unknown library: " + lib)
	}
}

var libs = []string{"encoding/json", "json-iterator/go", "bytedance/sonic", "mailru/easyjson"}

var testdata = []struct {
	name          string
	new           func() any
	rawDataGetter func() []byte
}{
	{
		name:          "canada_geometry",
		new:           func() any { return new(canadaRoot) },
		rawDataGetter: func() []byte { return gzip.UnzipFile("testdata/canada_geometry.json.gz") },
	},
	{
		name:          "citm_catalog",
		new:           func() any { return new(citmRoot) },
		rawDataGetter: func() []byte { return gzip.UnzipFile("testdata/citm_catalog.json.gz") },
	},
	{
		name:          "golang_source",
		new:           func() any { return new(golangRoot) },
		rawDataGetter: func() []byte { return gzip.UnzipFile("testdata/golang_source.json.gz") },
	},
	{
		name:          "string_unicode",
		new:           func() any { return new(stringRoot) },
		rawDataGetter: func() []byte { return gzip.UnzipFile("testdata/string_unicode.json.gz") },
	},
	{
		name:          "synthea_fhir",
		new:           func() any { return new(syntheaRoot) },
		rawDataGetter: func() []byte { return gzip.UnzipFile("testdata/synthea_fhir.json.gz") },
	},
	{
		name:          "twitter_status",
		new:           func() any { return new(twitterRoot) },
		rawDataGetter: func() []byte { return gzip.UnzipFile("testdata/twitter_status.json.gz") },
	},
	{
		name:          "nested_structure_1",
		new:           func() any { return new(NestedStructure) },
		rawDataGetter: func() []byte { return gzip.UnzipFile("testdata/nested_structure_1.json.gz") },
	},
	{
		name:          "nested_structure_2",
		new:           func() any { return new(NestedStructure) },
		rawDataGetter: func() []byte { return gzip.UnzipFile("testdata/nested_structure_2.json.gz") },
	},
	{
		name:          "nested_structure_3",
		new:           func() any { return new(NestedStructure) },
		rawDataGetter: func() []byte { return gzip.UnzipFile("testdata/nested_structure_3.json.gz") },
	},
	{
		name:          "nested_structure_4",
		new:           func() any { return new(NestedStructure) },
		rawDataGetter: func() []byte { return gzip.UnzipFile("testdata/nested_structure_4.json.gz") },
	},
	{
		name:          "nested_structure_5",
		new:           func() any { return new(NestedStructure) },
		rawDataGetter: func() []byte { return gzip.UnzipFile("testdata/nested_structure_5.json.gz") },
	},
	{
		name:          "nested_structure_6",
		new:           func() any { return new(NestedStructure) },
		rawDataGetter: func() []byte { return gzip.UnzipFile("testdata/nested_structure_6.json.gz") },
	},
	{
		name:          "nested_structure_7",
		new:           func() any { return new(NestedStructure) },
		rawDataGetter: func() []byte { return gzip.UnzipFile("testdata/nested_structure_7.json.gz") },
	},
	{
		name:          "nested_structure_8",
		new:           func() any { return new(NestedStructure) },
		rawDataGetter: func() []byte { return gzip.UnzipFile("testdata/nested_structure_8.json.gz") },
	},
	{
		name:          "nested_structure_9",
		new:           func() any { return new(NestedStructure) },
		rawDataGetter: func() []byte { return gzip.UnzipFile("testdata/nested_structure_9.json.gz") },
	},
	{
		name:          "nested_structure_10",
		new:           func() any { return new(NestedStructure) },
		rawDataGetter: func() []byte { return gzip.UnzipFile("testdata/nested_structure_10.json.gz") },
	},
	{
		name:          "number_structure_2",
		new:           func() any { return new(TestData2) },
		rawDataGetter: func() []byte { return gzip.UnzipFile("testdata/number_structure_2.json.gz") },
	},
	{
		name:          "number_structure_4",
		new:           func() any { return new(TestData4) },
		rawDataGetter: func() []byte { return gzip.UnzipFile("testdata/number_structure_4.json.gz") },
	},
	{
		name:          "number_structure_8",
		new:           func() any { return new(TestData8) },
		rawDataGetter: func() []byte { return gzip.UnzipFile("testdata/number_structure_8.json.gz") },
	},
	{
		name:          "number_structure_12",
		new:           func() any { return new(TestData12) },
		rawDataGetter: func() []byte { return gzip.UnzipFile("testdata/number_structure_12.json.gz") },
	},
}
