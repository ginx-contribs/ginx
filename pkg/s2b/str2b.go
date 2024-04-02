package s2b

import "unsafe"

// StringToBytes return []byte representation of string, make sure underlying data will not be modified.
func StringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// BytesToString return string representation of []byte
func BytesToString(bs []byte) string {
	return unsafe.String(unsafe.SliceData(bs), len(bs))
}
