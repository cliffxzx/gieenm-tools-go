package group

import (
	"errors"
	"fmt"

	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/base/scalars"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/database"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/common"
	"github.com/cliffxzx/gieenm-tools/pkg/utils"
)

// Adds group to database and return to parameter
//
// Require Name, Subnet, NusoftID, FirewallID fields.
func Adds(destSource *[]Group) error {
	type Input struct {
		Name       string          `db:"name"`
		NusoftID   common.NusoftID `db:"nusoft_id"`
		Subnet     scalars.IPAddr  `db:"subnet"`
		FirewallID int             `db:"firewall_id"`
	}

	inputs := []Input{}
	for _, ds := range *destSource {
		if ds.Name == nil ||
			ds.Subnet == nil ||
			ds.NusoftID == nil ||
			ds.FirewallID == nil {
			return errors.New("all array element require Name, NusoftID, FirewallID fields")
		}

		inputs = append(inputs, Input{
			Name:       *ds.Name,
			Subnet:     *ds.Subnet,
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

// GetByID group ...
func GetByID(ID int) (*Group, error) {
	group := Group{}
	sql := Sqls{}.GetByID()
	err := database.GetDB().Get(&group, sql, ID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New("can't find group by uid")
		}

		return nil, err
	}

	return &group, nil
}

// GetByUID group ...
func GetByUID(UID string) (*Group, error) {
	group := Group{}
	sql := Sqls{}.GetByUID()
	err := database.GetDB().Get(&group, sql, UID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New("can't find group by uid")
		}

		return nil, err
	}

	return &group, nil
}
