package test

import (
	"os"
	"regexp"
	"testing"

	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/database"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall"
	"github.com/cliffxzx/gieenm-tools/pkg/utils"
	"github.com/joho/godotenv"
)

// LoadEnv loads env vars from .env
func LoadEnv() {
	re := regexp.MustCompile(`^(.*` + "gieenm-tools-go" + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	path := string(rootPath) + `/.env`
	godotenv.Load(path)
}

func TestOnDebug(t *testing.T) {
	LoadEnv()

	database.Init()
	err := firewall.InitFirewalls()
	if err != nil {
		t.Logf("%s", err)
	}

	t.Logf("\n%s", utils.SliceToJSON(*firewall.GetFirewalls()))
}
