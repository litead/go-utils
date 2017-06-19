package conv

import (
	"bytes"
	"net"
	"strconv"
	"strings"
)

func InetNtoa(n uint32) string {
	fields := []string{
		strconv.Itoa(int(byte(n >> 24))),
		strconv.Itoa(int(byte(n >> 16))),
		strconv.Itoa(int(byte(n >> 8))),
		strconv.Itoa(int(byte(n))),
	}
	return strings.Join(fields, ".")
}

func InetAton(s string) uint32 {
	ip := net.ParseIP(s)
	if ip == nil {
		return 0
	}
	return uint32(ip[12])<<24 + uint32(ip[13])<<16 + uint32(ip[14])<<8 + uint32(ip[15])
}

func SliceToStringInt(slice []int, sep string) string {
	if len(slice) == 0 {
		return ""
	}

	s := strconv.Itoa(slice[0])
	if len(slice) == 1 {
		return s
	}

	var buf bytes.Buffer
	buf.WriteString(s)
	bs := []byte(sep)
	for _, i := range slice[1:] {
		buf.Write(bs)
		buf.WriteString(strconv.Itoa(i))
	}
	return buf.String()
}

func SliceToStringUint32(slice []uint32, sep string) string {
	if len(slice) == 0 {
		return ""
	}

	s := strconv.Itoa(int(slice[0]))
	if len(slice) == 1 {
		return s
	}

	var buf bytes.Buffer
	buf.WriteString(s)
	bs := []byte(sep)
	for _, i := range slice[1:] {
		buf.Write(bs)
		buf.WriteString(strconv.Itoa(int(i)))
	}
	return buf.String()
}

func StringToSliceInt(s string, seq string) []int {
	ss := strings.Split(s, seq)
	r := make([]int, len(ss))
	for i, s := range ss {
		v, _ := strconv.Atoi(s)
		r[i] = v
	}
	return r
}

func StringToSliceUint32(s string, seq string) []uint32 {
	ss := strings.Split(s, seq)
	r := make([]uint32, len(ss))
	for i, s := range ss {
		v, _ := strconv.Atoi(s)
		r[i] = uint32(v)
	}
	return r
}
