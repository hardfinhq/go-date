# `go-date`

[![GoDoc][1]][2]
[![Go ReportCard][3]][4]

The `go-date` package provides a dedicated `Date{}` struct to emulate the
standard library `time.Time{}` behavior.

## API

This package provides helpers for:

- conversion: `ToTime()`, `date.FromTime()`, `date.FromString()`
- serialization: JSON and SQL
- emulating `time.Time{}`: `After()`, `Before()`, `Sub()`, etc.
- explicit null handling: `NullDate{}` and an analog of `sql.NullTime{}`
- emulating `time` helpers: `Today()` as an analog of `time.Now()`

## Background

The Go standard library contains no native type for dates without times.
Instead, common convention is to use a `time.Time{}` with only the year, month,
and day set. For example, this convention is followed when a timestamp of the
form YYYY-MM-DD is parsed via `time.Parse(time.DateOnly, s)`.

## Alternatives

This package is intended to be simple to understand and only needs to cover
"modern" dates (i.e. dates between 1900 and 2100). As a result, the core
`Date{}` struct directly exposes the year, month, and day as fields.

There are several alternative date packages which cover wider date ranges.
(These packages all use the [proleptic Gregorian calendar][6] to cover the
historical date ranges.) Some existing packages:

- `github.com/fxtlabs/date` [package][7]
- `github.com/rickb777/date` [package][5]

[1]: https://godoc.org/github.com/hardfinhq/go-date?status.svg
[2]: http://godoc.org/github.com/hardfinhq/go-date
[3]: https://goreportcard.com/badge/hardfinhq/go-date
[4]: https://goreportcard.com/report/hardfinhq/go-date
[5]: https://pkg.go.dev/github.com/rickb777/date
[6]: https://en.wikipedia.org/wiki/Proleptic_Gregorian_calendar
[7]: https://pkg.go.dev/github.com/fxtlabs/date
