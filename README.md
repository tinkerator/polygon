# polygon - a package to manipulate non self-intersecting polygons

## Overview

This [package](http://zappem.net/pub/math/polygon/) provides routines for 2D polygons. That is non
self-intersecting multi-edged shapes with more than 2 vertex points.

## API

The API provided by this package is avalble using `go doc
zappem.net/pub/math/polygon`. It can also be browsed on the
[go.dev](http://go.dev) website: [package
zappem.net/pub/math/polygon](https://pkg.go.dev/zappem.net/pub/math/polygon).

## Example

See the https://zappem.net/pub/project/polygons/ example.

## TODO

- A `2` in the
  [`zappem.net/pub/graphics/hershey`](https://zappem.net/pub/graphics/hershey/)
  `rowmand` font contains a dot of imperfection.

- Still working through some issues with concave polygons. Problem
  case: an E shape with a vertical line through the tines that do not
  extend out of the top or bottom of the shape.

## License info

The `polygon` package is distributed with the same BSD 3-clause
license as that used by [golang](https://golang.org/LICENSE) itself.

## Reporting bugs

Use the [github `polygon` bug
tracker](https://github.com/tinkerator/polygon/issues).
