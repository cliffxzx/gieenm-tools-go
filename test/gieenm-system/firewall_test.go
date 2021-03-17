package test

import (
	"net"
	"net/url"
	"os"
	"regexp"
	"testing"

	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/database"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/graphql/scalars"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/user"
	"github.com/cliffxzx/gieenm-tools/pkg/utils"
	"github.com/joho/godotenv"
	"github.com/openlyinc/pointy"
)

// LoadEnv loads env vars from .env
func LoadEnv() {
	re := regexp.MustCompile(`^(.*` + "gieenm-tools-go" + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	path := string(rootPath) + `/.env`
	godotenv.Load(path)
}

func init() {
	LoadEnv()
	database.Init()
}

func fakeStudent() *user.User {
	u := user.User{
		Role:      &user.STUDENT,
		Name:      pointy.String("張大三"),
		Email:     pointy.String("clfxzx@gmail.com"),
		StudentID: pointy.String("108502123"),
		Password:  pointy.String("12345678"),
	}

	u.AddOrGet()

	return &u
}

func fakeFirewall() *firewall.Firewall {
	url, _ := url.Parse("192.0.0.1")
	uurl := firewall.URL(*url)
	f := firewall.Firewall{
		Name:     pointy.String("test"),
		Username: pointy.String("test"),
		Password: pointy.String("test"),
		Enable:   pointy.Bool(true),
		Host:     &uurl,
	}

	return &f
}

func fakeGroup() *firewall.RecordGroup {
	g := &[]firewall.RecordGroup{
		{
			Name:       pointy.String("testFake"),
			NusoftID:   &firewall.NusoftID{Serial: pointy.Int64(1), Time: pointy.Int64(1)},
			FirewallID: fakeFirewall().ID,
			Enable:     pointy.Bool(true),
		},
	}

	firewall.AddRecordGroups(g)

	return &(*g)[0]
}

func fakeRecords(us *user.User, group *firewall.RecordGroup) *[]firewall.Record {
	ipAddr, ip, _ := net.ParseCIDR("192.0.0.1/22")
	ip.IP = ipAddr
	fip := scalars.IPAddr(*ip)
	ipAddr2, ip2, _ := net.ParseCIDR("192.0.0.2/22")
	ip2.IP = ipAddr2
	fip2 := scalars.IPAddr(*ip2)
	mac, _ := net.ParseMAC("00-00-5e-00-53-02")
	fmac := scalars.MacAddr(mac)
	mac2, _ := net.ParseMAC("00-00-5e-00-53-02")
	fmac2 := scalars.MacAddr(mac2)
	u := []firewall.Record{{
		Name:     pointy.String("test-1"),
		NusoftID: &firewall.NusoftID{Serial: pointy.Int64(1), Time: pointy.Int64(1)},
		IPAddr:   &fip,
		MacAddr:  &fmac,
		Enable:   pointy.Bool(true),
		User:     us,
		Group:    group,
	}, {
		Name:     pointy.String("test-2"),
		NusoftID: &firewall.NusoftID{Serial: pointy.Int64(2), Time: pointy.Int64(2)},
		IPAddr:   &fip2,
		MacAddr:  &fmac2,
		Enable:   pointy.Bool(true),
		User:     us,
		Group:    group,
	}}

	return &u
}

func TestInitFirewalls(t *testing.T) {
	fw, err := firewall.InitFirewalls()
	if err != nil {
		t.Logf("%s", err)
	}

	t.Logf("\n%s", utils.SliceToJSON(*fw))
}

// TODO: need implements
func TestAddUserRecordLimits(t *testing.T) {
	err := firewall.AddUserRecordLimits(&[]firewall.UserRecordLimit{
		{
			MaxCount: pointy.Int(3),
			User:     fakeStudent(),
			Group:    fakeGroup(),
		},
	})
	if err != nil {
		t.Logf("%s", err)
	}
}

func TestGetRecordsByUserID(t *testing.T) {
	u := fakeStudent()
	_ = u
	records, err := firewall.GetRecordsByUserID(1)
	if err != nil {
		t.Logf("%s", err)
	}

	t.Logf(utils.SliceToJSON(*records))
}

func TestAddRecords(t *testing.T) {
	u := fakeStudent()
	g := fakeGroup()

	r := fakeRecords(u, g)
	err := firewall.AddRecords(r)
	if err != nil {
		t.Logf("%s", err)
	}
}

func TestDelRecords(t *testing.T) {
	u := fakeStudent()
	g := fakeGroup()
	r := fakeRecords(u, g)
	err := firewall.DelRecords(r)
	if err != nil {
		t.Logf("%s", err)
	}

	t.Logf("Success")
}

func TestNusoftToFirewall(t *testing.T) {
	fws, _ := firewall.InitFirewalls()
	err := fws.NusoftToFirewall()
	if err != nil {
		t.Logf("%s", err)
	}
}
