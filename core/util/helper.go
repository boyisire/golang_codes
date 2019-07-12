package util

//获取随机数
import (
	"bytes"
	"crypto/subtle"

	"crypto/md5"
	"encoding/hex"
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"math/rand"
	"net/http"
	"path"
	"regexp"
	"strings"
	"time"
)

func IntInArray(val int, array []int) (exist bool) {
	exist = false
	for _, v := range array {
		if val == v {
			exist = true
			return
		}
	}
	return
}

func StrInArray(val string, array []string) (exist bool) {
	exist = false
	for _, v := range array {
		if val == v {
			exist = true
			return
		}
	}
	return
}

//获取文件扩展名
func GetFileExt(disposition string) string {
	reg := regexp.MustCompile(`"(.*)"`)
	disposition = reg.FindString(disposition)
	reg = regexp.MustCompile(`"`)
	disposition = reg.ReplaceAllString(disposition, "")
	return path.Ext(disposition)
}

// Md5加密
func Md5Sum(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// Md5加密
func Md5(str string) string {
	return Md5Sum(str)
}

//获取随机数
func RandStringBytesMaskImprSrc(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)
	var src = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

//获取随机数字
func RandCode(n int) string {
	const letterBytes = "123456789"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)
	var src = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

//获取图片宽高
func GetImageWH(uri string) (width, height int, err error) {
	resp, err := http.Get(uri)
	if err != nil {
		return
	}
	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "image") {
		err = errors.New("not image")
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fd := bytes.NewReader(data)
	im, _, err := image.DecodeConfig(fd)
	if err != nil {
		return
	}
	return im.Width, im.Height, nil
}

//获取手机型号
func GetMobileType(http_user_agent string) map[string]string {
	if http_user_agent == "" {
		return nil
	}
	r := regexp.MustCompile(`\([^\(]+\)`)
	all := r.FindAllString(http_user_agent, -1)
	r = regexp.MustCompile(`(\(|\))`)
	str1 := r.ReplaceAllString(all[0], "")
	//str2 := r.ReplaceAllString(all[1], "")
	tmp := strings.Replace(str1, "Linux; ", "", -1)
	tmp = strings.Replace(tmp, "_", ".", -1)
	tmp = strings.Replace(tmp, " CPU iPhone OS ", "", -1)
	tmp = strings.Replace(tmp, " like Mac OS X", "", -1)
	tmp = strings.Replace(tmp, " Build", "", -1)
	osarr := strings.Split(tmp, ";")
	ios := []string{"iPhone", "iPad", "iPod", "iTouch"}
	//判断是否ios系统
	for _, v := range ios {
		if v == osarr[0] {
			return map[string]string{
				"plateform": osarr[0] + " OS",
				"version":   osarr[0] + " OS",
				"system":    osarr[1] + " OS",
			}
		}
	}
	//其他
	mobile := strings.Split(osarr[1], "/")
	if len(mobile) > 1 {
		return map[string]string{
			"plateform": mobile[0],
			"version":   mobile[1],
			"system":    osarr[0],
		}
	}

	if strings.Replace(osarr[0], " ", "", -1) == "U" {
		tmp := strings.Replace(str1, "Linux; U; ", "", -1)
		tmpsplit1 := strings.Split(tmp, " Build/")
		tmpsplit2 := strings.Split(tmpsplit1[0], "; zh-cn; ")
		return map[string]string{
			"plateform": tmpsplit2[1],
			"version":   "",
			"system":    tmpsplit2[0],
		}
	}
	return nil
}

//比较字符串相同
func Strcmp(in, out string) bool {
	if subtle.ConstantTimeEq(int32(len(in)), int32(len(out))) == 1 {
		if subtle.ConstantTimeCompare([]byte(in), []byte(out)) == 1 {
			return true
		}
		return false
	}
	// Securely compare out to itself to keep constant time, but always return false
	if subtle.ConstantTimeCompare([]byte(out), []byte(out)) == 1 {
		return false
	}
	return false
}

// 生成范围随机值
//
//   len(num)=1; min=0, max=num[0]
//
//   len(num)>1; min=num[0], max=num[1]
func Random(num ...int64) int64 {
	var min, max int64 = 0, 0
	if len(num) == 1 {
		max = num[0]
	} else {
		min = num[0]
		max = num[1]
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(max-min+1) + min
}
