package record

import (
	"errors"
	"fmt"

	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/base/scalars"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/database"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/common"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/group"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/limit"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/user"
	"github.com/cliffxzx/gieenm-tools/pkg/utils"
)

// Adds record to database and return to parameter
//
// Require Name, NusoftID, IPAddr, MacAddr, User.ID, Group.ID fields.
func Adds(destSource *[]Record) error {
	type Input struct {
		Name     string          `db:"name"`
		NusoftID common.NusoftID `db:"nusoft_id"`
		IPAddr   scalars.IPAddr  `db:"ip_addr"`
		MacAddr  scalars.MacAddr `db:"mac_addr"`
		User     int             `db:"user_id"`
		Group    int             `db:"group_id"`
	}

	inputs := []Input{}
	for _, ds := range *destSource {
		if ds.Name == nil ||
			ds.NusoftID == nil ||
			ds.IPAddr == nil ||
			ds.MacAddr == nil ||
			ds.User == nil ||
			ds.User.ID == nil ||
			ds.Group == nil ||
			ds.Group.ID == nil {
			return errors.New("all array element require Name, NusoftID, IPAddr, MacAddr, User.ID, Group.ID fields")
		}

		l := limit.Limit{User: ds.User, Group: ds.Group}
		isFull, err := l.IsFull()
		if err != nil {
			return err
		}

		if *isFull {
			return errors.New("one of record user's records is full")
		}

		inputs = append(inputs, Input{
			Name:     *ds.Name,
			NusoftID: *ds.NusoftID,
			IPAddr:   *ds.IPAddr,
			MacAddr:  *ds.MacAddr,
			User:     *ds.User.ID,
			Group:    *ds.Group.ID,
		})
	}

	col, row, arg, err := utils.BulkInsertParameter(inputs)
	if err != nil {
		return err
	}

	sql := fmt.Sprintf(Sqls{}.Adds(), *col, *row)

	converter := []struct {
		Record
		UserID  *int `db:"user_id"`
		GroupID *int `db:"group_id"`
	}{}

	err = database.GetDB().Select(&converter, sql, arg...)
	if err != nil {
		return err
	}

	if len(converter) != len(*destSource) {
		return errors.New("input count not equals result count")
	}

	for idx, record := range converter {
		(*destSource)[idx] = record.Record
		(*destSource)[idx].ID = record.Record.ID
		(*destSource)[idx].UID = record.Record.UID
		(*destSource)[idx].CreatedAt = record.Record.CreatedAt
		(*destSource)[idx].ModifiedAt = record.Record.ModifiedAt
		(*destSource)[idx].User = &user.User{ID: record.UserID}
		(*destSource)[idx].Group = &group.Group{ID: record.GroupID}
		(*destSource)[idx].User = &user.User{ID: record.UserID}
	}

	return nil
}

// Sets record to database and return to parameter
//
// Require UID, Name, NusoftID, IPAddr, MacAddr, User.ID, Group.ID fields.
func Sets(destSource *[]Record) error {
	type Input struct {
		UID      string          `db:"uid"`
		Name     string          `db:"name"`
		NusoftID common.NusoftID `db:"nusoft_id"`
		IPAddr   scalars.IPAddr  `db:"ip_addr"`
		MacAddr  scalars.MacAddr `db:"mac_addr"`
		User     int             `db:"user_id"`
		Group    int             `db:"group_id"`
	}

	inputs := []Input{}
	for _, ds := range *destSource {
		if ds.UID == nil ||
			ds.Name == nil ||
			ds.NusoftID == nil ||
			ds.IPAddr == nil ||
			ds.MacAddr == nil ||
			ds.User == nil ||
			ds.User.ID == nil ||
			ds.Group == nil ||
			ds.Group.ID == nil {
			return errors.New("all array element require UID, Name, NusoftID, IPAddr, MacAddr, User.ID, Group.ID fields")
		}

		inputs = append(inputs, Input{
			UID:      *ds.UID,
			Name:     *ds.Name,
			NusoftID: *ds.NusoftID,
			IPAddr:   *ds.IPAddr,
			MacAddr:  *ds.MacAddr,
			User:     *ds.User.ID,
			Group:    *ds.Group.ID,
		})
	}

	col, row, arg, err := utils.BulkInsertParameter(inputs)
	if err != nil {
		return err
	}

	sql := fmt.Sprintf(Sqls{}.Sets(), *col, *row)

	converter := []struct {
		Record
		UserID  *int `db:"user_id"`
		GroupID *int `db:"group_id"`
	}{}

	err = database.GetDB().Select(&converter, sql, arg...)
	if err != nil {
		return err
	}

	if len(converter) != len(*destSource) {
		return errors.New("input count not equals result count")
	}

	for idx, record := range converter {
		(*destSource)[idx] = record.Record
		(*destSource)[idx].Group = &group.Group{ID: record.GroupID}
		(*destSource)[idx].User = &user.User{ID: record.UserID}
	}

	return nil
}

