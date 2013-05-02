// webQQ.go
/*
 *quanql:本程序来自网络！
 *GO语言讨论群:102319854
 *GO语言官网:www.golang.org
 *作者: 不死爬虫
 *主页: http://www.gososo.org http://www.daohang361.com/news/index.html
 *见证的轨迹
 *乱码的解决方法
 *	1、dos执行chcp 65001 //修改代码页为utf-8，否则无法通过编译 
 *	2、修改dos窗口字体为Lucida Console，否则显示的字符为乱码 
 *	http://bbs.golang-china.org/viewtopic.php?f=4&t=8&start=10#p93
 *  chcp 命令:
 *  chcp 65001  就是换成UTF-8代码页，在命令行标题栏上点击右键，选择"属性"->"字体"，将字体修改为True Type字体"Lucida Console"，然后点击确定将属性应用到当前窗口
 *  chcp 936 可以换回默认的GBK
 *  chcp 437 是美国英语
 */
package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"hash"
	"http"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strings"
)

var mcookies = make(map[string]*http.Cookie) //保存cookie
var qq string
var pass string
var ptwebqq string               //cookie值登陆和后面用
var skey string                  //cookie值登陆和后面用
var clientid string = "63690451" //随机取
var vfwebqq string               //第二次返回的参数
var psessionid string            //第二次登陆返回的参数
var refer string

func main() {
	fmt.Println("If you run this First time,Please run 'chcp 65001' in then cmd ,then set front 'Lucida Console'")
	fmt.Println("And run 'chcp 936 to back default'")
	if len(os.Args) != 3 {
		fmt.Println("Usage: ", os.Args[0], "365183440 pass")
		in := bufio.NewReader(os.Stdin)
		line, _ := in.ReadString('\n')
		qq = strings.Split(line, " ")[0]
		pass = strings.Split(line, " ")[1]
		pass = strings.TrimSpace(pass)
	} else {
		qq = os.Args[1]
		pass = os.Args[2]
	}
	fmt.Println("qq:" + qq + "pass:" + pass)
	refer = "http://web2-b.qq.com/proxy.html"
	s := getUrl("http://ptlogin2.qq.com:80/check?appid=1003903&uin="+qq, refer)
	var r = string(s[0:len(s)])
	fmt.Println(r)
	var checkType string = ""
	re, _ := regexp.Compile("([^//']*)")
	s2 := re.FindAllString(r, -1)
	checkType = s2[3]
	fmt.Println("checkType:" + checkType)
	var checkcode string = ""
checkImg:
	if strings.Index(checkType, "!") == 0 {
		fmt.Println("Not Need checkcode")
		checkcode = checkType
	} else {
		fmt.Println("Need checkcode")
		checkcode = getImg("http://captcha.qq.com:80/getimage?aid=1003903&uin="+qq+"&vc_type="+checkType, refer)
	}
	fmt.Println("checkcode:" + checkcode)
	//接下来开始登陆了
	//第1次登陆
	loginUrl := "http://ptlogin2.qq.com:80/login?u=" + qq + "&" +
		"p=" + GetCryPass(pass, checkcode) +
		"&verifycode=" + checkcode + "&remember_uin=1&aid=1003903" +
		"&u1=http%3A%2F%2Fweb2.qq.com%2Floginproxy.html%3Fstrong%3Dtrue" +
		"&h=1&ptredirect=0&ptlang=2052&from_ui=1&pttype=1&dumy=&fp=loginerroralert"
	rb := getUrl(loginUrl, refer)
	var loginS string = string(rb[:])
	fmt.Println(loginS)
	if strings.Index(loginS, "成功") > 0 {
		fmt.Println("第一次登陆成功!")
	} else if strings.Index(loginS, "验证码有误") > 0 {
		fmt.Println("验证码有误!重新获取验证码图片")
		goto checkImg
	} else if strings.Index(loginS, "密码有误") > 0 {
		fmt.Println("密码有误!")
		os.Exit(0)
	} else {
		fmt.Println("未知错误!")
		os.Exit(0)
	}
	ptwebqq = mcookies["ptwebqq"].Value
	skey = mcookies["skey"].Value
	fmt.Println("skey:" + skey)
	fmt.Println("ptwebqq:" + ptwebqq)

	//再次登陆，只有这次登陆，才算真正登陆qq，这个时候，如果你qq已经登陆，会把你的qq踢下线，而且此次登陆才算上线
	channelLoginUrl := "http://d.web2.qq.com:80/channel/login2"
	content := "{\"status\":\"\",\"ptwebqq\":\"" + ptwebqq + "\",\"passwd_sig\":\"\",\"clientid\":\"" + clientid + "\"}"
	content = urlEncode(content)                            //urlencode
	content = "r=" + content                                //post的数据
	res := PostUrl(channelLoginUrl, refer, []byte(content)) //post
	re_twice := string(res[:])                              //第二次返回是个json格式,如下,我们要获取 psessionid vfwebqq值
	/*
		{"retcode":0,"result":{"uin":526868457,"cip":1987728778,"index":1075,"port":40831,"status":"online","vfwebqq":"6c47630a8cd98902d38d919420ffb019141ba7f
		ebd6b8ce02b0b69d7b85e0ad2205c8f141a66f364","psessionid":"8368046764001d636f6e6e7365727665725f77656271714031302e3133342e362e31333800006023000005ee026e0
		400e95f671f6d0000000a4058474d6430776c35476d000000286c47630a8cd98902d38d919420ffb019141ba7febd6b8ce02b0b69d7b85e0ad2205c8f141a66f364","user_state":0,"f
		":0}}*/
	fmt.Print(re_twice)
	re, _ = regexp.Compile("(vfwebqq\":\")([^\"]*)(\")")
	s2 = re.FindStringSubmatch(re_twice)
	vfwebqq = s2[2]
	re, _ = regexp.Compile("(psessionid\":\")([^\"]*)(\")")
	s2 = re.FindStringSubmatch(re_twice)
	psessionid = s2[2]
	fmt.Println("psessionid:" + psessionid)
	fmt.Println("vfwebqq:" + vfwebqq)
	if len(vfwebqq) == 0 || len(psessionid) == 0 {
		fmt.Println("登陆失败")
		os.Exit(0)
	}
	//到此登陆成功  调用poll消息函数
	poll()
}

