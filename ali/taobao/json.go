package taobao

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Time struct {
	Value time.Time
}

const timeLayout = "2006-01-02 15:04:05"

func (t *Time) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		t.Value = time.Time{}
	} else if tt, e := time.ParseInLocation(timeLayout, s, time.Local); e != nil {
		return e
	} else if tt.Year() == 1970 {
		t.Value = time.Time{}
	} else {
		t.Value = tt
	}
	return nil
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return []byte(t.Value.Format("\"" + timeLayout + "\"")), nil
}

type StringSlice []string
type jsonStringSlice struct {
	String []string `json:"string"`
}

func (ss *StringSlice) UnmarshalJSON(b []byte) (err error) {
	var jss jsonStringSlice
	if e := json.Unmarshal(b, &jss); e != nil {
		return e
	}
	*ss = StringSlice(jss.String)
	return nil
}

func (ss *StringSlice) MarshalJSON() ([]byte, error) {
	var jss jsonStringSlice
	jss.String = []string(*ss)
	return json.Marshal(&jss)
}

type CouponInfo struct {
	Spent float64
	Back  float64
}

func (ci *CouponInfo) UnmarshalJSON(b []byte) (err error) {
	s := string(b)
	if s == "\"null\"" {
		ci.Spent = 0
		ci.Back = 0
		return nil
	}
	_, e := fmt.Sscanf(s, `"满%g元减%g元"`, &ci.Spent, &ci.Back)
	return e
}

func (ci *CouponInfo) MarshalJSON() ([]byte, error) {
	return []byte("\"" + ci.String() + "\""), nil
}

func (ci *CouponInfo) String() string {
	if ci.Spent <= 0 || ci.Back <= 0 || ci.Back > ci.Spent {
		return ""
	}
	return fmt.Sprintf("满%g元减%g元", ci.Spent, ci.Back)
}
