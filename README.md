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

Here are some known issues potentially gating release of v1.0.0.

- Observations from [`zappem.net/pub/graphics/hershey`](https://zappem.net/pub/graphics/hershey/) font rendering:
  - `rowmant` and `timesrb` font renders `a` with a split top left corner.
  - sometimes, these fonts look like they render `@` with the same
    split at the center right of the low part of the `a` inside the
    `@` symbol. This is really hard to reproduce.

Some things to look into after v1.0.0.

- Performance

## License info

The `polygon` package is distributed with the same BSD 3-clause
license as that used by [golang](https://golang.org/LICENSE) itself.

## Reporting bugs

Use the [github `polygon` bug
tracker](https://github.com/tinkerator/polygon/issues).
