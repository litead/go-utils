package gpsconv

import "math"

//
// Krasovsky 1940
//
// gA = 6378245.0, 1/f = 298.3
// gB = gA * (1 - f)
// gEE = (gA^2 - gB^2) / gA^2;
const gA = 6378245.0
const gEE = 0.00669342162296594323
const gXPI = math.Pi * 3000.0 / 180.0

// Location is GPS location with Longitude & Latitude
type Location struct {
	Longitude float64
	Latitude  float64
}

// NewLocation creates a new Location object
func NewLocation(longitude, latitude float64) Location {
	return Location{
		Longitude: longitude,
		Latitude:  latitude,
	}
}

// InsideChina returns true if the location is inside China, and false otherwise
func InsideChina(lng, lat float64) bool {
	if lng < 72.004 {
		return false
	}
	if lng > 137.8347 {
		return false
	}
	if lat < 0.8293 {
		return false
	}
	if lat > 55.8271 {
		return false
	}
	return true
}

func (l Location) InsideChina() bool {
	return InsideChina(l.Longitude, l.Latitude)
}

// DistanceFrom calculates the distance between two location, the result is in meter
func DistanceFrom(lng1, lat1, lng2, lat2 float64) float64 {
	rad1 := lat1 / 180.0 * math.Pi
	rad2 := lat2 / 180.0 * math.Pi

	a := rad1 - rad2
	b := (lng1 - lng2) / 180.0 * math.Pi

	s := math.Pow(math.Sin(a/2), 2)
	s += math.Cos(rad1) * math.Cos(rad2) * math.Pow(math.Sin(b/2), 2)

	return 2 * math.Asin(math.Sqrt(s)) * 6378137.0
}

func (l Location) DistanceFrom(l2 Location) float64 {
	return DistanceFrom(l.Longitude, l.Latitude, l2.Longitude, l2.Latitude)
}

func transform(lng, lat float64) (float64, float64) {
	x := lng - 105.0
	y := lat - 35.0

	rad := lat / 180.0 * math.Pi
	sin := math.Sin(rad)
	cos := math.Cos(rad)
	magic := 1 - gEE*sin*sin
	sqrt := math.Sqrt(magic)

	dLng := 20.0*math.Sin(6.0*x*math.Pi) + 20.0*math.Sin(2.0*x*math.Pi)
	dLat := dLng

	dLng += 20.0*math.Sin(x*math.Pi) + 40.0*math.Sin(x/3.0*math.Pi)
	dLng += 150.0*math.Sin(x/12.0*math.Pi) + 300.0*math.Sin(x/30.0*math.Pi)
	dLng = dLng * 2.0 / 3.0
	dLng += 300.0 + x + 2.0*y + 0.1*x*x + 0.1*x*y + 0.1*math.Sqrt(math.Abs(x))
	dLng = (dLng * 180.0) / (gA / sqrt * cos * math.Pi)

	dLat += 20.0*math.Sin(y*math.Pi) + 40.0*math.Sin(y/3.0*math.Pi)
	dLat += 160.0*math.Sin(y/12.0*math.Pi) + 320*math.Sin(y*math.Pi/30.0)
	dLat = dLat * 2.0 / 3.0
	dLat += -100.0 + 2.0*x + 3.0*y + 0.2*y*y + 0.1*x*y + 0.2*math.Sqrt(math.Abs(x))
	dLat = (dLat * 180.0) / ((gA * (1 - gEE)) / (magic * sqrt) * math.Pi)

	return lng + dLng, lat + dLat
}

// Wgs84ToGcj02 converts WGS84 location to GCJ02
func Wgs84ToGcj02(lng, lat float64) (float64, float64) {
	return transform(lng, lat)
}

func (l Location) Wgs84ToGcj02() Location {
	return NewLocation(transform(l.Longitude, l.Latitude))
}

// Gcj02ToWgs84 converts GCJ02 location to WGS84
func Gcj02ToWgs84(lng, lat float64) (float64, float64) {
	lng1, lat1 := transform(lng, lat)
	return lng*2 - lng1, lat*2 - lat1
}

func (l Location) Gcj02ToWgs84() Location {
	return NewLocation(Gcj02ToWgs84(l.Longitude, l.Latitude))
}

// Gcj02ToBaidu converts GCJ02 location to BAIDU
func Gcj02ToBaidu(lng, lat float64) (float64, float64) {
	z := math.Sqrt(lng*lng+lat*lat) + 0.00002*math.Sin(lat*gXPI)
	theta := math.Atan2(lat, lng) + 0.000003*math.Cos(lng*gXPI)
	return z*math.Cos(theta) + 0.0065, z*math.Sin(theta) + 0.006
}

func (l Location) Gcj02ToBaidu() Location {
	return NewLocation(Gcj02ToBaidu(l.Longitude, l.Latitude))
}

// BaiduToGcj02 converts BAIDU location to GCJ02
func BaiduToGcj02(lng, lat float64) (float64, float64) {
	x := lng - 0.0065
	y := lat - 0.006
	z := math.Sqrt(x*x+y*y) - 0.00002*math.Sin(y*gXPI)
	theta := math.Atan2(y, x) - 0.000003*math.Cos(x*gXPI)
	return z * math.Cos(theta), z * math.Sin(theta)
}

func (l Location) BaiduToGcj02() Location {
	return NewLocation(BaiduToGcj02(l.Longitude, l.Latitude))
}

// Wgs84ToBaidu converts WGS84 location to BAIDU
func Wgs84ToBaidu(lng, lat float64) (float64, float64) {
	return Gcj02ToBaidu(Wgs84ToGcj02(lng, lat))
}

func (l Location) Wgs84ToBaidu() Location {
	return NewLocation(Wgs84ToBaidu(l.Longitude, l.Latitude))
}

// BaiduToWgs84 converts BAIDU location to WGS84
func BaiduToWgs84(lng, lat float64) (float64, float64) {
	return Gcj02ToWgs84(BaiduToGcj02(lng, lat))
}

func (l Location) BaiduToWgs84() Location {
	return NewLocation(BaiduToWgs84(l.Longitude, l.Latitude))
}
