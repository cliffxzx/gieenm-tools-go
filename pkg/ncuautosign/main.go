// Source: https://github.com/youngle37/Auto_clock/blob/master/clock.py

package ncuautosign

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	var headers map[string]string = map[string]string{
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		"Accept-Encoding":           "gzip, deflate, br",
		"Accept-Language":           "zh-TW,zh;q=0.9,en-US;q=0.8,en;q=0.7",
		"Connection":                "keep-alive",
		"Host":                      "portal.ncu.edu.tw",
		"Upgrade-Insecure-Requests": "1",
		"User-Agent":                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36",
	}

	Login(headers, "", "")

	// session = requests.session()
	// login(session, headers, config['account'], config['password'])

	// if args.action == "signin":
	// 		signin(session, config["jobs"][args.jobname]["partTimeId"])
	// 		pass
	// elif args.action == "signout":
	// 		signout(session, config["jobs"][args.jobname]["partTimeId"], config["jobs"][args.jobname]["attendWork"])
	// 		pass
}

// Login ...
func Login(headers map[string]string, account string, password string) (*http.Header, error) {
	req, _ := http.NewRequest("GET", "https://portal.ncu.edu.tw/login", nil)
	for key, val := range headers {
		req.Header.Set(key, val)
	}

	payload := map[string][]string{
		"username": {account},
		"password": {password},
		"language": {"CHINESE"},
	}

	jar, _ := cookiejar.New(nil)
	client := http.Client{Jar: jar}

	resp, _ := client.Do(req)
	doc, _ := goquery.NewDocumentFromResponse(resp)
	defer resp.Body.Close()

	// Get CSRF token from portal login page
	csrf, ok := doc.Find("input[name=\"_csrf\"]").Attr("value")
	payload["_crsf"] = []string{csrf}

	if !ok {
		return nil, errors.New("Couldn't find or parsing csrf value")
	}

	resp, _ = http.PostForm("https://portal.ncu.edu.tw/login", url.Values(payload))
	cookies := resp.Header.Values("Set-Cookie")

	req, _ = http.NewRequest("GET", "https://cis.ncu.edu.tw/HumanSys/login", nil)
	for key, val := range headers {
		req.Header.Set(key, val)
	}

	for _, cookie := range cookies {
		req.Header.Add("Set-Cookie", cookie)
	}

	resp, _ = client.Do(req)
	doc, _ = goquery.NewDocumentFromResponse(resp)
	defer resp.Body.Close()
	csrf, ok = doc.Find("input[name=\"_csrf\"]").Attr("value")
	payload["_crsf"] = []string{csrf}

	if !ok {
		return nil, errors.New("Couldn't find or parsing csrf value")
	}

	resp, _ = http.PostForm("https://portal.ncu.edu.tw/leaving", url.Values{"_csrf": {csrf}})

	return &resp.Header, nil
}

// Signin ...
func Signin(headers *http.Header, partTimeID string) (*http.Header, error) {
	body := map[string][]string{
		"functionName":      {"doSign"},
		"idNo":              {""},
		"ParttimeUsuallyId": {partTimeID},
		"AttendWork":        {""},
	}

	req, _ := http.NewRequest(
		"POST",
		fmt.Sprintf("https://cis.ncu.edu.tw/HumanSys/student/stdSignIn/create?ParttimeUsuallyId=%s", partTimeID),
		strings.NewReader(url.Values(body).Encode()),
	)

	req.Header = *headers
	jar, _ := cookiejar.New(nil)
	client := http.Client{Jar: jar}
	resp, _ := client.Do(req)
	doc, _ := goquery.NewDocumentFromResponse(resp)
	defer resp.Body.Close()
	token, ok := doc.Find("input[name=\"_token\"]").Attr("value")
	if !ok {
		return nil, errors.New("Couldn't find or parsing token value")
	}

	body["_token"] = []string{token}

	req, _ = http.NewRequest(
		"POST",
		"https://cis.ncu.edu.tw/HumanSys/student/stdSignIn_detail",
		strings.NewReader(url.Values(body).Encode()),
	)

	req.Header = *headers
	resp, _ = client.Do(req)

	return &resp.Header, nil
}

// Signout ...
func Signout(headers *http.Header, partTimeID string, attendWork string) (*http.Header, error) {
	req, _ := http.NewRequest(
		"POST",
		fmt.Sprintf("https://cis.ncu.edu.tw/HumanSys/student/stdSignIn/create?ParttimeUsuallyId=%s", partTimeID),
		nil,
	)

	req.Header = *headers
	jar, _ := cookiejar.New(nil)
	client := http.Client{Jar: jar}
	resp, _ := client.Do(req)
	doc, _ := goquery.NewDocumentFromResponse(resp)
	defer resp.Body.Close()

	token, ok := doc.Find("input[name=\"_token\"]").Attr("value")
	if !ok {
		return nil, errors.New("Couldn't find or parsing token value")
	}

	idNo, ok := doc.Find("*[id=\"idNo\"]").Attr("value")
	if !ok {
		return nil, errors.New("Couldn't find or parsing idNo value")
	}

	body := map[string][]string{
		"functionName":      {"doSign"},
		"idNo":              {idNo},
		"ParttimeUsuallyId": {partTimeID},
		"AttendWork":        {attendWork},
		"_token":            {token},
	}

	req, _ = http.NewRequest(
		"POST",
		"https://cis.ncu.edu.tw/HumanSys/student/stdSignIn_detail",
		strings.NewReader(url.Values(body).Encode()),
	)

	req.Header = *headers
	resp, _ = client.Do(req)

	return &resp.Header, nil
}
