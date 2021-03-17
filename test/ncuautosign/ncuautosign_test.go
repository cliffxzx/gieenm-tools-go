package ncuautosign_test

import (
	"testing"

	"github.com/cliffxzx/gieenm-tools/pkg/ncuautosign"
)

func TestLogin(t *testing.T) {
	ncuautosign.Login(map[string]string{
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		"Accept-Encoding":           "gzip, deflate, br",
		"Accept-Language":           "zh-TW,zh;q=0.9,en-US;q=0.8,en;q=0.7",
		"Connection":                "keep-alive",
		"Host":                      "portal.ncu.edu.tw",
		"Upgrade-Insecure-Requests": "1",
		"User-Agent":                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36",
	}, "108502584", "Nn2053163214+")
}
