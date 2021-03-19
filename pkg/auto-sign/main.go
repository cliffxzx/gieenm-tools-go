// Source: https://github.com/youngle37/Auto_clock/blob/master/clock.py

package autosign

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type AutoSign struct {
	client   *http.Client
	headers  *http.Header
	account  string
	password string
}

func New(account, password string) (*AutoSign, error) {
	headers := &http.Header{
		"Accept-Language":           {"zh-TW,zh;q=0.9,en-US;q=0.8,en;q=0.7"},
		"Host":                      {"portal.ncu.edu.tw"},
		"Upgrade-Insecure-Requests": {"1"},
		"Accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		"Accept-Encoding":           {"gzip, deflate, br"},
		"Content-Type":              {"application/x-www-form-urlencoded"},
		"Cache-Control":             {"max-age=0"},
		"Connection":                {"keep-alive"},
		"Origin":                    {"https://portal.ncu.edu.tw"},
		"Referer":                   {"https://portal.ncu.edu.tw/login"},
		"sec-ch-ua":                 {"\"Google Chrome\";v=\"89\", \"Chromium\";v=\"89\", \";Not A Brand\";v=\"99\""},
		"sec-ch-ua-mobile":          {"?0"},
		"Sec-Fetch-Dest":            {"document"},
		"Sec-Fetch-Mode":            {"navigate"},
		"Sec-Fetch-Site":            {"same-origin"},
		"Sec-Fetch-User":            {"?1"},
		"User-Agent":                {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36"},
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	autoSign := AutoSign{
		client: &http.Client{
			Jar: jar,
		},
		headers:  headers,
		account:  account,
		password: password,
	}

	return &autoSign, nil
}

func (a *AutoSign) SetHeaders(headers *http.Header) {
	for key, header := range *headers {
		a.headers.Set(key, header[0])
	}
}

// Login ...
func (a *AutoSign) Login() error {
	req, err := http.NewRequest("GET", "https://portal.ncu.edu.tw/login", nil)
	if err != nil {
		return err
	}

	req.Header = *a.headers

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// Get CSRF token from portal login page
	csrf, ok := doc.Find("input[name=\"_csrf\"]").Attr("value")
	if !ok {
		return errors.New("couldn't find or parsing csrf value")
	}

	payload := map[string][]string{
		"account":  {a.account},
		"password": {a.password},
		"language": {"CHINESE"},
		"_crsf":    {csrf},
	}

	req, err = http.NewRequest("POST", "https://portal.ncu.edu.tw/login", strings.NewReader(url.Values(payload).Encode()))
	if err != nil {
		return err
	}

	req.Header = *a.headers

	resp, err = a.client.Do(req)
	if err != nil {
		return err
	}

	doc, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf(doc.Html())

	req, err = http.NewRequest("GET", "https://cis.ncu.edu.tw/HumanSys/login", nil)
	if err != nil {
		return err
	}

	resp, err = a.client.Do(req)
	if err != nil {
		return err
	}

	doc, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	csrf, ok = doc.Find("input[name=\"_csrf\"]").Attr("value")
	if !ok {
		return errors.New("couldn't find or parsing csrf value")
	}

	payload["_crsf"] = []string{csrf}

	resp, err = http.PostForm("https://portal.ncu.edu.tw/leaving", url.Values{"_csrf": {csrf}})
	if err != nil {
		return err
	}

	return nil
}

// Signin ...
func (a *AutoSign) Signin(partTimeID string) error {
	body := map[string][]string{
		"functionName":      {"doSign"},
		"idNo":              {""},
		"ParttimeUsuallyId": {partTimeID},
		"AttendWork":        {""},
	}

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://cis.ncu.edu.tw/HumanSys/student/stdSignIn/create?ParttimeUsuallyId=%s", partTimeID),
		strings.NewReader(url.Values(body).Encode()),
	)

	if err != nil {
		return err
	}

	req.Header = *a.headers

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	token, ok := doc.Find("input[name=\"_token\"]").Attr("value")
	if !ok {
		return errors.New("couldn't find or parsing token value")
	}

	body["_token"] = []string{token}

	req, err = http.NewRequest(
		"POST",
		"https://cis.ncu.edu.tw/HumanSys/student/stdSignIn_detail",
		strings.NewReader(url.Values(body).Encode()),
	)
	if err != nil {
		return err
	}

	req.Header = *a.headers
	_, err = a.client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

// Signout ...
func Signout(headers *http.Header, cookies *cookiejar.Jar, partTimeID string, attendWork string) (*http.Header, error) {
	req, _ := http.NewRequest(
		"POST",
		fmt.Sprintf("https://cis.ncu.edu.tw/HumanSys/student/stdSignIn/create?ParttimeUsuallyId=%s", partTimeID),
		nil,
	)

	req.Header = *headers
	client := http.Client{Jar: cookies}
	resp, _ := client.Do(req)
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	defer resp.Body.Close()

	token, ok := doc.Find("input[name=\"_token\"]").Attr("value")
	if !ok {
		return nil, errors.New("couldn't find or parsing token value")
	}

	idNo, ok := doc.Find("*[id=\"idNo\"]").Attr("value")
	if !ok {
		return nil, errors.New("couldn't find or parsing idNo value")
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
