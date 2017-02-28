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

func (cfg Config) GetInt32(path string, dflt int32) int32 {
	path = strings.ToLower(path)
	if v, ok := cfg[path]; ok {
		if i, e := strconv.ParseInt(v, 10, 32); e == nil {
			return int32(i)
		}
	}
	return dflt
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

func (cfg Config) GetFloat(path string, dflt float64) float64 {
	path = strings.ToLower(path)
	if v, ok := cfg[path]; ok {
		if f, e := strconv.ParseFloat(v, 64); e == nil {
			return f
		}
	}
	return dflt
}

func (cfg Config) GetString(path string, dflt string) string {
	path = strings.ToLower(path)
	if v, ok := cfg[path]; ok {
		return v
	}
	return dflt
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

func GetInt32(path string, dflt int32) int32 {
	return Default.GetInt32(path, dflt)
}

func GetUint32(path string, dflt uint32) uint32 {
	return Default.GetUint32(path, dflt)
}

func GetFloat(path string, dflt float64) float64 {
	return Default.GetFloat(path, dflt)
}

func GetString(path string, dflt string) string {
	return Default.GetString(path, dflt)
}

func GetBool(path string, dflt bool) bool {
	return Default.GetBool(path, dflt)
}

func ForEach(prefix string, fx func(key, val string) bool) {
	Default.ForEach(prefix, fx)
}
