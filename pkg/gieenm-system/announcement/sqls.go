package announcement

// Sqls ...
type Sqls struct{}

func (Sqls) Gets() string {
	return `
	select
		id, uid, title, content, level, announcer, created_at, modified_at
	from
		announcements
`
}

func (Sqls) Adds() string {
	return `
	insert into
		announcements
		%s
	values
		%s
	on conflict
		(uid)
	do nothing
	returning
		id, uid, title, content, level, announcer, created_at, modified_at
`
}
