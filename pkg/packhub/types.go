package packhub

import "time"

type Numbers interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64
}

type Generic interface {
	~bool | ~[]bool | ~string | ~[]string | Numbers | time.Time
}
