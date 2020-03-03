// Package {{ .Package }} is a ambedded resources file
// generated with carton.
//
// See https://github.com/renard/carton
package {{ .Package }}

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"

	"compress/gzip"
	"encoding/ascii85"
	"io/ioutil"
)

// file stores all cartonFS Files information.
type file struct {
	// content is an ascii85 encoded string of the gzipped original file
	// content.
	content string
	// modTime is the original file last modification time in nanoseconds.
	modTime int64
}

// cartonFS is a map abstraction of the embedded resources. The map keys is
// the file path and the value is a *file struct.
type cartonFS map[string]*file

// Files returns a list of all files embedded in cartonFS
func (b *cartonFS) Files() []string {
	f := make([]string, len(*b))
	i := 0
	for k := range *b {
		f[i] = k
		i++
	}
	return f
}

// isLocalRecent compares the ModTime of local version and File
// file. Returns true only if local file is more recent than the cartonFS.
//
// If local file does not exists, false is returned.
func (b *cartonFS) isLocalRecent(path string, cartonfile *file) bool {
	fi, err := os.Lstat(path)
	if err != nil {
		return false
	}
	if fi.ModTime().UnixNano() > cartonfile.modTime {
		return true
	}
	return false
}

// getFileLocal returns the local file content intead of cartonFS embedded file.
func (b *cartonFS) getFileLocal(path string) (ret []byte, err error) {
	fh, err := os.Open(path)
	if err != nil {
		return
	}
	defer fh.Close()

	var buf bytes.Buffer

	_, err = io.Copy(&buf, fh)
	if err != nil {
		return
	}

	ret = buf.Bytes()
	return
}

// getFileLocal returns the cartonFS file content. However is a local file
// exists and is more recent this later's content is returned.
func (b *cartonFS) getFileFromCarton(path string) (ret []byte, err error) {
	f, ok := (*b)[path]
	if !ok {
		return nil, errors.New("File " + path + " not found.")
	}

	if b.isLocalRecent(path, f) {
		return b.getFileLocal(path)
	}

	// Replace back all tilde chars with backtick.
	decoder := ascii85.NewDecoder(
		strings.NewReader(
			strings.ReplaceAll(f.content, "~", "`")))
	gz, err := gzip.NewReader(decoder)
	if err != nil {
		return
	}
	ret, err = ioutil.ReadAll(gz)
	gz.Close()
	if err != nil {
		return
	}

	return
}

// GetFile return path files content. First try the cartonFS then local storage.
func (b *cartonFS) GetFile(path string) (ret []byte, err error) {
	ret, err = b.getFileFromCarton(path)
	if err != nil {
		ret, err = b.getFileLocal(path)
	}
	return
}

// Begin of dynamic content

// {{ .Name }} conatains the carton data.
var {{ .Name }} = &cartonFS{
{{ range .Files }}
	`{{ .Path }}`: &file{
		content: `
{{ .Content }}
`,
		modTime: {{ .ModTime }},
	},
{{ end -}}
}

// End of dynamic content
