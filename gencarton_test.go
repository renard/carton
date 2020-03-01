// +build !gencartonForCartonPackageOnly

package carton

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"
)

func TestGenerateCarton(t *testing.T) {
	t.Log("Generating a new carton.")
	err := New("carton", "Carton", "carton.tpl", "/dev/null")
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkGenerateCarton(b *testing.B) {
	b.Logf("benchmarking %d GenerateCarton loops", b.N)
	for i := 0; i < b.N; i++ {
		New("carton", "Carton", "carton.tpl", "/dev/null")
	}
}

func TestGetFile(t *testing.T) {
	t.Log("Reading file from carton.go.")
	b, err := carton.GetFile("carton.tpl")
	if err != nil {
		t.Error(err)
	}
	t.Log("Checking read size.")
	if len(b) == 0 {
		t.Error(errors.New("empty content for carton.tpl"))
	}
	t.Log("Compairing with local file.")
	fh, err := os.Open("carton.tpl")
	if err != nil {
		t.Error(err)
	}
	defer fh.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, fh)
	if err != nil {
		return
	}

	if !bytes.Equal(buf.Bytes(), b) {
		t.Error(errors.New("carton content and local file differ"))
	}
}

func BenchmarkGetFile(b *testing.B) {
	b.Logf("benchmark %d GetFile loops", b.N)
	for i := 0; i < b.N; i++ {
		carton.GetFile("Carton.tpl")
	}
}
