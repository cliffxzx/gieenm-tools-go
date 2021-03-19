package ncuautosign_test

import (
	"testing"

	autosign "github.com/cliffxzx/gieenm-tools/pkg/auto-sign"
)

func TestSignin(t *testing.T) {
	autoSign, err := autosign.New("108502584", "Nn2053163214+")
	if err != nil {
		t.Logf("%s", err)
	}

	err = autoSign.Login()
	if err != nil {
		t.Logf("%s", err)
	}

	err = autoSign.Signin("144469")
	if err != nil {
		t.Logf("%s", err)
	}

	t.Logf("%s", err)
}

func TestSignOut(t *testing.T) {
	// headers, cookies, err := autosign.Login(map[string]string{
	// 	"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
	// 	"Accept-Encoding":           "gzip, deflate, br",
	// 	"Accept-Language":           "zh-TW,zh;q=0.9,en-US;q=0.8,en;q=0.7",
	// 	"Connection":                "keep-alive",
	// 	"Host":                      "portal.ncu.edu.tw",
	// 	"Upgrade-Insecure-Requests": "1",
	// 	"User-Agent":                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36",
	// }, "108502584", "Nn2053163214+")

	// if err != nil {
	// 	t.Logf("%s", err)
	// }

	// headers, err = autosign.Signout(headers, cookies, "144469", "設定mac卡號 防火牆查看異常流量紀錄(1F~4F，WIFI) (所辦&永續中心)筆電：清理download和desktop，掃毒(with更新)，系統更新")
	// t.Logf("%v\n%s", headers, err)
}
