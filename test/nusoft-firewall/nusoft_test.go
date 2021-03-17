package test

import (
	"encoding/json"
	"math/rand"
	"net"
	"net/url"
	"testing"

	"github.com/cliffxzx/gieenm-tools/pkg/nusoft-firewall"
	"github.com/cliffxzx/gieenm-tools/pkg/utils"
	"github.com/thoas/go-funk"
)

func fakeFirewall() nusoft.Firewall {
	url, _ := url.Parse("http://140.115.151.237:8080")

	f := *nusoft.New("wifi", *url, "network", "3465134651")

	return f
}

func fakeRecord() nusoft.Record {
	ip, ipAddr, _ := net.ParseCIDR("192.168.4.220/24")
	ipAddr.IP = ip

	macAddr, _ := net.ParseMAC("00:00:5e:00:53:01")
	record := nusoft.Record{
		Name:    funk.PtrOf("gieenm-system-test").(*string),
		IPAddr:  ipAddr,
		MacAddr: &macAddr,
	}

	return record
}

func fakeRecordGroupInfo() nusoft.RecordGroupInfo {
	f := fakeFirewall()

	infos, _ := f.GetRecordGroupInfos()

	return infos[0]
}

func fakeRecordGroup() nusoft.RecordGroup {
	f := fakeFirewall()

	infos, _ := f.GetRecordGroupInfos()
	group, _ := f.GetRecordGroup(infos[0])

	return group
}

func TestGetRecords(t *testing.T) {
	firewall := fakeFirewall()

	records, err := firewall.GetRecords()
	if err != nil || records == nil {
		t.Fatal("Error: ", err)
	}

	t.Logf("Records Count: %d, Value: \n%s", len(records), utils.SliceToJSON(records))
}

func TestAddRecord(t *testing.T) {
	firewall := fakeFirewall()

	record, err := firewall.AddRecord(fakeRecord())
	if err != nil {
		t.Fatal("Error: ", err)
	}

	t.Logf("Add Success: %t", true)
	t.Logf("Records: \n%+v", record.String())

	err = firewall.DelRecord(*record.ID)
	if err != nil {
		t.Fatal("Error: ", err)
	}
}

func TestDelRecord(t *testing.T) {
	firewall := fakeFirewall()

	record, err := firewall.AddRecord(fakeRecord())
	if err != nil {
		t.Fatal("Error: ", err)
	}

	err = firewall.DelRecord(*record.ID)
	if err != nil {
		t.Fatal("Error: ", err)
	}

	t.Logf("Del Success: %t", true)
}

func TestGetRecordGroupInfos(t *testing.T) {
	f := fakeFirewall()

	infos, err := f.GetRecordGroupInfos()
	if err != nil || infos == nil {
		t.Fatal("Error: ", err)
	}

	t.Logf("RecordGroupInfos Count: %d, Value: \n%s", len(infos), utils.SliceToJSON(infos))
}

func TestGetRecordGroup(t *testing.T) {
	f := fakeFirewall()

	info := fakeRecordGroupInfo()

	group, err := f.GetRecordGroup(info)
	if err != nil {
		t.Fatal("Error: ", err)
	}

	t.Logf("RecordGroup: \n%s", group.String())
}

func TestSetRecordGroup(t *testing.T) {
	f := fakeFirewall()
	info := fakeRecordGroupInfo()
	group := fakeRecordGroup()
	originJSON, err := json.Marshal(*group.Includes)

	var deletedRecord []nusoft.Record
	json.Unmarshal(originJSON, &deletedRecord)
	deletedIncludeIdx := rand.Intn(len(deletedRecord))
	deletedRecord = append(deletedRecord[:deletedIncludeIdx], deletedRecord[deletedIncludeIdx+1:]...)

	if err = f.SetRecordGroup(info, deletedRecord); err != nil {
		t.Fatal("Error: ", err)
	}

	deletedGroup, err := f.GetRecordGroup(info)
	if err != nil {
		t.Fatal("Error: ", err)
	}

	t.Logf("Deleted Group: \n%s", deletedGroup.String())

	if err = f.SetRecordGroup(info, *group.Includes); err != nil {
		t.Fatal("Error: ", err)
	}

	group, err = f.GetRecordGroup(info)
	if err != nil {
		t.Fatal("Error: ", err)
	}

	t.Logf("Final Status: \n%s", group.String())
}
