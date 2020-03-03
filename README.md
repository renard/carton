carton
=====


[![Go Report Card][goreport-img]][goreport-url]
[![Build status][build-img]][build-url]
[![Coverage report][cover-img]][cover-url]
[![GoDoc][godoc-img]][godoc-url]
[![License: WTFPL][license-img]][license-url]

Carton is a *very* simple resource embedding system without the whistles and
bells for golang applications.

There are many similar packages so what distinguish this one from others?

* Small footprint: files are gzipped and ascii85 encoded (instead of base64).
* Use local file if more recent than embedded one (useful customization or
  testing).
* No need to use a specific binary, a simple *go generate* command is enough.

# Example


The `generate_test.go` used to generate the resources is:

```go
// +build gencarton

package main

import (
	"testing"
	"github.com/renard/carton"
)

func TestGenBox(t *testing.T) {
    // Generate a carton.go file with all files from the templates
    // directory. The carton.go is usable within the main package and
	// the variable containing all resources is CartonFiles.
	err := carton.New("main", "CartonFiles", "templates", "carton.go")
	if err != nil {
		panic(err)
	}
}
```

At the beginning of your `main.go` add following lines (do not forget the
empty line). It is important this line appears before the `package main`
line.

```
//go:generate go test -tags gencarton

```

The `gencarton` can be almost any valid golang build tag as long as it does
not clash with any other.

To build the resource simply run `go generate`.

If you get errors complaining about undefined variable:

```
./main.go:17:22: undefined: CartonFiles
```

Add a build directive in `main.go`:

```
//go:generate go test -tags gencarton
// +build !gencarton

package main

```



To use it in your application you can run something like:

```go
tpl, err := CartonFiles.GetFile("templates/mytemplate.txt")
if err != nil {
	panic(err)
}
mytemplate, err = template.New("MyTemplate").Parse(string(tpl))
if err != nil {
	panic(err)
}
// do something with mytemplate
err = myemplate.Execute(fd, myvars)
// ...
```

You can also use the `Files()` method to list all embedded files. If you
want to sort files within a tree structure you have to do it yourself.

# License

Copyright © 2020 Sébastien Gross <seb•ɑƬ•chezwam•ɖɵʈ•org> 

This program is free software. It comes without any warranty, to the extent
permitted by applicable law. You can redistribute it and/or modify it under
the terms of the Do What The Fuck You Want To Public License, Version 2, as
published by Sam Hocevar. See http://sam.zoy.org/wtfpl/COPYING for more
details.


[goreport-img]: https://goreportcard.com/badge/github.com/renard/carton
[goreport-url]: https://goreportcard.com/report/github.com/renard/carton
[build-img]: https://travis-ci.org/renard/carton.svg?branch=master
[build-url]: https://travis-ci.org/renard/carton
[cover-img]: https://coveralls.io/repos/github/renard/carton/badge.svg?branch=master
[cover-url]: https://coveralls.io/github/renard/carton
[godoc-img]: https://godoc.org/github.com/renard/carton?status.svg
[godoc-url]: https://godoc.org/github.com/renard/carton
[license-img]: https://img.shields.io/badge/License-WTFPL-brightgreen.svg
[license-url]: http://www.wtfpl.net/about/
