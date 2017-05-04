package config

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strconv"
	"strings"
)

type Config map[string]string

func removeUtf8Bom(data []byte) []byte {
	if len(data) > 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		data = data[3:]
	}
	return data
}

func (cfg Config) ParseIniStream(reader io.Reader) error {
	section, lastKey := "/", ""
	firstLine, scanner := true, bufio.NewScanner(reader)

	for scanner.Scan() {
		s := scanner.Bytes()
		if firstLine {
			s = removeUtf8Bom(s)
			firstLine = false
		}

		s = bytes.TrimSpace(s)
		if len(s) == 0 || s[0] == '#' { // empty or comment
			continue
		}

		if s[0] == '[' && s[len(s)-1] == ']' { // section
			s = bytes.TrimSpace(s[1 : len(s)-1])
			if len(s) >= 0 {
				section = "/" + string(bytes.ToLower(s))
			}
			continue
		}

		k, v := "", ""
		if i := bytes.IndexByte(s, '='); i != -1 {
			k = string(bytes.ToLower(bytes.TrimSpace(s[:i])))
			v = string(bytes.TrimSpace(s[i+1:]))
		}

		if len(k) > 0 {
			lastKey = section + "/" + k
			cfg[lastKey] = v
			continue
		} else if len(lastKey) == 0 {
			continue
		}

		c, lv := byte(128), cfg[lastKey]
		if len(lv) > 0 {
			c = lv[len(lv)-1]
		}

		if len(v) == 0 { // empty value means a new line
			cfg[lastKey] = lv + "\n"
		} else if c < 128 && c != '-' && v[0] < 128 { // need a white space?
			// not good enough, but should be ok in most cases
			cfg[lastKey] = lv + " " + v
		} else {
			cfg[lastKey] = lv + v
		}
	}

	if e := scanner.Err(); e != nil {
		return e
	}

	return nil
}

func (cfg Config) ParseIniFile(path string) error {
	f, e := os.Open(path)
	if e != nil {
		return e
	}
	defer f.Close()

	return cfg.ParseIniStream(f)
}

func (cfg Config) GetInt(path string, dflt int) int {
	path = strings.ToLower(path)
	if v, ok := cfg[path]; ok {
		if i, e := strconv.Atoi(v); e == nil {
			return i
		}
	}
	return dflt
}

func (cfg Config) Int(path string) int {
	v := cfg.String(path)
	i, e := strconv.Atoi(v)
	if e != nil {
		panic(e)
	}
	return i
}

func (cfg Config) GetInt64(path string, dflt int64) int64 {
	path = strings.ToLower(path)
	if v, ok := cfg[path]; ok {
		if i, e := strconv.ParseInt(v, 10, 64); e == nil {
			return i
		}
	}
	return dflt
}

func (cfg Config) Int64(path string) int64 {
	v := cfg.String(path)
	i, e := strconv.ParseInt(v, 10, 64)
	if e != nil {
		panic(e)
	}
	return i
}

func (cfg Config) GetUint64(path string, dflt uint64) uint64 {
	path = strings.ToLower(path)
	if v, ok := cfg[path]; ok {
		if u, e := strconv.ParseUint(v, 10, 64); e == nil {
			return u
		}
	}
	return dflt
}

func (cfg Config) Uint64(path string) uint64 {
	v := cfg.String(path)
	u, e := strconv.ParseUint(v, 10, 64)
	if e != nil {
		panic(e)
	}
	return u
}

func (cfg Config) GetInt32(path string, dflt int32) int32 {
	path = strings.ToLower(path)
	if v, ok := cfg[path]; ok {
		if i, e := strconv.ParseInt(v, 10, 32); e == nil {
			return int32(i)
		}
	}
	return dflt
}

func (cfg Config) Int32(path string) int32 {
	v := cfg.String(path)
	i, e := strconv.ParseInt(v, 10, 32)
	if e != nil {
		panic(e)
	}
	return int32(i)
}

func (cfg Config) GetUint32(path string, dflt uint32) uint32 {
	path = strings.ToLower(path)
	if v, ok := cfg[path]; ok {
		if u, e := strconv.ParseUint(v, 10, 32); e == nil {
			return uint32(u)
		}
	}
	return dflt
}

