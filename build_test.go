// +build gencartonForCartonPackageOnly

package carton

import (
	"testing"
)

func TestGenBox(t *testing.T) {
	err := New("carton", "carton", "carton.tpl", "carton.go")
	if err != nil {
		panic(err)
	}
}