func urlEncode(urlin string) string {
	return http.URLEscape(urlin)
}

func poll() {
	var pollUrl string = "http://d.web2.qq.com:80/channel/poll2?clientid=" + clientid + "&psessionid=" + psessionid
	for {
		fmt.Println("获取消息......................................")
		fmt.Println("Loop Get Message ......................................")
		s := getUrl(pollUrl, refer)
		fmt.Println(string(s[:]))
	}
}

/*
 *获取验证码图片,并且返回输入验证码
 */
func getImg(urlin string, refer string) string {
	imgByte := getUrl(urlin, refer)
	fmt.Print("write img begin and the img length is:")
	fmt.Println(len(imgByte))
	err2 := ioutil.WriteFile("d:/aa.jpg", imgByte, 0)
	if err2 != nil {
		fmt.Println(err2)
		return ""
	}
	fmt.Println("请打开文件:d:/aa.jpg 并输入验证码回车")
	fmt.Println("Need CheckCode Open :d:/aa.jpp input CheckCode")
	var line string
	_, err := fmt.Scanln(&line)
	if err != nil {
		return ""
	}
	return line
}

/*
 *带cookie获取url
 */
func getUrl(urlin string, refer string) []byte {
	url, err := http.ParseURL(urlin)
	checkError(err)
	// build a TCP connection first
	host := url.Host
	conn, err := net.Dial("tcp", host)
	checkError(err)

	// then wrap an HTTP client connection around it
	clientConn := http.NewClientConn(conn, nil)
	if clientConn == nil {
		fmt.Println("Can't build connection")
		os.Exit(1)
	}
	// define the additional HTTP header fields
	header := map[string][]string{
		"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*,q=0.8"},
		"Accept-Language": {"zh-cn,zh;q=0.5"},
		"Accept-Charset":  {"UTF-8,utf-8;q=0.7,*;q=0.7"},
		"Connection":      {"keep-alive"},
		"Referer":         {refer},
		"User-Agent":      {"Mozilla/5.0 (Windows NT 5.1; rv:5.0) Gecko/20100101 Firefox/5.0"},
	}
	// and build the request
	request := http.Request{Method: "GET", URL: url, Header: header}

	for _, value := range mcookies {
		request.AddCookie(value)
	}

	dump, _ := http.DumpRequest(&request, false)
	fmt.Println(string(dump))
	// send the request
	err = clientConn.Write(&request)
	checkError(err)
	// and get the response
	response, err := clientConn.Read(&request)
	checkError(err)
	if response.Status != "200 OK" {
		fmt.Println(response.Status)
		os.Exit(2)
	}
	//Set-cookie
	for i := 0; i < len(response.Cookies()); i++ {
		mcookies[response.Cookies()[i].Name] = response.Cookies()[i]
	}
	const NBUF = 1
	var buf = make([]byte, NBUF)
	//http://code.google.com/p/golang-china/wiki/go_tutorial
	reader := response.Body
	result := make([]byte, 0)
	//我们一个byte取可以有效避免数据覆盖的问题,不过与整块读取略微慢点，暂且这样写吧
	for {
		switch nr, _ := reader.Read(buf[:]); true {
		case nr < 0:
			goto a1
			break
		case nr == 0: // EOF              
			goto a1
			break
		case nr > 0:
			result = append(result, buf[0])
		}
	}
