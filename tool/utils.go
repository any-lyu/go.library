package tool

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"regexp"
	"time"
	"unsafe"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// D2S struct to string
func D2S(data interface{}) (r string) {
	bs, _ := json.Marshal(data)
	return string(bs)
}

// D2B struct to bytes
func D2B(data interface{}) []byte {
	bs, _ := json.Marshal(data)
	return bs
}

//B2D bytes to struct
func B2D(b []byte, data interface{}) error {
	return json.Unmarshal(b, data)
}

//B2S bytes to string
func B2S(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// S2B string to bytes
func S2B(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

// B2I bytes to int
func B2I(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	_ = binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int(x)
}

// I2B int to bytes
func I2B(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

// Md5 cz
func Md5(msg string) string {
	h := md5.New()
	h.Write([]byte(msg))
	return hex.EncodeToString(h.Sum(nil))
}

//TimeSubDay time sub time return days
func TimeSubDay(t1, t2 time.Time) int {
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, time.Local)
	t2 = time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, time.Local)
	return int(t1.Sub(t2).Hours() / 24)
}

// TimeSubHour time sub time return hours
func TimeSubHour(t1, t2 time.Time) int {
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour(), 0, 0, 0, time.Local)
	t2 = time.Date(t2.Year(), t2.Month(), t2.Day(), t2.Hour(), 0, 0, 0, time.Local)
	return int(t1.Sub(t2).Hours())
}

// TimeSubWeek time sub time return weeks
func TimeSubWeek(t1, t2 time.Time, switchDay time.Weekday) int {
	weekDay := t2.Weekday() - switchDay
	if weekDay > 0 {
		t2 = t2.Add(-time.Hour * 24 * time.Duration(weekDay))
	}
	t2 = time.Date(t2.Year(), t2.Month(), t2.Day(), t2.Hour(), 0, 0, 0, time.Local)
	return TimeSubDay(t1, t2) / 7
}

// GbkToUtf8 gbk To utf-8
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// Utf8ToGbk utf-8 to gbk
func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// CompressStr remove blank character
func CompressStr(str string) string {
	if str == "" {
		return ""
	}
	//匹配一个或多个空白符的正则表达式
	reg := regexp.MustCompile("\\s+")
	return reg.ReplaceAllString(str, "")
}
