package announcement

import (
	"errors"
	"fmt"

	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/database"
	"github.com/cliffxzx/gieenm-tools/pkg/utils"
)

// Gets
func Gets() (*[]Announcement, error) {
	announcements := []Announcement{}
	sql := Sqls{}.Gets()
	err := database.GetDB().Select(&announcements, sql)
	if err != nil {
		return nil, err
	}

	return &announcements, nil
}

// Adds announcement to database and return to parameter
//
// Require Title, Content, Annnouncer, Level fields.
func Adds(destSource *[]Announcement) error {
	type Input struct {
		Title     string        `db:"title"`
		Content   string        `db:"content"`
		Announcer string        `db:"announcer"`
		Level     AnnounceLevel `db:"level"`
	}

	inputs := []Input{}
	for _, ds := range *destSource {
		if ds.Title == nil ||
			ds.Content == nil ||
			ds.Level == nil ||
			ds.Announcer == nil {
			return errors.New("all array element require Title, Content, Annnouncer, Level fields")
		}

		inputs = append(inputs, Input{
			Title:     *ds.Title,
			Content:   *ds.Content,
			Announcer: *ds.Announcer,
			Level:     *ds.Level,
		})
	}

	col, row, arg, err := utils.BulkInsertParameter(inputs)
	if err != nil {
		return err
	}

	coverter := []Announcement{}
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
