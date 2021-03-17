package nusoft

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/cliffxzx/gieenm-tools/pkg/utils"
	"github.com/thoas/go-funk"
)

// RecordGroupInfo ...
type RecordGroupInfo struct {
	ID          *NusoftID
	Name        *string
	MemberNames *[]string
}

func (r RecordGroupInfo) String() string {
	return fmt.Sprintf(
		`{ ID: { Serial: %s, Time: %s }, Name: %s}`,
		utils.MustToJSON(r.ID.Serial),
		utils.MustToJSON(r.ID.Time),
		utils.MustToJSON(r.Name),
	)
}

// RecordGroup ...
type RecordGroup struct {
	Includes *[]Record
	Excludes *[]Record
}

func (r RecordGroup) String() string {
	return fmt.Sprintf(
		"Excludes: %s,\nIncludes: %s\n",
		utils.SliceToJSON(*r.Excludes),
		utils.SliceToJSON(*r.Includes),
	)
}

func parseGroupInfoTableRow(fn func(info *RecordGroupInfo)) func(i int, s *goquery.Selection) {
	return func(i int, s *goquery.Selection) {
		IDStr := regexp.MustCompile(`[(,]`).Split(regexp.MustCompile(`[)\s';]`).ReplaceAllString(s.Find("button[type=submit]").AttrOr("onclick", ""), ""), -1)
		idSerial, _ := strconv.ParseInt(IDStr[3], 10, 64)
		idTime, _ := strconv.ParseInt(IDStr[2], 10, 64)

		id := NusoftID{
			Serial: &idSerial,
			Time:   &idTime,
		}
		name := s.Find("td").Eq(0).Text()
		memberNames := strings.Split(s.Find("td").Eq(1).Text(), ", ")

		fn(&RecordGroupInfo{ID: &id, Name: &name, MemberNames: &memberNames})
	}
}

// GetRecordGroupInfos ...
func (f *Firewall) GetRecordGroupInfos() ([]RecordGroupInfo, error) {
	infos := []RecordGroupInfo{}
	err := f.runPages(func(page uint64) (*goquery.Document, error) {
		url := utils.SetQuery(getAddressURL(*f.Host), map[string]string{
			"t":          "1",
			"menu":       "click_v=23\nclick_v=24\nclick_v=26\n",
			"MULTI_LANG": "ch",
			"se":         strconv.FormatUint(page, 10),
		})

		doc, err := f.request("GET", url, "")
		if err != nil {
			return nil, fmt.Errorf("Can't request firewall page: %s", err.Error())
		}

		err = eachTableRow(doc, parseGroupInfoTableRow(func(info *RecordGroupInfo) { infos = append(infos, *info) }))
		if err != nil {
			return nil, fmt.Errorf("Can't parsing table row on firewall: %s", err.Error())
		}

		return doc, nil
	})

	if err != nil {

	}

	return infos, nil
}

// GetRecordGroup ...
func (f *Firewall) GetRecordGroup(info RecordGroupInfo) (RecordGroup, error) {
	doc, err := f.request("POST", getAddressURL(*f.Host), url.Values{
		"q":          {"2"},
		"t":          {"1"},
		"menu":       {"click_v=23\nclick_v=24\nclick_v=26"},
		"MULTI_LANG": {"ch"},
		"id":         {strconv.FormatInt(*info.ID.Time, 10)},
		"n":          {strconv.FormatInt(*info.ID.Serial, 10)},
	}.Encode())

	if err != nil {
		return RecordGroup{}, fmt.Errorf("Can't request firewall page: %s", err.Error())
	}

	var includes, excludes []Record
	table := doc.Find("body > center > form > table.MainTable > tbody > tr.Col > td > table")
	table.Find("select#avail_members > option").Each(func(idx int, s *goquery.Selection) {
		if idx > 0 {
			IDStr, exist := s.Attr("value")
			var id *int64
			if exist {
				id = funk.PtrOf(utils.CheckErr(strconv.ParseInt(IDStr, 10, 64))).(*int64)
			}

			excludes = append(excludes, Record{
				ID:   &NusoftID{Time: id},
				Name: funk.PtrOf(s.Text()).(*string),
			})
		}
	})

	table.Find("select#select_members > option").Each(func(idx int, s *goquery.Selection) {
		if idx > 0 {
			IDStr, exist := s.Attr("value")
			var id *int64
			if exist {
				id = funk.PtrOf(utils.CheckErr(strconv.ParseInt(IDStr, 10, 64))).(*int64)
			}

			includes = append(includes, Record{
				ID:   &NusoftID{Time: id},
				Name: funk.PtrOf(s.Text()).(*string),
			})
		}
	})

	group := RecordGroup{
		Includes: &includes,
		Excludes: &excludes,
	}

	return group, nil
}

// SetRecordGroup ...
func (f *Firewall) SetRecordGroup(info RecordGroupInfo, includes []Record) error {
	body := url.Values{
		"q":          {"2"},
		"t":          {"1"},
		"s":          {"2"},
		"menu":       {"click_v=23\nclick_v=24\nclick_v=26"},
		"MULTI_LANG": {"ch"},
		"id":         {strconv.FormatInt(*info.ID.Time, 10)},
		"n":          {strconv.FormatInt(*info.ID.Serial, 10)},
		"name":       {*info.Name},
	}.Encode()

	body = fmt.Sprintf("%s&%s", body, strings.Join(funk.Map(includes, func(record Record) string {
		return "select_members=" + strconv.FormatInt(*record.ID.Time, 10)
	}).([]string), "&"))

	doc, err := f.request("POST", getAddressURL(*f.Host), body)

	if err != nil {
		return fmt.Errorf("Can't request firewall page: %s", err.Error())
	}

	success := true
	err = eachTableRow(doc, parseGroupInfoTableRow(func(ifo *RecordGroupInfo) {
		if success && *ifo.ID.Serial == *info.ID.Serial && *ifo.ID.Time == *info.ID.Time {
			diff := funk.SubtractString(*ifo.MemberNames, funk.Map(includes, func(record Record) string { return *record.Name }).([]string))
			if len(diff) != 0 {
				success = false
			}
		}
	}))

	if !success {
		return fmt.Errorf("SetRecordGroup fail, the result and input is different")
	}

	if err != nil {
		return fmt.Errorf("Can't parsing table row on firewall: %s", err.Error())
	}

	return nil
}
