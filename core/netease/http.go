package netease

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	_ "github.com/joho/godotenv/autoload"
)

var (
	VideoHost = os.Getenv("VIDEOURL")
	AppKey    = os.Getenv("NIM_APP_KEY")     //网易云appkey
	AppSecret = os.Getenv("NIM_APP_SECREPT") //网易云appsecret
)

var letters = []rune("0123456789abcdef")

//获取随机数
func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

//sha1加密
func b(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

//返回验签
func getCheckSum(timestamp string) (string, string) {
	var nonce string
	for a := 0; a < 128; a++ {
		nonce += randSeq(1)
	}
	checksum := AppSecret + string(nonce) + timestamp
	return string(nonce), b(checksum)
}

func setCommonHeader(r *http.Request, contenType string) {
	timestamp := strconv.FormatInt(time.Now().Unix()-30, 10)
	nonce, checksum := getCheckSum(timestamp)
	r.Header.Set("AppKey", AppKey)
	r.Header.Set("Nonce", nonce)
	r.Header.Set("CurTime", timestamp)
	r.Header.Set("CheckSum", checksum)
	r.Header.Set("Content-Type", contenType)
}

//post请求
func HttpPost(url, data string) (res []byte, err error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return
	}
	setCommonHeader(req, "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"url":  url,
			"data": data,
			"err":  err,
		}).Info("netease post")
		return
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

//postjson请求
func HttpJson(url string, data []byte) (res []byte, err error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return
	}
	setCommonHeader(req, "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"url":  url,
			"data": data,
			"err":  err,
		}).Info("netease post")
		return
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
