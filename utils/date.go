package utils

import (
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// StrToIntMonth 字符串月份转整数月份
func StrToIntMonth(month string) int {
	var data = map[string]int{
		"January":   0,
		"February":  1,
		"March":     2,
		"April":     3,
		"May":       4,
		"June":      5,
		"July":      6,
		"August":    7,
		"September": 8,
		"October":   9,
		"November":  10,
		"December":  11,
	}
	return data[month]
}

// GetTodayYMD 得到以sep为分隔符的年、月、日字符串(今天)
func GetTodayYMD(sep string) string {
	now := time.Now()
	year := now.Year()
	month := StrToIntMonth(now.Month().String())
	date := now.Day()

	var monthStr string
	var dateStr string
	if month < 9 {
		monthStr = "0" + strconv.Itoa(month+1)
	} else {
		monthStr = strconv.Itoa(month + 1)
	}

	if date < 10 {
		dateStr = "0" + strconv.Itoa(date)
	} else {
		dateStr = strconv.Itoa(date)
	}
	return strconv.Itoa(year) + sep + monthStr + sep + dateStr
}

// GetTodayYM 得到以sep为分隔符的年、月字符串(今天所属于的月份)
func GetTodayYM(sep string) string {
	now := time.Now()
	year := now.Year()
	month := StrToIntMonth(now.Month().String())

	var monthStr string
	if month < 9 {
		monthStr = "0" + strconv.Itoa(month+1)
	} else {
		monthStr = strconv.Itoa(month + 1)
	}
	return strconv.Itoa(year) + sep + monthStr
}

// GetYesterdayYMD 得到以sep为分隔符的年、月、日字符串(昨天)
func GetYesterdayYMD(sep string) string {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	todaySec := today.Unix()            //秒
	yesterdaySec := todaySec - 24*60*60 //秒
	yesterdayTime := time.Unix(yesterdaySec, 0)
	yesterdayYMD := yesterdayTime.Format("2006-01-02")
	return strings.Replace(yesterdayYMD, "-", sep, -1)
}

// GetTomorrowYMD 得到以sep为分隔符的年、月、日字符串(明天)
func GetTomorrowYMD(sep string) string {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	todaySec := today.Unix()           //秒
	tomorrowSec := todaySec + 24*60*60 //秒
	tomorrowTime := time.Unix(tomorrowSec, 0)
	tomorrowYMD := tomorrowTime.Format("2006-01-02")
	return strings.Replace(tomorrowYMD, "-", sep, -1)
}

// GetTodayTime 返回今天零点的time
func GetTodayTime() time.Time {
	now := time.Now()
	// now.Year(), now.Month(), now.Day() 是以本地时区为参照的年、月、日
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	return today
}

// GetYesterdayTime 返回昨天零点的time
func GetYesterdayTime() time.Time {
	now := time.Now()
	// now.Year(), now.Month(), now.Day() 是以本地时区为参照的年、月、日
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	yesterdaySec := today.Unix() - 24*60*60
	return time.Unix(yesterdaySec, 0)
}

func isEmpty(a interface{}) bool {
	v := reflect.ValueOf(a)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v.Interface() == reflect.Zero(v.Type()).Interface()
}

func Format(time *time.Time, format string) string {
	if time == nil {
		return ""
	}
	return time.Format(format)
}

func FormatS(time time.Time, format string) string {
	if time.IsZero() {
		return ""
	}
	return time.Format(format)
}

func Split(val string, sep string) []string {
	var arr []string
	val = strings.Trim(val, "")
	if len(val) > 0 {
		return strings.Split(val, sep)
	} else {
		return arr
	}
}

func shiftJulianToNoon(julianDays, julianFraction float64) (float64, float64) {
	switch {
	case -0.5 < julianFraction && julianFraction < 0.5:
		julianFraction += 0.5
	case julianFraction >= 0.5:
		julianDays++
		julianFraction -= 0.5
	case julianFraction <= -0.5:
		julianDays--
		julianFraction += 1.5
	}
	return julianDays, julianFraction
}

func julianDateToGregorianTime(part1, part2 float64) time.Time {
	part1I, part1F := math.Modf(part1)
	part2I, part2F := math.Modf(part2)
	julianDays := part1I + part2I
	julianFraction := part1F + part2F
	julianDays, julianFraction = shiftJulianToNoon(julianDays, julianFraction)
	day, month, year := doTheFliegelAndVanFlandernAlgorithm(int(julianDays))
	hours, minutes, seconds, nanoseconds := fractionOfADay(julianFraction)
	return time.Date(year, time.Month(month), day, hours, minutes, seconds, nanoseconds, time.UTC)
}

func fractionOfADay(fraction float64) (hours, minutes, seconds, nanoseconds int) {

	const (
		c1us  = 1e3
		c1s   = 1e9
		c1day = 24 * 60 * 60 * c1s
	)

	frac := int64(c1day*fraction + c1us/2)
	nanoseconds = int((frac%c1s)/c1us) * c1us
	frac /= c1s
	seconds = int(frac % 60)
	frac /= 60
	minutes = int(frac % 60)
	hours = int(frac / 60)
	return
}

func doTheFliegelAndVanFlandernAlgorithm(jd int) (day, month, year int) {
	l := jd + 68569
	n := (4 * l) / 146097
	l = l - (146097*n+3)/4
	i := (4000 * (l + 1)) / 1461001
	l = l - (1461*i)/4 + 31
	j := (80 * l) / 2447
	d := l - (2447*j)/80
	l = j / 11
	m := j + 2 - (12 * l)
	y := 100*(n-49) + i + l
	return d, m, y
}

func TimeFromExcelTime(excelTime float64, date1904 bool) time.Time {
	const MDD int64 = 106750 // Max time.Duration Days, aprox. 290 years
	var date time.Time
	var intPart = int64(excelTime)
	// Excel uses Julian dates prior to March 1st 1900, and Gregorian
	// thereafter.
	if intPart <= 61 {
		const OFFSET1900 = 15018.0
		const OFFSET1904 = 16480.0
		const MJD0 float64 = 2400000.5
		var date time.Time
		if date1904 {
			date = julianDateToGregorianTime(MJD0, excelTime+OFFSET1904)
		} else {
			date = julianDateToGregorianTime(MJD0, excelTime+OFFSET1900)
		}
		return date
	}
	var floatPart = excelTime - float64(intPart)
	var dayNanoSeconds float64 = 24 * 60 * 60 * 1000 * 1000 * 1000
	if date1904 {
		date = time.Date(1904, 1, 1, 0, 0, 0, 0, time.UTC)
	} else {
		date = time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
	}

	// Duration is limited to aprox. 290 years
	for intPart > MDD {
		durationDays := time.Duration(MDD) * time.Hour * 24
		date = date.Add(durationDays)
		intPart = intPart - MDD
	}
	durationDays := time.Duration(intPart) * time.Hour * 24
	durationPart := time.Duration(dayNanoSeconds * floatPart)
	return date.Add(durationDays).Add(durationPart)
}
