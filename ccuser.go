package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mozillazg/request"
)

const (
	version  = "0.1.1"
	loginURL = "http://8.8.8.8:90/login"
	homeURL  = "http://8.8.8.8:90"
	testURL  = "http://www.baidu.com/"
)

type md6 struct {
	raw []byte
}
type ccuser struct {
	username     string
	raw_password string
	password     string
}

func (*md6) mc(a int) string {
	ret := ""
	b := strings.Split("0123456789ABCDEF", "")
	switch {
	case a == int(' '):
		ret = "+"
	case (a < int('0') && a != int('-') && a != int('.')),
		(a < int('A') && a > int('9')),
		(a > int('Z') && a < int('a') && a != int('_')),
		(a > int('z')):
		ret = "%"
		ret += b[a>>4]
		ret += b[a&15]
	default:
		ret = string(int(a))
	}
	return ret
}

func (*md6) m(a int) int {
	return ((a & 1) << 7) |
		((a & (0x2)) << 5) |
		((a & (0x4)) << 3) |
		((a & (0x8)) << 1) |
		((a & (0x10)) >> 1) |
		((a & (0x20)) >> 3) |
		((a & (0x40)) >> 5) |
		((a & (0x80)) >> 7)
}

func (self *md6) digest() string {
	b := ""
	c := 0xbb
	for n, v := range self.raw {
		c = self.m(int(v))
		c = c ^ (0x35 ^ (n & 0xff))
		b += self.mc(c)
	}
	return b
}

func (u *ccuser) status() bool {
	a := request.NewArgs(new(http.Client))
	resp, err := request.Get(testURL, a)
	if err != nil {
		panic(err)
	}
	url, err := resp.URL()
	if err != nil {
		panic(err)
	}

	if strings.HasPrefix(url.String(), homeURL) {
		fmt.Println("Current Status: can't visit", testURL, "need login")
	} else {
		fmt.Println("Current Status: can visit", testURL)
		return true
	}
	return false
}

func (u *ccuser) login() *request.Args {
	a := request.NewArgs(new(http.Client))
	resp, err := request.Get(loginURL, a)
	if err != nil {
		panic(err)
	}
	url, err := resp.URL()
	if err != nil {
		panic(err)
	}
	if url.String() == loginURL {
		fmt.Println("already logged in")
		return a
	}

	uri := strings.Split(url.String(), "?")[1]
	a.Data = map[string]string{
		"username":             u.username,
		"password":             u.password,
		"uri":                  uri,
		"terminal":             "pc",
		"login_type":           "login",
		"check_passwd":         "0",
		"show_tip":             "block",
		"show_change_password": "block",
		"short_message":        "none",
		"show_captcha":         "none",
		"show_read":            "block",
		"show_assure":          "none",
		"assure_phone":         "",
		"password1":            u.raw_password,
		"new_password":         "",
		"retype_newpassword":   "",
		"captcha_value":        "",
		"save_user":            "1",
		"save_pass":            "1",
		"read":                 "1",
	}

	resp, err = request.Post(loginURL, a)
	if err != nil {
		panic(err)
	}
	html, err := resp.Text()
	if err != nil {
		panic(err)
	}

	if strings.Contains(html, u.username) {
		fmt.Println("login success")
	} else {
		fmt.Println("login fail")
		os.Exit(1)
	}

	return a
}

func (u *ccuser) logout() {
	a := request.NewArgs(new(http.Client))
	a.Data = map[string]string{
		"login_type": "logout",
	}
	_, err := request.Post(loginURL, a)
	if err != nil {
		panic(err)
	}
	time.Sleep(1 * time.Second)
	fmt.Println("logout success")
	u.status()
}

func (u *ccuser) beat(a *request.Args) {
	a.Data = map[string]string{
		"login_type": "heartbeat",
	}
	fmt.Println("send heartbeat.")
	n := 1
	for {
		request.Post(loginURL, a)
		fmt.Printf(
			"\r%s%s", strings.Repeat(".", n), strings.Repeat(" ", 10-n),
		)
		n++
		if n > 10 {
			n = 1
		}
		time.Sleep(50 * time.Second)
	}
}

func main() {
	usage := func() {
		fmt.Fprintf(os.Stderr, "ccuser %s\n", version)
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "ccuser -u USERNAME -p PASSWORD ACTION\n")
		fmt.Fprintf(os.Stderr, "options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "  ACTION: login, logout, status\n")
	}

	v := flag.Bool("V", false, "show version info")
	help := flag.Bool("h", false, "show help info")
	beat := flag.Bool("b", false, "heartbeat mode")
	u1 := flag.String("u", "", "username")
	p1 := flag.String("p", "", "password")
	flag.Usage = usage
	flag.Parse()
	actions := flag.Args()

	if *help {
		usage()
		return
	}
	if *v {
		fmt.Printf("ccuser v%s\n", version)
		return
	}
	username := *u1
	password := *p1
	if username == "" {
		username = os.Getenv("CCUSER_USERNAME")
	}
	if password == "" {
		password = os.Getenv("CCUSER_PASSWORD")
	}

	switch {
	case len(actions) != 1,
		actions[0] != "status" && username == "",
		actions[0] != "status" && password == "":
		usage()
		os.Exit(2)
	}

	m := md6{[]byte(password)}
	cc := ccuser{
		username,
		password,
		m.digest(),
	}

	switch actions[0] {
	case "login":
		s := cc.status()
		a := request.NewArgs(new(http.Client))
		if !s {
			a = cc.login()
		}
		if *beat {
			cc.beat(a)
		}
	case "logout":
		cc.logout()
	case "status":
		cc.status()
	default:
		usage()
		os.Exit(2)
	}
}
