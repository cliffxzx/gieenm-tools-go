package test

import (
	"net"
	"net/url"
	"os"
	"regexp"
	"testing"

	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/base/scalars"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/database"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/common"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/group"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/limit"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/record"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/sync"
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
		Host:     &uurl,
	}

	return &f
}

func fakeGroup() *group.Group {
	g := &[]group.Group{
		{
			Name:       pointy.String("testFake"),
			NusoftID:   &common.NusoftID{Serial: pointy.Int64(1), Time: pointy.Int64(1)},
			FirewallID: fakeFirewall().ID,
		},
	}

	group.Adds(g)

	return &(*g)[0]
}

func fakeRecords(us *user.User, group *group.Group) *[]record.Record {
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
	u := []record.Record{{
		Name:     pointy.String("test-1"),
		NusoftID: &common.NusoftID{Serial: pointy.Int64(1), Time: pointy.Int64(1)},
		IPAddr:   &fip,
		MacAddr:  &fmac,
		User:     us,
		Group:    group,
	}, {
		Name:     pointy.String("test-2"),
		NusoftID: &common.NusoftID{Serial: pointy.Int64(2), Time: pointy.Int64(2)},
		IPAddr:   &fip2,
		MacAddr:  &fmac2,
		User:     us,
		Group:    group,
	}}

	return &u
}

func TestInitFirewalls(t *testing.T) {
	err := firewall.InitFirewalls()
	if err != nil {
		t.Logf("%s", err)
	}

	t.Logf("\n%s", utils.SliceToJSON(firewall.GetFirewalls()))
}

// TODO: need implements
func TestAddUserRecordLimits(t *testing.T) {
	err := limit.Adds(&[]limit.Limit{
		{
			RecordMaxCount: pointy.Int(3),
			User:           fakeStudent(),
			Group:          fakeGroup(),
		},
	})
	if err != nil {
		t.Logf("%s", err)
	}
}

func TestGetRecordsByUserID(t *testing.T) {
	u := fakeStudent()
	_ = u
	records, err := record.GetsByUserID(1)
	if err != nil {
		t.Logf("%s", err)
	}

	t.Logf(utils.SliceToJSON(*records))
}

func TestAddRecords(t *testing.T) {
	u := fakeStudent()
	g := fakeGroup()

	r := fakeRecords(u, g)
	err := record.Adds(r)
	if err != nil {
		t.Logf("%s", err)
	}
}

func TestDelRecords(t *testing.T) {
	u := fakeStudent()
	g := fakeGroup()
	r := fakeRecords(u, g)
	err := record.Dels(r)
	if err != nil {
		t.Logf("%s", err)
	}

	t.Logf("Success")
}

func TestNusoftToFirewall(t *testing.T) {
	err := firewall.InitFirewalls()
	if err != nil {
		t.Logf("%s", err)
	}

	subnetRaw := map[string]string{
		// 109
		"20180205214124,3": "192.168.1.0/24",
		// 114
		"20180203234922,1": "192.168.2.0/24",
		// 115-1,115-2
		"20180204001228,2": "192.168.3.0/24",
		// 1f 管制群組
		"20190925091129,4": "192.168.0.0/16",
		// 206
		"20170629102323,1": "192.168.1.0/24",
		// 208
		"20170629102325,2": "192.168.2.0/24",
		// 311A
		"20170629102327,3": "192.168.3.0/24",
		// 2f 管制群組
		"20200427155607,4": "192.168.0.0/16",
		// 305-306
		"20170629105419,1": "192.168.1.0/24",
		// 307-310
		"20170629114059,2": "192.168.2.0/24",
		// 405
		"20170629200244,3": "192.168.3.0/24",
		// 405-1
		"20180202235136,2": "192.168.1.0/24",
		// 407-408
		"20170629110845,3": "192.168.2.0/24",
		// 410
		"20180202225649,1": "192.168.3.0/24",
		// 4f 管制群組
		"20190615104456,1": "192.168.0.0/16",
		// gieenm-tools-tokenStr
		"20210319014124,2": "192.168.2.0/16",
	}

	defaultSubnets := map[string]*scalars.IPAddr{}
	for key, raw := range subnetRaw {
		_, subnet, _ := net.ParseCIDR(raw)
		tmp := scalars.IPAddr(*subnet)
		defaultSubnets[key] = &tmp
	}

	err = sync.SyncNusoftToDatabase(defaultSubnets)
	if err != nil {
		t.Logf("%s", err)
	}
}
