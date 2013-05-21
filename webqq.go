/*
* 尝试新特性，或者尝试新语法
 */
package main

import (
	// "crypto/md5"
	// "encoding/hex"
	// "fmt"
	// "hash"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	// "time"
	// "net/http/cookiejar"
	"os"
	"strings"
)

var (
	QQ = "2690371552"

// pass := "123456"
// vc := "EKWJ"
)

var BaseHeader = map[string]string{
	"Accept":          "*/*",
	"Accept-Encoding": "gzip,deflate",
	"Accept-Language": "zh-cn,zh;q=0.8,en-us;q=0.5,en;q=0.3",
	"Connection":      "keep-alive",
	"User-Agent":      "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:20.0) Gecko/20100101 Firefox/20.0",
}

func main() {

	log.SetFlags(log.Lshortfile)
	//appid=1003903

	resp, err := CheckVerifyCode()
	if err != nil {
		log.Println(err)
	}
	code := GetVerifyCode(*resp)
	log.Println(code)

}

func CheckVerifyCode() (*http.Response, error) {
	log.Println("start func checkVC()")
	v := url.Values{}
	v.Add("uin", QQ)
	v.Add("r", 0.12431234)
	log.Println(v.Encode())
	// url := "https://ssl.ptlogin2.qq.com/check?uin=357088531&appid=1003903&js_ver=10029&js_type=0&login_sig=f1ll0J3TOEMtw6nwmDO832a--xFelS-IMsbL8CGIBOHMWPhdak8IIlSD6USs0aMm&u1=http%3A%2F%2Fweb.qq.com%2Floginproxy.html&r=0.9890235673936297"
	url := "https://ssl.ptlogin2.qq.com/check?uin=" + QQ // + "&r=0.9890235673936297"

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}

	for key, value := range BaseHeader {
		req.Header.Set(key, value)
	}

	// req.Header.Set("Cookie", "pgv_pvid=4335859410; pgv_info=pgvReferrer=&ssid=s3944121500; uikey=40ae17426019d961b53020d402710924fb62a93019e805af496463433c1069c2; chkuin=357088531")
	// req.Header.Set("Host", "ui.ptlogin2.qq.com")
	// req.Header.Set("Referer", "https://ui.ptlogin2.qq.com/cgi-bin/login?target=self&style=5&mibao_css=m_webqq&appid=1003903&enable_qlogin=0&no_verifyimg=1&s_url=http%3A%2F%2Fweb.qq.com%2Floginproxy.html&f_url=loginerroralert&strong_login=1&login_state=10&t=20130417001")

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	if resp.StatusCode == 200 {
		// body, _ := ioutil.ReadAll(resp.Body)
		// bodystr := string(body)
		// log.Println(bodystr)
	}
	// log.Printf("\n%T %v\n", resp.Cookies(), resp.Cookies())
	// for _, b := range resp.Cookies() {
	//     log.Printf("%T  %v", b, b)
	// }
	return resp, nil
}

func GetVerifyCode(resp http.Response) string {
	log.Println("start func GetVerifyCode")
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	bodystr := string(body)
	log.Println(bodystr)

	//返回结果0，表示不需要输入验证码
	ok := strings.Contains(strings.ToLower(bodystr), "ptui_checkvc('0'")
	if ok {
		substr := strings.Split(bodystr, "','")
		return substr[1]
	} else {
		return getVCImg()
	}
}

func getVCImg() string {
	log.Println("start func getImg()")
	file := "qq-verifycode.jpeg"

	url := "https://ssl.captcha.qq.com/getimage?&uin=357088531&aid=1003903&r=0.08379581950067883"
	// url := "https://ssl.ptlogin2.qq.com/check?uin=357088531"

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip,deflate")
	req.Header.Set("Accept-Language", "zh-cn,zh;q=0.8,en-us;q=0.5,en;q=0.3")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:20.0) Gecko/20100101 Firefox/20.0")

	req.Header.Set("Cookie", "pgv_pvid=4335859410; pgv_info=pgvReferrer=&ssid=s3944121500; ptisp=ctc; verifysession=h00d963252d82920a46ba445ab48cc6e1044b68058ef89cc93d8c1539245e735b9e0ee01ee9304cbaf1a38a18c58bdba11c; ptui_loginuin=357088531")
	req.Header.Set("Host", "ssl.captcha.qq.com")
	req.Header.Set("Referer", "https://ui.ptlogin2.qq.com/cgi-bin/login?target=self&style=5&mibao_css=m_webqq&appid=1003903&enable_qlogin=0&no_verifyimg=1&s_url=http%3A%2F%2Fweb.qq.com%2Floginproxy.html&f_url=loginerroralert&strong_login=1&login_state=10&t=20130417001")

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	log.Println(resp.StatusCode) //200
	log.Println(resp.Status)     //200 OK
	log.Println(resp.Body)       //pointer

	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		fout, err := os.Create(file)
		if err != nil {
			log.Println("err:", err)
		}

		fout.Write(body)
		// log.Println(body)
		defer func() {
			resp.Body.Close()
			fout.Close()
		}()

	}
	return ""
}
