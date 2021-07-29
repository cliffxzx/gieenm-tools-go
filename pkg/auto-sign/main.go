// Source: https://github.com/youngle37/Auto_clock/blob/master/clock.py

package autosign

import (
	"errors"
	"fmt"

	"github.com/gocolly/colly/v2"
)

type AutoSign struct {
	colly    *colly.Collector
	username string
	password string
}

func New(username, password string) (*AutoSign, error) {
	autoSign := AutoSign{
		colly: colly.NewCollector(
			colly.IgnoreRobotsTxt(),
			colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 999) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/999 Safari/537.36"),
		),
		username: username,
		password: password,
	}

	return &autoSign, nil
}

// Login ...
func (a *AutoSign) Login() error {
	csrf := ""

	// a.colly.OnRequest(func(r *colly.Request) {
	// 	fmt.Println(*r.Headers)
	// })

	// a.colly.OnResponse(func(r *colly.Response) {
	// 	fmt.Println(string(r.Body))
	// 	fmt.Println(*r.Headers)
	// })

	a.colly.OnHTML(`input[name="_csrf"]`, func(h *colly.HTMLElement) { csrf = h.Attr("value") })
	err := a.colly.Visit("https://portal.ncu.edu.tw/login")
	if err != nil {
		return err
	}

	a.colly.Wait()

	if csrf == "" {
		return errors.New("couldn't find or parsing csrf value")
	}

	form := map[string]string{
		"_csrf":    csrf,
		"language": "CHINESE",
		"password": a.password,
		"username": a.username,
	}

	err = a.colly.Post("https://portal.ncu.edu.tw/login", form)
	if err != nil {
		return err
	}

	a.colly.Wait()

	csrf = ""

	err = a.colly.Visit("https://cis.ncu.edu.tw/HumanSys/login")
	if err != nil {
		return err
	}

	a.colly.Wait()

	if csrf == "" {
		return errors.New("couldn't find or parsing csrf value")
	}

	err = a.colly.Post("https://portal.ncu.edu.tw/leaving", map[string]string{"_csrf": csrf})
	if err != nil {
		return err
	}

	a.colly.Wait()

	a.colly.OnHTMLDetach(`input[name="_csrf"]`)

	return nil
}

// Signin ...
func (a *AutoSign) Signin(partTimeID string) error {
	token := ""

	a.colly.OnHTML(`input[name="_token"]`, func(h *colly.HTMLElement) { token = h.Attr("value") })
	err := a.colly.Visit(fmt.Sprintf("https://cis.ncu.edu.tw/HumanSys/student/stdSignIn/create?ParttimeUsuallyId=%s", partTimeID))
	if err != nil {
		return err
	}

	a.colly.Wait()

	if token == "" {
		return errors.New("couldn't find or parsing token value")
	}

	form := map[string]string{
		"functionName":      "doSign",
		"ParttimeUsuallyId": partTimeID,
		"_token":            token,
	}

	err = a.colly.Post("https://cis.ncu.edu.tw/HumanSys/student/stdSignIn_detail", form)
	if err != nil {
		return err
	}

	a.colly.Wait()

	a.colly.OnHTMLDetach(`input[name="_token"]`)

	return nil
}

// Signout ...
func (a *AutoSign) Signout(partTimeID string, attendWork string) error {
	token, idNo := "", ""

	a.colly.OnHTML(`input[name="_token"]`, func(h *colly.HTMLElement) { token = h.Attr("value") })
	a.colly.OnHTML(`input[id="idNo"]`, func(h *colly.HTMLElement) { idNo = h.Attr("value") })
	err := a.colly.Visit(fmt.Sprintf("https://cis.ncu.edu.tw/HumanSys/student/stdSignIn/create?ParttimeUsuallyId=%s", partTimeID))
	if err != nil {
		return err
	}

	a.colly.Wait()

	if token == "" {
		return errors.New("couldn't find or parsing token value")
	}

	if idNo == "" {
		return errors.New("couldn't find or parsing idNo value")
	}

	form := map[string]string{
		"functionName":      "doSign",
		"ParttimeUsuallyId": partTimeID,
		"AttendWork":        attendWork,
		"idNo":              idNo,
		"_token":            token,
	}

	err = a.colly.Post("https://cis.ncu.edu.tw/HumanSys/student/stdSignIn_detail", form)
	if err != nil {
		return err
	}

	a.colly.Wait()

	a.colly.OnHTMLDetach(`input[name="_token"]`)
	a.colly.OnHTMLDetach(`input[id="idNo"]`)

	return nil
}
