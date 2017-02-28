package gpsconv

import (
	"math"
	"testing"
)

func Test_DistanceFrom(t *testing.T) {
	loc1 := NewLocation(1, 1)
	loc2 := NewLocation(-1, -1)
	loc3 := NewLocation(116.3080027, 40.0359261)
	if math.Abs(loc1.DistanceFrom(loc2)-314851.074216824) > 0.1 {
		t.Fail()
	}
	if math.Abs(loc1.DistanceFrom(loc3)-12069547.1060307) > 0.1 {
		t.Fail()
	}
	if math.Abs(loc2.DistanceFrom(loc1)-314851.074216824) > 0.1 {
		t.Fail()
	}
}

func Test_Wgs84ToGcj02(t *testing.T) {
	origin := NewLocation(116.308016, 40.035937)
	result := origin.Wgs84ToGcj02()
	if math.Abs(result.Longitude-116.314122) > 0.00002 {
		t.Fail()
	}
	if math.Abs(result.Latitude-40.037216) > 0.0002 {
		t.Fail()
	}
}

func Test_Gcj02ToWgs84(t *testing.T) {
	origin := NewLocation(116.308016, 40.035937)
	result := origin.Gcj02ToWgs84()
	if math.Abs(result.Longitude-116.301910) > 0.00002 {
		t.Fail()
	}
	if math.Abs(result.Latitude-40.034658) > 0.0002 {
		t.Fail()
	}
}

func Test_Gcj02ToBaidu(t *testing.T) {
	origin := NewLocation(116.308016, 40.035937)
	result := origin.Gcj02ToBaidu()
	if math.Abs(result.Longitude-116.314490) > 0.00002 {
		t.Fail()
	}
	if math.Abs(result.Latitude-40.041968) > 0.0002 {
		t.Fail()
	}
}

func Test_BaiduToGcj02(t *testing.T) {
	origin := NewLocation(116.308016, 40.035937)
	result := origin.BaiduToGcj02()
	if math.Abs(result.Longitude-116.301577) > 0.00002 {
		t.Fail()
	}
	if math.Abs(result.Latitude-40.029790) > 0.0002 {
		t.Fail()
	}
}
