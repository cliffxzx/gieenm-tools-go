package announcement

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

type AnnounceLevel string

// Announce Level
const (
	INFO      AnnounceLevel = "INFO"
	NOTICE    AnnounceLevel = "NOTICE"
	WARNING   AnnounceLevel = "WARNING"
	EMERGENCY AnnounceLevel = "EMERGENCY"
)

// Scan ...
func (r *AnnounceLevel) Scan(value interface{}) error {
	*r = AnnounceLevel(string(value.([]byte)))
	switch *r {
	case INFO, NOTICE, WARNING, EMERGENCY:
	default:
		return errors.New("unknown enum key")
	}

	return nil
}

// Value ...
func (r AnnounceLevel) Value() (driver.Value, error) {
	return string(r), nil
}

// MarshalAnnounceLevel ...
func MarshalAnnounceLevel(r AnnounceLevel) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(fmt.Sprintf(`"%s"`, r)))
	})
}

// UnmarshalAnnounceLevel ...
func UnmarshalAnnounceLevel(v interface{}) (AnnounceLevel, error) {
	switch v {
	case INFO, NOTICE, WARNING, EMERGENCY:
		announceLebel, ok := v.(AnnounceLevel)
		if !ok {
			return INFO, errors.New("can't parsing to string")
		}

		return announceLebel, nil
	default:
		return INFO, nil
	}
}

type Announcement struct {
	ID         string         `db:"id"`
	UID        string         `db:"uid" gqlgen:"id"`
	Title      *string        `db:"title" gqlgen:"title"`
	Content    *string        `db:"content" gqlgen:"content"`
	Announcer  *string        `db:"announcer" gqlgen:"announcer"`
	CreatedAt  *time.Time     `db:"created_at" gqlgen:"createdAt"`
	ModifiedAt *time.Time     `db:"modified_at" gqlgen:"modifiedAt"`
	Level      *AnnounceLevel `db:"level" gqlgen:"level"`
}

func (Announcement) IsNode() {}