// Dels record to database and return to parameter
//
// Require UID fields.
func Dels(destSource *[]Record) error {
	inputs := "("
	for _, ds := range *destSource {
		if ds.UID == nil {
			return errors.New("all array element require UID fields")
		}

		inputs = fmt.Sprintf("%s'%s',", inputs, *ds.UID)
	}
	inputs = fmt.Sprintf("%s)", inputs[:len(inputs)-1])

	sql := fmt.Sprintf(Sqls{}.Dels(), inputs)

	converter := []struct {
		Record
		UserID  *int `db:"user_id"`
		GroupID *int `db:"group_id"`
	}{}

	err := database.GetDB().Select(&converter, sql)
	if err != nil {
		return err
	}

	if len(converter) != len(*destSource) {
		return errors.New("input count not equals result count")
	}

	for idx, record := range converter {
		(*destSource)[idx] = record.Record
		(*destSource)[idx].Group = &group.Group{ID: record.GroupID}
		(*destSource)[idx].User = &user.User{ID: record.UserID}
	}

	return nil
}

// Gets records
func Gets() (*[]Record, error) {
	records := []Record{}
	converter := []struct {
		Record
		UserID  *int `db:"user_id"`
		GroupID *int `db:"group_id"`
	}{}

	err := database.GetDB().Select(&converter, Sqls{}.Gets())
	if err != nil {
		return nil, err
	}

	for _, r := range converter {
		r.Record.Group = &group.Group{ID: r.GroupID}
		r.Record.User = &user.User{ID: r.UserID}
		records = append(records, r.Record)
	}

	return &records, nil
}

// GetsByUserID record to database and return to parameter
//
// Require userID fields.
func GetsByUserID(userID int) (*[]Record, error) {
	records := []Record{}
	converter := []struct {
		Record
		UserID  *int `db:"user_id"`
		GroupID *int `db:"group_id"`
	}{}

	err := database.GetDB().Select(&converter, Sqls{}.GetsByUserID(), userID)
	if err != nil {
		return nil, err
	}

	for _, r := range converter {
		r.Record.Group = &group.Group{ID: r.GroupID}
		r.Record.User = &user.User{ID: r.UserID}
		records = append(records, r.Record)
	}

	return &records, nil
}

// GetsByUID record to database and return to parameter
//
// Require UID fields.
func GetsByUID(UID string) (*Record, error) {
	converter := struct {
		Record
		UserID  *int `db:"user_id"`
		GroupID *int `db:"group_id"`
	}{}

	err := database.GetDB().Get(&converter, Sqls{}.GetsByUID(), UID)
	if err != nil {
		return nil, err
	}

	converter.Record.Group = &group.Group{ID: converter.GroupID}
	converter.Record.User = &user.User{ID: converter.UserID}

	return &converter.Record, nil
}

// GetsByGroupID record to database and return to parameter
//
// Require GroupID fields.
func GetsByGroupID(groupID int) (*[]Record, error) {
	records := []Record{}
	converter := []struct {
		Record
		UserID  *int `db:"user_id"`
		GroupID *int `db:"group_id"`
	}{}

	err := database.GetDB().Select(&converter, Sqls{}.GetsByGroupID(), groupID)
	if err != nil {
		return nil, err
	}

	for _, r := range converter {
		r.Record.Group = &group.Group{ID: r.GroupID}
		r.Record.User = &user.User{ID: r.UserID}
		records = append(records, r.Record)
	}

	return &records, nil
}
