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
func NewLocation(longitude, latitude float64) *Location {
	return &Location{
		Longitude: longitude,
		Latitude:  latitude,
	}
}

// InsideChina returns true if the location is inside China, and false otherwise
func (this *Location) InsideChina() bool {
	if this.Longitude < 72.004 {
		return false
	}
	if this.Longitude > 137.8347 {
		return false
	}
	if this.Latitude < 0.8293 {
		return false
	}
	if this.Latitude > 55.8271 {
		return false
	}
	return true
}

// DistanceFrom calculates the distance between two location, the result is in meter
func (this *Location) DistanceFrom(loc *Location) float64 {
	rad1 := this.Latitude / 180.0 * math.Pi
	rad2 := loc.Latitude / 180.0 * math.Pi

	a := rad1 - rad2
	b := (this.Longitude - loc.Longitude) / 180.0 * math.Pi

	s := math.Pow(math.Sin(a/2), 2)
	s += math.Cos(rad1) * math.Cos(rad2) * math.Pow(math.Sin(b/2), 2)

	return 2 * math.Asin(math.Sqrt(s)) * 6378137.0
}

func (this *Location) transform() *Location {
	x := this.Longitude - 105.0
	y := this.Latitude - 35.0

	rad := this.Latitude / 180.0 * math.Pi
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

	return NewLocation(this.Longitude+dLng, this.Latitude+dLat)
}

// Wgs84ToGcj02 converts WGS84 location to GCJ02
func (this *Location) Wgs84ToGcj02() *Location {
	if this.InsideChina() {
		return this.transform()
	}
	return NewLocation(this.Longitude, this.Latitude)
}

// Gcj02ToWgs84 converts GCJ02 location to WGS84
func (this *Location) Gcj02ToWgs84() *Location {
	if !this.InsideChina() {
		return NewLocation(this.Longitude, this.Latitude)
	}
	loc := this.transform()
	return NewLocation(this.Longitude*2-loc.Longitude, this.Latitude*2-loc.Latitude)
}

// Gcj02ToBaidu converts GCJ02 location to BAIDU
func (this *Location) Gcj02ToBaidu() *Location {
	x := this.Longitude
	y := this.Latitude
	z := math.Sqrt(x*x+y*y) + 0.00002*math.Sin(y*gXPI)
	theta := math.Atan2(y, x) + 0.000003*math.Cos(x*gXPI)
	return NewLocation(z*math.Cos(theta)+0.0065, z*math.Sin(theta)+0.006)
}

// BaiduToGcj02 converts BAIDU location to GCJ02
func (this *Location) BaiduToGcj02() *Location {
	x := this.Longitude - 0.0065
	y := this.Latitude - 0.006
	z := math.Sqrt(x*x+y*y) - 0.00002*math.Sin(y*gXPI)
	theta := math.Atan2(y, x) - 0.000003*math.Cos(x*gXPI)
	return NewLocation(z*math.Cos(theta), z*math.Sin(theta))
}

// Wgs84ToBaidu converts WGS84 location to BAIDU
func (this *Location) Wgs84ToBaidu() *Location {
	return this.Wgs84ToGcj02().Gcj02ToBaidu()
}

// BaiduToWgs84 converts BAIDU location to WGS84
func (this *Location) BaiduToWgs84() *Location {
	return this.BaiduToGcj02().Gcj02ToWgs84()
}
