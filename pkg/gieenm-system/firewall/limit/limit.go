package limit

import (
	"errors"
	"fmt"

	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/database"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/group"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/user"
	"github.com/cliffxzx/gieenm-tools/pkg/utils"
	"github.com/openlyinc/pointy"
)

// Adds require User.ID, Group.ID, RecordMaxCount fields.
func Adds(destSource *[]Limit) error {
	type Input struct {
		UserID         *int `db:"user_id"`
		GroupID        *int `db:"group_id"`
		RecordMaxCount *int `db:"record_max_count"`
	}

	inputs := []Input{}
	for _, ds := range *destSource {
		if ds.User == nil ||
			ds.User.ID == nil ||
			ds.Group == nil ||
			ds.Group.ID == nil ||
			ds.RecordMaxCount == nil {
			return errors.New("all array element require User.ID, Group.ID, RecordMaxCount fields")
		}

		inputs = append(inputs, Input{
			UserID:         ds.User.ID,
			GroupID:        ds.Group.ID,
			RecordMaxCount: ds.RecordMaxCount,
		})
	}

	col, row, arg, err := utils.BulkInsertParameter(inputs)
	if err != nil {
		return err
	}

	sql := fmt.Sprintf(Sqls{}.Adds(), *col, *row)

	converter := []struct {
		Limit
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

	for idx, Limit := range converter {
		(*destSource)[idx] = Limit.Limit
		(*destSource)[idx].Group = &group.Group{ID: Limit.GroupID}
		(*destSource)[idx].User = &user.User{ID: Limit.UserID}
	}

	return nil
}

// GetRecordLimit require User.ID, Group.ID fields
func (u *Limit) Get() error {
	if u == nil || u.User == nil || u.User.ID == nil || u.Group == nil || u.Group.ID == nil {
		return errors.New("User require ID, Group require ID")
	}

	destSource := struct {
		Limit
		UserID  *int `db:"user_id"`
		GroupID *int `db:"group_id"`
	}{}

	err := database.GetDB().Get(&destSource, Sqls{}.Get(), u.User.ID, u.Group.ID)
	if err != nil {
		return err
	}

	*u = destSource.Limit
	u.Group = &group.Group{ID: destSource.GroupID}
	u.User = &user.User{ID: destSource.UserID}

	return nil
}

// GetRecordLimits require UserID field
func GetByUserID(userID int) (*[]Limit, error) {
	sql := Sqls{}.GetByUserID()
	converter := []struct {
		Limit
		UserID  *int `db:"user_id"`
		GroupID *int `db:"group_id"`
	}{}

	err := database.GetDB().Select(&converter, sql, userID)
	if err != nil {
		return nil, err
	}

	destSource := []Limit{}
	for _, ds := range converter {
		ds.Limit.Group = &group.Group{ID: ds.GroupID}
		ds.Limit.User = &user.User{ID: ds.UserID}
		destSource = append(destSource, ds.Limit)
	}

	return &destSource, nil
}

// IsFull if records count return true
//
// required User.ID, Group.ID fields
func (u *Limit) IsFull() (*bool, error) {
	if u == nil || u.User == nil || u.User.ID == nil || u.Group == nil || u.Group.ID == nil {
		return nil, errors.New("required User.ID, Group.ID fields")
	}

	err := u.Get()
	if err != nil {
		return nil, err
	}

	recordsCount := 0
	err = database.GetDB().Get(&recordsCount, Sqls{}.IsFull(), u.User.ID, u.Group.ID)
	if err != nil {
		return nil, err
	}

	switch {
	case recordsCount > *u.RecordMaxCount:
		if *u.RecordMaxCount == -1 {
			return pointy.Bool(false), nil
		} else {
			return nil, LimitExceedError{MaxCount: *u.RecordMaxCount, NowCount: recordsCount}
		}
	case recordsCount == *u.RecordMaxCount:
		return pointy.Bool(true), nil
	case recordsCount < *u.RecordMaxCount:
		return pointy.Bool(false), nil
	}

	panic("not possible!")
}
