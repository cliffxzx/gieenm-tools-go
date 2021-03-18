package nusoft

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/cliffxzx/gieenm-tools/pkg/utils"
)

// Record ...
type Record struct {
	ID      *NusoftID
	Name    *string
	IPAddr  *net.IPNet
	MacAddr *net.HardwareAddr
}

func (r Record) String() string {
	return fmt.Sprintf(
		`{ ID: { Serial: %s, Time: %s }, Name: %s, IPAddr: %s, MacAddr: %s }`,
		utils.MustToJSON(r.ID.Serial),
		utils.MustToJSON(r.ID.Time),
		utils.MustToJSON(r.Name),
		utils.IPAddrToJSON(r.IPAddr),
		utils.MacAddrToJSON(r.MacAddr),
	)
}

func parseRecordTableRow(fn func(record *Record)) func(i int, s *goquery.Selection) {
	return func(i int, s *goquery.Selection) {
		col := s.Find("td")
		IPMaskStr := strings.Split(strings.Replace(col.Eq(3).Text(), " ", "", -1), "/")
		IDStr := regexp.MustCompile(`[(,]`).Split(regexp.MustCompile(`[)\s';]`).ReplaceAllString(col.Eq(5).Find("button").AttrOr("onclick", "null"), ""), -1)
		idSerial, _ := strconv.ParseInt(IDStr[3], 10, 64)
		idTime, _ := strconv.ParseInt(IDStr[2], 10, 64)

		id := NusoftID{
			Serial: &idSerial,
			Time:   &idTime,
		}
		name := col.Eq(0).Text()
		ip := net.ParseIP(IPMaskStr[0])

		var mask net.IPMask
		if len(IPMaskStr) > 1 {
			mask = utils.ParseSubmask(IPMaskStr[1])
		}

		ipAddr := net.IPNet{IP: ip, Mask: mask}

		var macAddr *net.HardwareAddr
		if mac, err := net.ParseMAC(col.Eq(4).Text()); err == nil {
			macAddr = &mac
		}

		fn(&Record{&id, &name, &ipAddr, macAddr})
	}
}

// GetRecords ...
func (f *Firewall) GetRecords() ([]Record, error) {
	records := []Record{}
	err := f.runPages(func(page uint64) (*goquery.Document, error) {
		url := utils.SetQuery(getAddressURL(*f.Host), map[string]string{
			"menu":       "click_v=23\nclick_v=24\nclick_v=25\n",
			"MULTI_LANG": "ch",
			"se":         strconv.FormatUint(page, 10),
		})

		doc, err := f.request("GET", url, "")
		if err != nil {
			return nil, fmt.Errorf("Can't request firewall page: %s", err.Error())
		}

		err = eachTableRow(doc, parseRecordTableRow(func(record *Record) { records = append(records, *record) }))

		if err != nil {
			return nil, fmt.Errorf("Can't parsing table row on firewall: %s", err.Error())
		}

		return doc, nil
	})

	if err != nil {
		return nil, err
	}

	return records, nil
}

// AddRecord ...
func (f *Firewall) AddRecord(record Record) (Record, error) {
	if f.pageRowCount == nil || f.recordCount == nil {
		return Record{}, errors.New("operation add record must from use New func generate firewall")
	}

	var result *Record

	doc, err := f.request("POST", getAddressURL(*f.Host), url.Values{
		"q":          {"1"},
		"s":          {"1"},
		"MULTI_LANG": {"ch"},
		"menu":       {"click_v=23\nclick_v=24\nclick_v=25\n"},
		"ipv":        {"0"},
		"adstartip4": {"0.0.0.0"},
		"interface":  {"All"},
		"n":          {strconv.FormatUint(*f.recordCount+1, 10)},
		"id":         {time.Now().UTC().Format("20060102150405")},
		"name":       {*record.Name},
		"ip":         {record.IPAddr.IP.String()},
		"ip4":        {record.IPAddr.IP.String()},
		"nm4":        {net.IP(record.IPAddr.Mask).String()},
		"mask":       {net.IP(record.IPAddr.Mask).String()},
		"mac":        {record.MacAddr.String()},
		"admac":      {record.MacAddr.String()},
	}.Encode())

	if err != nil {
		return Record{}, fmt.Errorf("Can't request firewall page: %s", err.Error())
	}

	// url := utils.SetQuery(getAddressURL(*f.Host), map[string]string{
	// 	"menu":       "click_v=23\nclick_v=24\nclick_v=25\n",
	// 	"MULTI_LANG": "ch",
	// 	"se":         strconv.FormatUint(*f.recordCount+1, 10),
	// })

	// doc, err = f.request("GET", url, "")
	// if err != nil {
	// 	return Record{}, fmt.Errorf("Can't request firewall page: %s", err.Error())
	// }

	err = eachTableRow(doc, parseRecordTableRow(func(rec *Record) {
		if *rec.Name == *record.Name {
			result = rec
		}
	}))

	if err != nil {
		return Record{}, fmt.Errorf("Can't parsing table row on firewall: %s", err.Error())
	}

	if result == nil {
		return Record{}, fmt.Errorf("Unknown Error on AddRecord in Firewall")
	}

	*f.recordCount++

	return *result, nil
}

// DelRecord ...
func (f *Firewall) DelRecord(id NusoftID) error {
	if f.pageRowCount == nil || f.recordCount == nil {
		return errors.New("operation add record must from use New func generate firewall")
	}

	doc, err := f.request("POST", getAddressURL(*f.Host), url.Values{
		"q":          {"3"},
		"s":          {"3"},
		"MULTI_LANG": {"ch"},
		"menu":       {"click_v=23\nclick_v=24\nclick_v=25\n"},
		"id":         {strconv.FormatInt(*id.Time, 10)},
		"n":          {strconv.FormatInt(*id.Serial, 10)},
	}.Encode())

	if err != nil {
		return fmt.Errorf("Can't request firewall page: %s", err.Error())
	}

	success := true
	url := utils.SetQuery(getAddressURL(*f.Host), map[string]string{
		"menu":       "click_v=23\nclick_v=24\nclick_v=25\n",
		"MULTI_LANG": "ch",
		"se":         strconv.FormatInt(*id.Serial, 10),
	})

	doc, err = f.request("GET", url, "")
	if err != nil {
		return fmt.Errorf("Can't request firewall page: %s", err.Error())
	}

	err = eachTableRow(doc, parseRecordTableRow(func(record *Record) {
		if success && *id.Serial == *record.ID.Serial && *id.Time == *record.ID.Time {
			success = false
		}
	}))

	if err != nil {
		return fmt.Errorf("Can't parsing table row on firewall: %s", err.Error())
	}

	if !success {
		return fmt.Errorf("Unknown Error on AddRecord in Firewall")
	}

	*f.recordCount--

	return nil
}
