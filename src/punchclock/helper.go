package punchclock

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"time"
)

var header = map[string]string{
	"Host":            "my.ntu.edu.tw",
	"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.63 Safari/537.36",
	"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
	"Accept-Language": "zh-TW,zh;q=0.8,en-US;q=0.5,en;q=0.3",
	"Connection":      "keep-alive",
}

func set_request_header(req *http.Request) {
	for k, v := range header {
		req.Header.Set(k, v)
	}
}

func Dump_response_body(res *http.Response) {
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error : %s", err)
	}
	fmt.Printf("Response body:\n %s\n", resBody)
}

func Dump_status(res *http.Response) {
	fmt.Println(res.Status)
}

func Dump_header(res *http.Response) {
	for k, v := range res.Header {
		fmt.Println(k, "value is", v)
	}
}

func Dump_cookie(client *http.Client) {
	u_myntu, err := url.Parse("https://my.ntu.edu.tw")

	if err != nil {
		log.Fatalf("Paring NTU url runs into error %s\n", err)
	}

	cookies_myntu := client.Jar.Cookies(u_myntu)
	fmt.Println("Cookies stored in the site ", u_myntu.Hostname())
	for k, v := range cookies_myntu {
		fmt.Println(k, "value is", v)
	}

	u_portail, err := url.Parse("https://web2.cc.ntu.edu.tw")
	if err != nil {
		log.Fatalf("Paring NTU url runs into error %s\n", err)
	}
	cookies_portail := client.Jar.Cookies(u_portail)
	fmt.Println("Cookies stored in the site ", u_portail.Hostname())
	for k, v := range cookies_portail {
		fmt.Println(k, "value is", v)
	}
}

// The client log the SSLKey and equiped with CookieJar
func getClient(debug bool) (*http.Client, error) {
	// Create client with cookiejar
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}
	if debug {
		sslkeyfile := "SSLKEYLOGFILE"
		keyLogWriter, err := os.OpenFile(sslkeyfile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
		if err != nil {
			log.Fatalf("Opening SSLKEYLOGFILE run into error %s", err)
		}
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				KeyLogWriter:       keyLogWriter,
			},
		}
	}
	return client, nil
}

func wait_a_while(w_min int32, w_sec int32) {
	a := rand.Int31n(w_min)
	b := rand.Int31n(w_sec)
	time.Sleep(time.Duration(a)*time.Minute + time.Duration(b)*time.Second)
}
