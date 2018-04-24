package conv

import (
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

	s := strconv.FormatInt(int64(slice[0]), 10)
	if len(slice) == 1 {
		return s
	}

	var sb strings.Builder
	sb.WriteString(s)
	for _, i := range slice[1:] {
		sb.WriteString(sep)
		sb.WriteString(strconv.FormatInt(int64(i), 10))
	}
	return sb.String()
}

func SliceToStringUint32(slice []uint32, sep string) string {
	if len(slice) == 0 {
		return ""
	}

	s := strconv.FormatUint(uint64(slice[0]), 10)
	if len(slice) == 1 {
		return s
	}

	var sb strings.Builder
	sb.WriteString(s)
	for _, i := range slice[1:] {
		sb.WriteString(sep)
		sb.WriteString(strconv.FormatUint(uint64(i), 10))
	}
	return sb.String()
}

func SliceToStringUint16(slice []uint16, sep string) string {
	if len(slice) == 0 {
		return ""
	}

	s := strconv.FormatUint(uint64(slice[0]), 10)
	if len(slice) == 1 {
		return s
	}

	var sb strings.Builder
	sb.WriteString(s)
	for _, i := range slice[1:] {
		sb.WriteString(sep)
		sb.WriteString(strconv.FormatUint(uint64(i), 10))
	}
	return sb.String()
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
		v, _ := strconv.ParseUint(s, 10, 32)
		r[i] = uint32(v)
	}
	return r
}

func StringToSliceUint16(s string, seq string) []uint16 {
	ss := strings.Split(s, seq)
	r := make([]uint16, len(ss))
	for i, s := range ss {
		v, _ := strconv.ParseUint(s, 10, 16)
		r[i] = uint16(v)
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
