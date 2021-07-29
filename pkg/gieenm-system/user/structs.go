package user

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

// Role ...
type Role string

//User ...
type User struct {
	ID         *int       `db:"id"`
	UID        *string    `db:"uid"                 gqlgen:"id"`
	Role       *Role      `db:"role,type:role_enum" gqlgen:"role"`
	Name       *string    `db:"name"                gqlgen:"name"`
	Email      *string    `db:"email"               gqlgen:"email"`
	StudentID  *string    `db:"student_id"          gqlgen:"studentID"`
	Password   *string    `db:"password"            gqlgen:"password"`
	CreatedAt  *time.Time `db:"created_at"          gqlgen:"createdAt"`
	ModifiedAt *time.Time `db:"modified_at"         gqlgen:"modifiedAt"`
}

// IsNode ...
func (User) IsNode() {}

// User Role
var (
	GUEST   Role = "GUEST"
	STUDENT Role = "STUDENT"
	TEACHER Role = "TEACHER"
	MANAGER Role = "MANAGER"
)

// Scan ...
func (r *Role) Scan(value interface{}) error {
	*r = Role(string(value.([]byte)))
	switch *r {
	case GUEST, STUDENT, TEACHER, MANAGER:
	default:
		*r = GUEST
	}

	return nil
}

// Value ...
func (r Role) Value() (driver.Value, error) {
	return string(r), nil
}

// MarshalRole ...
func MarshalRole(r Role) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(fmt.Sprintf(`"%s"`, r)))
	})
}

// UnmarshalRole ...
func UnmarshalRole(v interface{}) (Role, error) {
	switch v {
	case GUEST, STUDENT, TEACHER, MANAGER:
		role, ok := v.(Role)
		if !ok {
			return GUEST, errors.New("Can't parsing to string")
		}

		return role, nil
	default:
		return GUEST, nil
	}
}
