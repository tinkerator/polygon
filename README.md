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

- Still working through some issues with concave polygons. Problem
  case: TestConcaveE() currently configured to be skipped.

- Observations from [`zappem.net/pub/graphics/hershey`](https://zappem.net/pub/graphics/hershey/) font rendering:
  - A `2` in the `rowmand` font contains a dot of imperfection. A
    potential issue with `8`, but it renders occasionally, so the
    issue with that may not be the polygon code.
  - `T` and `5` in the `astrology` font have the same dot issue.
  - `J` in the `rowmant` is worth taking a closer look at. `T` renders
    with a stray dot. Also, there may be an issue with `t`.
  - `symbolic` and `timesg` fail to render and crash with my test program.
  - `T` and `Z` have extra dots in `timesib`, `timesr` and `timesrb` fonts.

## License info

The `polygon` package is distributed with the same BSD 3-clause
license as that used by [golang](https://golang.org/LICENSE) itself.

## Reporting bugs

Use the [github `polygon` bug
tracker](https://github.com/tinkerator/polygon/issues).
