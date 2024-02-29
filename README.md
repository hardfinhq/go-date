# `go-date`

[![GoDoc][1]][2]
[![Go ReportCard][3]][4]
[![Build Status][8]][9]

The `go-date` package provides a dedicated `Date{}` struct to emulate the
standard library `time.Time{}` behavior.

## API

This package provides helpers for:

- conversion: `ToTime()`, `date.FromTime()`, `date.FromString()`
- serialization: text, JSON, and SQL
- emulating `time.Time{}`: `After()`, `Before()`, `Sub()`, etc.
- explicit null handling: `NullDate{}` and an analog of `sql.NullTime{}`
- emulating `time` helpers: `Today()` as an analog of `time.Now()`

## Background

The Go standard library contains no native type for dates without times.
Instead, common convention is to use a `time.Time{}` with only the year, month,
and day set. For example, this convention is followed when a timestamp of the
form YYYY-MM-DD is parsed via `time.Parse(time.DateOnly, value)`.

## Conversion

For cases where existing code produces a "conventional"
`time.Date(YYYY, MM, DD, 0, 0, 0, 0, time.UTC)` value, it can be validated
and converted to a `Date{}` via:

```go
t := time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC)
d, err := date.FromTime(t)
fmt.Println(d, err)
// 2024-03-01 <nil>
```

If there is any deviation from the "conventional" format, this will error.
For example:

```text
timestamp contains more than just date information; 2020-05-11T01:00:00Z
timestamp contains more than just date information; 2022-01-31T00:00:00-05:00
```

For cases where we have a discrete timestamp (e.g. "last updated datetime") and
a relevant timezone for a given request, we can extract the date within that
timezone:

```go
t := time.Date(2023, time.April, 14, 3, 55, 4, 777000100, time.UTC)
tz, _ := time.LoadLocation("America/Chicago")
d := date.InTimezone(t, tz)
fmt.Println(d)
// 2023-04-13
```

For conversion in the **other** direction, a `Date{}` can be converted back
into a `time.Time{}`:

```go
d := date.NewDate(2017, time.July, 3)
t := d.ToTime()
fmt.Println(t)
// 2017-07-03 00:00:00 +0000 UTC
```

By default this will use the "conventional" format, but any of the values
(other than year, month, day) can also be set:

```go
d := date.NewDate(2017, time.July, 3)
tz, _ := time.LoadLocation("America/Chicago")
t := d.ToTime(date.OptConvertHour(12), date.OptConvertTimezone(tz))
fmt.Println(t)
// 2017-07-03 12:00:00 -0500 CDT
```

## Equivalent methods

There are a number of methods from `time.Time{}` that directly translate over:

```go
d := date.NewDate(2020, time.February, 29)
fmt.Println(d.Year)
// 2020
fmt.Println(d.Month)
// February
fmt.Println(d.Day)
// 29
fmt.Println(d.ISOWeek())
// 2020 9
fmt.Println(d.Weekday())
// Saturday

fmt.Println(d.IsZero())
// false
fmt.Println(d.String())
// 2020-02-29
fmt.Println(d.Format("Jan 2006"))
// Feb 2020
fmt.Println(d.GoString())
// date.NewDate(2020, time.February, 29)

d2 := date.NewDate(2021, time.February, 28)
fmt.Println(d2.Equal(d))
// false
fmt.Println(d2.Before(d))
// false
fmt.Println(d2.After(d))
// true
fmt.Println(d2.Compare(d))
// 1
```

However, some methods translate over only approximately. For example, it's much
more natural for `Sub()` to return the **number of days** between two dates:

```go
d := date.NewDate(2020, time.February, 29)
d2 := date.NewDate(2021, time.February, 28)
fmt.Println(d2.Sub(d))
// 365
```

## Divergent methods

We've elected to **translate** the `time.Time{}.AddDate()` method rather
than providing it directly:

```go
d := date.NewDate(2020, time.February, 29)
fmt.Println(d.AddDays(1))
// 2020-03-01
fmt.Println(d.AddDays(100))
// 2020-06-08
fmt.Println(d.AddMonths(1))
// 2020-03-29
fmt.Println(d.AddMonths(3))
// 2020-05-29
fmt.Println(d.AddYears(1))
// 2021-02-28
```

This is in part because of the behavior of the standard library's
`AddDate()`. In particular, it "overflows" a target month if the number
of days in that month is less than the number of desired days. As a result,
we provide `*Stdlib()` variants of the date addition helpers:

```go
d := date.NewDate(2020, time.February, 29)
fmt.Println(d.AddMonths(12))
// 2021-02-28
fmt.Println(d.AddMonthsStdlib(12))
// 2021-03-01
fmt.Println(d.AddYears(1))
// 2021-02-28
fmt.Println(d.AddYearsStdlib(1))
// 2021-03-01
```

In the same line of thinking as the divergent `AddMonths()` behavior, a
`MonthEnd()` method is provided that can pinpoint the number of days in
the current month:

```go
d := date.NewDate(2022, time.January, 14)
fmt.Println(d.MonthEnd())
// 2022-01-31
fmt.Println(d.MonthStart())
// 2022-01-01
```

## Integrating with `sqlc`

Out of the box, the `sqlc` [library][10] uses a Go `time.Time{}` both for
columns of type `TIMESTAMPTZ` and `DATE`. When reading `DATE` values (which come
over the wire in the form YYYY-MM-DD), the Go standard library produces values
of the form:

```go
time.Date(YYYY, MM, DD, 0, 0, 0, 0, time.UTC)
```

Instead, we can instruct `sqlc` to **globally** use `date.Date` and
`date.NullDate` when parsing `DATE` columns:

```yaml
---
version: '2'
overrides:
  go:
    overrides:
      - go_type:
          import: github.com/hardfinhq/go-date
          package: date
          type: NullDate
        db_type: date
        nullable: true
      - go_type:
          import: github.com/hardfinhq/go-date
          package: date
          type: Date
        db_type: date
        nullable: false
```

## Alternatives

This package is intended to be simple to understand and only needs to cover
"modern" dates (i.e. dates between 1900 and 2100). As a result, the core
`Date{}` struct directly exposes the year, month, and day as fields.

There are several alternative date packages which cover wider date ranges.
(These packages all use the [proleptic Gregorian calendar][6] to cover the
historical date ranges.) Some existing packages:

- `github.com/fxtlabs/date` [package][7]
- `github.com/rickb777/date` [package][5]

Additionally, there is a `Date{}` type provided by the `github.com/jackc/pgtype`
[package][11] that is part of the `pgx` ecosystem. However, this type is very
focused on being useful for database serialization and deserialization and
doesn't implement a wider set of methods present on `time.Time{}` (e.g.
`After()`).

[1]: https://godoc.org/github.com/hardfinhq/go-date?status.svg
[2]: http://godoc.org/github.com/hardfinhq/go-date
[3]: https://goreportcard.com/badge/hardfinhq/go-date
[4]: https://goreportcard.com/report/hardfinhq/go-date
[5]: https://pkg.go.dev/github.com/rickb777/date
[6]: https://en.wikipedia.org/wiki/Proleptic_Gregorian_calendar
[7]: https://pkg.go.dev/github.com/fxtlabs/date
[8]: https://github.com/hardfinhq/go-date/actions/workflows/ci.yaml/badge.svg?branch=main
[9]: https://github.com/hardfinhq/go-date/actions/workflows/ci.yaml
[10]: https://docs.sqlc.dev
[11]: https://pkg.go.dev/github.com/jackc/pgtype
