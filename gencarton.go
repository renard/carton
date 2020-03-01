//go:generate go test -tags gencartonForCartonPackageOnly

package carton

import (
	"bytes"
	"io"
	"os"
	"strings"

	"compress/gzip"
	"encoding/ascii85"

	"path/filepath"
	"text/template"
)

type cartonFile struct {
	Path    string
	Content string
	ModTime int64
}

type cartonSource struct {
	source      string
	destination string
	Package     string
	Name        string
	Files       []*cartonFile
}

var mycarton *cartonSource

// cartonTemplate is initialized with init() functions in both built-init.go
// when generating the carton.go file (when gencarton build flag IS set) or
// genbox-init.go (when gencarton build flag IS NOT set).
var cartonTemplate *template.Template

// New creates a new carton file
func New(pkg, name, path, out string) error {
	mycarton = &cartonSource{
		Package: pkg,
		Name:    name,
	}
	err := filepath.Walk(path, walker)

	os.Truncate(out, 0)
	fd, err := os.OpenFile(out, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer fd.Close()
	err = cartonTemplate.Execute(fd, mycarton)

	return err
}

func (b *cartonSource) addFile(f *cartonFile) {
	b.Files = append(b.Files, f)
}

// encode returns a string with path content gzipped and encoded to ascii85.
// To prevent from having ridiculously long lines the returned data is
// wrapped at 80 characters.
func encode(path string) (ret string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	var buf bytes.Buffer
	a85w := ascii85.NewEncoder(&buf)
	zw, err := gzip.NewWriterLevel(a85w, gzip.BestCompression)
	if err != nil {
		return
	}

	_, err = io.Copy(zw, f)
	if err != nil {
		return
	}
	zw.Close()
	a85w.Close()
	// backtick (`) is used for litteral strings by go. Here it is replaced
	// by tilde (~) to generate go literal string compatble data. Tilde is
	// not used by ascii85 so it is safe to do the replacement as long as it
	// is converted back to backtick before decoding the data.
	b := bytes.ReplaceAll(buf.Bytes(), []byte{'`'}, []byte("~"))

	c := chunk(b, 80)
	str := make([]string, len(c))
	for i, j := range c {
		str[i] = string(j)
	}
	ret = strings.Join(str, "\n")
	return
}

func walker(path string, info os.FileInfo, err error) (e error) {
	fi, e := os.Lstat(path)
	if e != nil || fi.IsDir() {
		return
	}

	str, e := encode(path)
	if e != nil {
		return
	}

	file := &cartonFile{
		Path:    path,
		Content: str,
		ModTime: fi.ModTime().UnixNano(),
	}

	mycarton.addFile(file)

	return nil
}

// Splits byte array in length-char lines.
//
// Returns an array of lines (each line is a byte array)
func chunk(data []byte, length int) [][]byte {
	// Compute initial parameters
	ldata := len(data)
	nChunks := ldata / length
	extraChunk := 0

	// In case of an incomplete last line
	if ldata > length*nChunks {
		extraChunk = 1
	}

	// Split data into chunks
	chunks := make([][]byte, nChunks+extraChunk)
	for i := 0; i < nChunks; i++ {
		chunks[i] = data[i*length : (i+1)*length]
	}

	// If need extraChunk, last line is incomplete. Read until end of data.
	if extraChunk > 0 {
		chunks[nChunks] = data[(nChunks)*length:]
	}

	return chunks
}
