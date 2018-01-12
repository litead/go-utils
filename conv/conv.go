package conv

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func InetNtoa(ip uint32) string {
	sip := make([]byte, 4)
	binary.BigEndian.PutUint32(sip, ip)
	return fmt.Sprintf("%d.%d.%d.%d", sip[0], sip[1], sip[2], sip[3])
}

func InetAton(ip string) uint32 {
	if v := net.ParseIP(ip); v != nil {
		return binary.BigEndian.Uint32(v[12:])
	}
	return 0
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

func ItoaUint32(v uint32) string {
	return strconv.FormatUint(uint64(v), 10)
}

func ItoaInt32(v int32) string {
	return strconv.FormatInt(int64(v), 10)
}

func ItoaUint64(v uint64) string {
	return strconv.FormatUint(v, 10)
}

func Float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func AtoiUint32(s string) (uint32, error) {
	v, e := strconv.ParseUint(s, 10, 32)
	return uint32(v), e
}

func AtoiInt32(s string) (int32, error) {
	v, e := strconv.ParseInt(s, 10, 32)
	return int32(v), e
}

func AtoiInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