func (cfg Config) Uint32(path string) uint32 {
	v := cfg.String(path)
	u, e := strconv.ParseUint(v, 10, 32)
	if e != nil {
		panic(e)
	}
	return uint32(u)
}

func (cfg Config) GetFloat64(path string, dflt float64) float64 {
	path = strings.ToLower(path)
	if v, ok := cfg[path]; ok {
		if f, e := strconv.ParseFloat(v, 64); e == nil {
			return f
		}
	}
	return dflt
}

func (cfg Config) Float64(path string) float64 {
	v := cfg.String(path)
	f, e := strconv.ParseFloat(v, 64)
	if e != nil {
		panic(e)
	}
	return f
}

func (cfg Config) GetFloat32(path string, dflt float32) float32 {
	path = strings.ToLower(path)
	if v, ok := cfg[path]; ok {
		if f, e := strconv.ParseFloat(v, 32); e == nil {
			return float32(f)
		}
	}
	return dflt
}

func (cfg Config) Float32(path string) float32 {
	v := cfg.String(path)
	f, e := strconv.ParseFloat(v, 32)
	if e != nil {
		panic(e)
	}
	return float32(f)
}

func (cfg Config) GetString(path string, dflt string) string {
	path = strings.ToLower(path)
	if v, ok := cfg[path]; ok {
		return v
	}
	return dflt
}

func (cfg Config) String(path string) string {
	path = strings.ToLower(path)
	if v, ok := cfg[path]; ok {
		return v
	}
	panic("path not found")
}

func (cfg Config) GetBool(path string, dflt bool) bool {
	path = strings.ToLower(path)
	if v, ok := cfg[path]; ok {
		if b, e := strconv.ParseBool(v); e == nil {
			return b
		}
	}
	return dflt
}

func (cfg Config) Bool(path string) bool {
	v := cfg.String(path)
	b, e := strconv.ParseBool(v)
	if e != nil {
		panic(e)
	}
	return b
}

func (cfg Config) ForEach(prefix string, fx func(key, val string) bool) {
	prefix = strings.ToLower(prefix)
	for k, v := range cfg {
		if strings.HasPrefix(k, prefix) && fx(k, v) {
			break
		}
	}
}

func New() Config {
	return Config(make(map[string]string))
}

var Default = New()

func Reset() {
	Default = New()
}

func ParseIniStream(reader io.Reader) error {
	return Default.ParseIniStream(reader)
}

func ParseIniFile(path string) error {
	return Default.ParseIniFile(path)
}

func GetInt(path string, dflt int) int {
	return Default.GetInt(path, dflt)
}

func Int(path string) int {
	return Default.Int(path)
}

func GetInt64(path string, dflt int64) int64 {
	return Default.GetInt64(path, dflt)
}

func Int64(path string) int64 {
	return Default.Int64(path)
}

func GetUint64(path string, dflt uint64) uint64 {
	return Default.GetUint64(path, dflt)
}

func Uint64(path string) uint64 {
	return Default.Uint64(path)
}

func GetInt32(path string, dflt int32) int32 {
	return Default.GetInt32(path, dflt)
}

func Int32(path string) int32 {
	return Default.Int32(path)
}

func GetUint32(path string, dflt uint32) uint32 {
	return Default.GetUint32(path, dflt)
}

func Uint32(path string) uint32 {
	return Default.Uint32(path)
}

func GetFloat64(path string, dflt float64) float64 {
	return Default.GetFloat64(path, dflt)
}

func Float64(path string) float64 {
	return Default.Float64(path)
}

func GetFloat32(path string, dflt float32) float32 {
	return Default.GetFloat32(path, dflt)
}

func Float32(path string) float32 {
	return Default.Float32(path)
}

func GetString(path string, dflt string) string {
	return Default.GetString(path, dflt)
}

func String(path string) string {
	return Default.String(path)
}

func GetBool(path string, dflt bool) bool {
	return Default.GetBool(path, dflt)
}

func Bool(path string) bool {
	return Default.Bool(path)
}

func ForEach(prefix string, fx func(key, val string) bool) {
	Default.ForEach(prefix, fx)
}
