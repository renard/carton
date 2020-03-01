// +build gencartonForCartonPackageOnly

package carton

import (
	"text/template"
)

// Solves chicken-and-egg problem.
//
// When carton.go does not exist the carton-generating template is not
// available as a cartonFS. Here we force to load it from the original
// template file.
//
// This init function is reach when generating carton.go using "go generate"
// function (it uses the gencarton build flag).
func init() {
	var err error
	cartonTemplate, err = template.ParseFiles("carton.tpl")
	if err != nil {
		panic(err)
	}
}
