package nusoft

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/cliffxzx/gieenm-tools/pkg/utils"
)

// Firewall ...
type Firewall struct {
	Host         *url.URL
	Username     *string
	Password     *string
	recordCount  *uint64
	pageRowCount *uint64
}

// String ...
func (f *Firewall) String() string {

	return fmt.Sprintf(
		`{ Host: %s, Username: %s, Password: %s, RecordCount: %s, PageRowCount: %s }`,
		utils.URLToJSON(f.Host),
		utils.MustToJSON(f.Username),
		utils.MustToJSON(f.Password),
		utils.MustToJSON(f.recordCount),
		utils.MustToJSON(f.pageRowCount),
	)
}

func (f *Firewall) runPages(fn func(page uint64) (*goquery.Document, error)) error {
	pageCursor := uint64(0)
	pageInfo := PageInfo{}

	for {
		doc, err := fn(pageCursor + 1)
		if err != nil {
			return err
		}

		getPageInfo(doc, &pageInfo)
		if pageInfo.now >= pageInfo.total {
			break
		}

		pageCursor += *f.pageRowCount
	}

	return nil
}

func (f *Firewall) request(method string, url *url.URL, body string) (*goquery.Document, error) {
	req, err := http.NewRequest(method, url.String(), strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	switch method {
	case "GET":
	case "POST":
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	default:
		return nil, errors.New("Undefined method, the method must be upper case")
	}

	req.SetBasicAuth(*f.Username, *f.Password)
	res, err := http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != 200 {
		return nil, err
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	return doc, err
}

// GetPageRowCountByHTML ...
func (f *Firewall) GetPageRowCountByHTML() *uint64 {
	url := utils.SetQuery(getSettingURL(*f.Host), map[string]string{
		"menu":       "click_v=1\nclick_v=6\nclick_v=7\n",
		"MULTI_LANG": "ch",
	})

	doc, err := f.request("GET", url, "")
	if err != nil {
		panic(fmt.Errorf("Can't request firewall page: %s", err.Error()))
	}

	rowCountStr, ok := doc.Find("input[name=lineperpage]").Attr("value")
	if !ok {
		panic(fmt.Errorf("Can't get page count: %s", err.Error()))
	}

	rowCount, err := strconv.ParseUint(rowCountStr, 10, 64)
	if err != nil {
		panic(fmt.Errorf("Can't parse page count: %s", err.Error()))
	}

	return &rowCount
}

// Firewalls ...
type Firewalls []Firewall

// String ...
func (f *Firewalls) String() string {
	result := "[\n"
	for _, firewall := range *f {
		result += fmt.Sprintf("\t%s", firewall.String())
		result += ",\n"
	}
	result += "]"
	return result
}

// PageInfo ...
type PageInfo struct{ now, total int }

// NusoftID ...
type NusoftID struct {
	Serial *int64
	Time   *int64
}

// New ...
func New(name string, host url.URL, username, password string) *Firewall {
	fw := &Firewall{
		Host:     &host,
		Username: &username,
		Password: &password,
	}

	fw.pageRowCount = fw.GetPageRowCountByHTML()

	records, _ := fw.GetRecords()
	recordsLength := uint64(len(records))
	fw.recordCount = &recordsLength

	return fw
}

func getAddressURL(url url.URL) *url.URL {
	page, err := url.Parse("cgi-bin/address.cgi")
	if err != nil {
		panic(fmt.Sprintf("Not valid page url: %s", err.Error()))
	}

	return url.ResolveReference(page)
}

func getSettingURL(url url.URL) *url.URL {
	page, err := url.Parse("cgi-bin/setting.cgi")
	if err != nil {
		panic(fmt.Sprintf("Not valid page url: %s", err.Error()))
	}

	return url.ResolveReference(page)
}

// TODO: go routine
func getPageInfo(doc *goquery.Document, pageInfo *PageInfo) {
	pageNowDoc := doc.Find(`input[name="cp1"]`)
	pageTotalDoc := doc.Find("tr.list_tool_text_attr > td")

	// Get Html Input Element Value
	if pageNowDoc.Length() > 0 {
		now, _ := strconv.ParseInt(pageNowDoc.AttrOr("value", "1"), 10, 32)
		pageInfo.now = int(now)
	}

	if pageTotalDoc.Length() > 0 {
		total, _ := strconv.ParseInt(strings.TrimSpace(strings.Split(pageTotalDoc.First().Text(), "/")[1]), 10, 32)
		pageInfo.total = int(total)
	}
}

func eachTableRow(doc *goquery.Document, fn func(i int, s *goquery.Selection)) error {
	cols := doc.Find("body > center > form > table.FixedTable tr.Col")
	if cols.Length() == 1 {
		match, err := regexp.Match(`^\s*沒有記錄！\s*$`, []byte(cols.Eq(0).Text()))
		if err != nil {
			return fmt.Errorf("Unknown Error on check row is empty in Firewall")
		}

		if match {
			return fmt.Errorf("Table row length is zero")
		}
	}

	cols.Each(fn)

	return nil
}
