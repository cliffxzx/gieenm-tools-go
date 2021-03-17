package group

import (
	"errors"
	"fmt"
	"net"

	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/database"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/common"
	"github.com/cliffxzx/gieenm-tools/pkg/utils"
)

// Adds group to database and return to parameter
//
// Require Name, NusoftID, FirewallID fields.
func Adds(destSource *[]Group) error {
	type Input struct {
		Name       string          `db:"name"`
		NusoftID   common.NusoftID `db:"nusoft_id"`
		FirewallID int             `db:"firewall_id"`
	}

	inputs := []Input{}
	for _, ds := range *destSource {
		if ds.Name == nil ||
			ds.NusoftID == nil ||
			ds.FirewallID == nil {
			return errors.New("all array element require Name, NusoftID, FirewallID fields")
		}

		inputs = append(inputs, Input{
			Name:       *ds.Name,
			NusoftID:   *ds.NusoftID,
			FirewallID: *ds.FirewallID,
		})
	}

	col, row, arg, err := utils.BulkInsertParameter(inputs)
	if err != nil {
		return err
	}

	coverter := []Group{}
	sql := fmt.Sprintf(Sqls{}.Adds(), *col, *row)
	err = database.GetDB().Select(&coverter, sql, arg...)
	if err != nil {
		return err
	}

	for idx, group := range coverter {
		(*destSource)[idx].ID = group.ID
		(*destSource)[idx].UID = group.UID
		(*destSource)[idx].CreatedAt = group.CreatedAt
		(*destSource)[idx].ModifiedAt = group.ModifiedAt
	}

	return nil
}

// Gets groups ...
func Gets() (*[]Group, error) {
	groups := []Group{}
	sql := Sqls{}.Gets()
	err := database.GetDB().Select(&groups, sql)
	if err != nil {
		return nil, err
	}

	return &groups, nil
}

func GetNextIP(g Group) (*net.IPNet, error) {

}
