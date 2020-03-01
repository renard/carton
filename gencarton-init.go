// +build !gencartonForCartonPackageOnly

package carton

import (
	"text/template"
)

// Solves chicken-and-egg problem.
//
// In regular operation carton package provides a carton.go with a cartonFS
// variable (carton). We force the use of the cartonFS over the local file
// template.
//
// This init function is reach when NOT generating carton.go when the
// gencarton build flag is not set.
func init() {
	var err error
	tpl, err := carton.GetFile("carton.tpl")
	if err != nil {
		panic(err)
	}
	cartonTemplate, err = template.New("Carton generator").Parse(string(tpl))
	if err != nil {
		panic(err)
	}
}