a1:
	return result[0:len(result)]
}

/*
 *带cookie post
 */
func PostUrl(urlin string, refer string, sendBytes []byte) []byte {
	url, err := http.ParseURL(urlin)
	checkError(err)
	// build a TCP connection first
	host := url.Host
	conn, err := net.Dial("tcp", host)
	checkError(err)

	// then wrap an HTTP client connection around it
	clientConn := http.NewClientConn(conn, nil)
	if clientConn == nil {
		fmt.Println("Can't build connection")
		os.Exit(1)
	}
	// define the additional HTTP header fields
	header := map[string][]string{
		"Accept":            {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*,q=0.8"},
		"Accept-Language":   {"zh-cn,zh;q=0.5"},
		"Accept-Charset":    {"UTF-8,utf-8;q=0.7,*;q=0.7"},
		"Content-Type":      {"application/x-www-form-urlencoded"},
		"Connection":        {"keep-alive"},
		"Transfer-Encoding": {"chunked"},
		"Referer":           {refer},
		"User-Agent":        {"Mozilla/5.0 (Windows NT 5.1; rv:5.0) Gecko/20100101 Firefox/5.0"},
	}
	// and build the request
	request := http.Request{Method: "POST", URL: url, Header: header}
	request.ContentLength = (int64)(len(sendBytes))
	for _, value := range mcookies {
		request.AddCookie(value)
	}

	request.Body = &ClosingBuffer{bytes.NewBuffer(sendBytes)}

	dump, _ := http.DumpRequest(&request, false)
	fmt.Println(string(dump))
	// send the request
	err = clientConn.Write(&request)
	checkError(err)
	// and get the response
	response, err := clientConn.Read(&request)
	checkError(err)
	if response.Status != "200 OK" {
		fmt.Println(response.Status)
	}
	//Set-cookie
	for i := 0; i < len(response.Cookies()); i++ {
		mcookies[response.Cookies()[i].Name] = response.Cookies()[i]
	}
	const NBUF = 1
	var buf = make([]byte, NBUF)
	//http://code.google.com/p/golang-china/wiki/go_tutorial
	reader := response.Body
	result := make([]byte, 0)
	//我们一个byte取可以有效避免数据覆盖的问题,不过与整块读取略微慢点，暂且这样写吧
	for {
		switch nr, _ := reader.Read(buf[:]); true {
		case nr < 0:
			goto a2
			break
		case nr == 0: // EOF              
			goto a2
			break
		case nr > 0:
			result = append(result, buf[0])
		}
	}
a2:
	return result[0:len(result)]
}

func checkError(err os.Error) {
	if err != nil {
		fmt.Println("Fatal error ", err.String())
		//os.Exit(1)
	}
}

/*
 *webqq加密方式
 */
func GetCryPass(pass string, code string) string {
	cryPss_3 := Getmd5_3(pass)
	cryPss_3 = strings.ToUpper(hex.EncodeToString([]byte(cryPss_3))) + strings.ToUpper(code)
	r := Getmd5(cryPss_3)
	return strings.ToUpper(hex.EncodeToString(r[0:len(r)]))
}

func Getmd5(original string) []byte {
	var h hash.Hash = md5.New()
	h.Write([]byte(original))
	//fmt.Printf("%x\n", h.Sum()) 
	return h.Sum()
}

func Getmd5_3(in string) string {
	cry := Getmd5(in)
	var r string = string(cry[0:len(cry)])
	cry = Getmd5(r)
	r = string(cry[0:len(cry)])
	cry = Getmd5(r)
	r = string(cry[0:len(cry)])
	return r
}

//实现io.ReadCloser接口
type ClosingBuffer struct {
	*bytes.Buffer
}

func (cb *ClosingBuffer) Close() (err os.Error) {
	//we don't actually have to do anything here, since the buffer is
	//just some data in memory
	//and the error is initialized to no-error
	return
}
