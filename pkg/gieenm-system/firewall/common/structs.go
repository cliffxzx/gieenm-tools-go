package common

import (
	"database/sql/driver"
	"fmt"
	"regexp"
	"strconv"

	"github.com/cliffxzx/gieenm-tools/pkg/nusoft-firewall"
	"github.com/cliffxzx/gieenm-tools/pkg/utils"
)

type NusoftID nusoft.NusoftID

// Scan ...
func (n *NusoftID) Scan(value interface{}) error {
	regex := regexp.MustCompile(`^\((?P<Time>.*),(?P<Serial>.*)\)$`)
	res := regex.FindAllStringSubmatch(string(value.([]byte)), -1)
	serial, _ := strconv.ParseInt(res[0][1], 10, 32)
	time, _ := strconv.ParseInt(res[0][2], 10, 32)
	n.Serial = &serial
	n.Time = &time
	return nil
}

// Value ...
func (n NusoftID) Value() (driver.Value, error) {
	return fmt.Sprintf(`(%d,%d)`, *n.Time, *n.Serial), nil
}

// String ...
func (n NusoftID) String() string {
	return fmt.Sprintf(
		"{ Time: %s, Serial: %s  }",
		utils.MustToJSON(n.Time),
		utils.MustToJSON(n.Serial),
	)
}