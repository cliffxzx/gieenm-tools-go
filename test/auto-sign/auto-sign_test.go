package ncuautosign_test

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	autosign "github.com/cliffxzx/gieenm-tools/pkg/auto-sign"
)

func TestLogin(t *testing.T) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("NCU Student ID Number: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		t.Logf("%s", err)
	}

	fmt.Printf("NCU Password: ")
	password, err := reader.ReadString('\n')
	if err != nil {
		t.Logf("%s", err)
	}

	autoSign, err := autosign.New(username, password)
	if err != nil {
		t.Logf("%s", err)
	}

	err = autoSign.Login()
	if err != nil {
		t.Logf("%s", err)
	}

	t.Logf("%s", err)
}

func TestSignin(t *testing.T) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("NCU Student ID Number: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		t.Logf("%s", err)
	}

	fmt.Printf("NCU Password: ")
	password, err := reader.ReadString('\n')
	if err != nil {
		t.Logf("%s", err)
	}

	autoSign, err := autosign.New(username, password)
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
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("NCU Student ID Number: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		t.Logf("%s", err)
	}

	fmt.Printf("NCU Password: ")
	password, err := reader.ReadString('\n')
	if err != nil {
		t.Logf("%s", err)
	}

	autoSign, err := autosign.New(username, password)
	if err != nil {
		t.Logf("%s", err)
	}

	err = autoSign.Login()
	if err != nil {
		t.Logf("%s", err)
	}

	err = autoSign.Signout("144469", "設定mac卡號 防火牆查看異常流量紀錄(1F~4F，WIFI) (所辦&永續中心)筆電：清理download和desktop，掃毒(with更新)，系統更新")
	if err != nil {
		t.Logf("%s", err)
	}

	t.Logf("%s", err)
}
