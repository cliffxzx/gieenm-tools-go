package user

import "fmt"

// SqlsGetType ...
type SqlsGetType uint

// Sqls ...
type Sqls struct {
}

const (
	// ByUserID ...
	ByUserID SqlsGetType = 1 << iota
	// ByEmail ...
	ByEmail SqlsGetType = 1 << iota
	// ByStudentID ...
	ByStudentID SqlsGetType = 1 << iota
)

// Get ...
func (Sqls) Get(ugType SqlsGetType) (string, error) {
	sql := `
select
	id, uid, name, role, email, student_id, password, created_at, modified_at
from
	users
`

	switch ugType {
	case ByEmail, ByStudentID, ByUserID:
		col := map[SqlsGetType]string{
			ByUserID:    "id",
			ByEmail:     "email",
			ByStudentID: "student_id",
		}[ugType]

		sql = fmt.Sprintf(`
%s
WHERE
	%s=$1
`, sql, col)
	case ByEmail + ByStudentID:
		sql = fmt.Sprintf(`
%s
WHERE
	email=$1
	OR
	student_id=$2
`, sql)
	default:
		return "", fmt.Errorf(`Undefined SqlsGetType: "%d"`, ugType)
	}

	return sql, nil
}

// Add ...
func (Sqls) Add() (string, error) {
	return `
insert into
	users
	(name, email, role, student_id, password)
values
	($1, $2, $3, $4, $5)
returning
	id, uid, name, role, email, student_id, password
`, nil
}
