package gzip

import (
	"bytes"
	"compress/gzip"
	"io"
	"os"
)

func UnzipFile(fileName string) []byte {
	fl, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer fl.Close()

	gz, err := gzip.NewReader(fl)
	if err != nil {
		panic(err)
	}
	defer gz.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}
