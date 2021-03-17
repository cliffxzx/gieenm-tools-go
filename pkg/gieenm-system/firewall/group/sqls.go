package group

// Sqls ...
type Sqls struct{}

// Gets ...
func (Sqls) Gets() string {
	return `
	select
		id, uid, name, nusoft_id, firewall_id, created_at, modified_at
	from
		record_groups
`
}

// Adds ...
func (Sqls) Adds() string {
	return `
	insert into
		record_groups
		%s
	values
		%s
	on conflict
		(nusoft_id, firewall_id)
	do update set
		uid = excluded.uid,
		name = excluded.name,
		modified_at = excluded.modified_at
	returning
		id, uid, name, nusoft_id, firewall_id, created_at, modified_at
`
}
