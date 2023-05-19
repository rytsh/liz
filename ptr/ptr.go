/*
Converting values to pointers.
*/
package ptr

import "time"

func Bool(b bool) *bool {
	return &b
}

func Int(i int) *int {
	return &i
}

func Int8(i int8) *int8 {
	return &i
}

func Int16(i int16) *int16 {
	return &i
}

func Int32(i int32) *int32 {
	return &i
}

func Int64(i int64) *int64 {
	return &i
}

func Uint(i uint) *uint {
	return &i
}

func Uint8(i uint8) *uint8 {
	return &i
}

func Uint16(i uint16) *uint16 {
	return &i
}

func Uint32(i uint32) *uint32 {
	return &i
}

func Uint64(i uint64) *uint64 {
	return &i
}

func Float32(f float32) *float32 {
	return &f
}

func Float64(f float64) *float64 {
	return &f
}

func Complex64(c complex64) *complex64 {
	return &c
}

func Complex128(c complex128) *complex128 {
	return &c
}

func String(s string) *string {
	return &s
}

func Rune(r rune) *rune {
	return &r
}

func Byte(b byte) *byte {
	return &b
}

func Interface(i interface{}) *interface{} {
	return &i
}

func Map(m map[string]interface{}) *map[string]interface{} {
	return &m
}

func Slice(s []interface{}) *[]interface{} {
	return &s
}

func Time(t time.Time) *time.Time {
	return &t
}
