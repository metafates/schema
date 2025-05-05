package constraint

import "time"

type Float interface {
	~float32 | ~float64
}

type Integer interface {
	Signed | Unsigned
}

type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Real interface{ Float | Integer }

// Text constraints types that can be converted to string
type Text interface{ ~string | ~[]rune | ~[]byte }

type Comparable[T any] interface{ Compare(other T) int }

type Time = Comparable[time.Time]
