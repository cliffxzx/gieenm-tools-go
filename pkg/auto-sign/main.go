// Source: https://github.com/youngle37/Auto_clock/blob/master/clock.py

package autosign

import (
	"errors"
	"fmt"

	"github.com/gocolly/colly/v2"
)

type AutoSign struct {
	colly    *colly.Collector
	account  string
	password string
}

func (AutoSign) newColly() *colly.Collector {
	return colly.NewCollector(
		colly.IgnoreRobotsTxt(),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36"),
	)
}

func New(account, password string) (*AutoSign, error) {
	autoSign := AutoSign{
		account:  account,
		password: password,
	}

	return &autoSign, nil
}

// Login ...
func (a *AutoSign) Login() error {
	csrf := ""

	a.colly = a.newColly()

	a.colly.OnHTML(`input[name="_csrf"]`, func(h *colly.HTMLElement) { csrf = h.Attr("value") })
	a.colly.Visit("https://portal.ncu.edu.tw/login")
	a.colly.Wait()

	if csrf == "" {
		return errors.New("couldn't find or parsing csrf value")
	}

	a.colly.OnRequest(func(r *colly.Request) {
		fmt.Println(*r.Headers)
	})

	a.colly.OnResponse(func(r *colly.Response) {
		fmt.Println(string(r.Body))
		fmt.Println(*r.Headers)
	})

	a.colly.Post("https://portal.ncu.edu.tw/login", map[string]string{
		"_csrf":    csrf,
		"language": "CHINESE",
		"account":  a.account,
		"password": a.password,
	})
	a.colly.Wait()

	a.colly.Visit("https://cis.ncu.edu.tw/HumanSys/login")
	a.colly.Wait()

	if csrf == "" {
		return errors.New("couldn't find or parsing csrf value")
	}

	a.colly.Post("https://portal.ncu.edu.tw/leaving", map[string]string{"_csrf": csrf})

	a.colly.Wait()

	return nil
}

// Signin ...
// func (a *AutoSign) Signin(partTimeID string) error {
// 	body := map[string][]string{
// 		"functionName":      {"doSign"},
// 		"idNo":              {""},
// 		"ParttimeUsuallyId": {partTimeID},
// 		"AttendWork":        {""},
// 	}

// 	req, err := http.NewRequest(
// 		"GET",
// 		fmt.Sprintf("https://cis.ncu.edu.tw/HumanSys/student/stdSignIn/create?ParttimeUsuallyId=%s", partTimeID),
// 		strings.NewReader(url.Values(body).Encode()),
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	req.Header = *a.headers

// 	resp, err := a.client.Do(req)
// 	if err != nil {
// 		return err
// 	}

// 	doc, err := goquery.NewDocumentFromReader(resp.Body)
// 	if err != nil {
// 		return err
// 	}

// 	defer resp.Body.Close()

// 	token, ok := doc.Find("input[name=\"_token\"]").Attr("value")
// 	if !ok {
// 		return errors.New("couldn't find or parsing token value")
// 	}

// 	body["_token"] = []string{token}

// 	req, err = http.NewRequest(
// 		"POST",
// 		"https://cis.ncu.edu.tw/HumanSys/student/stdSignIn_detail",
// 		strings.NewReader(url.Values(body).Encode()),
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	req.Header = *a.headers
// 	_, err = a.client.Do(req)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// // Signout ...
// func Signout(headers *http.Header, cookies *cookiejar.Jar, partTimeID string, attendWork string) (*http.Header, error) {
// 	req, _ := http.NewRequest(
// 		"POST",
// 		fmt.Sprintf("https://cis.ncu.edu.tw/HumanSys/student/stdSignIn/create?ParttimeUsuallyId=%s", partTimeID),
// 		nil,
// 	)

// 	req.Header = *headers
// 	client := http.Client{Jar: cookies}
// 	resp, _ := client.Do(req)
// 	doc, _ := goquery.NewDocumentFromReader(resp.Body)
// 	defer resp.Body.Close()

// 	token, ok := doc.Find("input[name=\"_token\"]").Attr("value")
// 	if !ok {
// 		return nil, errors.New("couldn't find or parsing token value")
// 	}

// 	idNo, ok := doc.Find("*[id=\"idNo\"]").Attr("value")
// 	if !ok {
// 		return nil, errors.New("couldn't find or parsing idNo value")
// 	}

// 	body := map[string][]string{
// 		"functionName":      {"doSign"},
// 		"idNo":              {idNo},
// 		"ParttimeUsuallyId": {partTimeID},
// 		"AttendWork":        {attendWork},
// 		"_token":            {token},
// 	}

// 	req, _ = http.NewRequest(
// 		"POST",
// 		"https://cis.ncu.edu.tw/HumanSys/student/stdSignIn_detail",
// 		strings.NewReader(url.Values(body).Encode()),
// 	)

// 	req.Header = *headers
// 	resp, _ = client.Do(req)

// 	return &resp.Header, nil
// }
