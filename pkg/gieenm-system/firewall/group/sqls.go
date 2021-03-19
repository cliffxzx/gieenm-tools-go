package group

// Sqls ...
type Sqls struct{}

// Gets ...
func (Sqls) Gets() string {
	return `
	select
		id, uid, name, subnet, nusoft_id, firewall_id, created_at, modified_at
	from
		record_groups
`
}

// GetByID ...
func (Sqls) GetByID() string {
	return `
	select
		id, uid, name, subnet, nusoft_id, firewall_id, created_at, modified_at
	from
		record_groups
	where
		id = $1
`
}

// GetByUID ...
func (Sqls) GetByUID() string {
	return `
	select
		id, uid, name, subnet, nusoft_id, firewall_id, created_at, modified_at
	from
		record_groups
	where
		uid = $1
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
		subnet = excluded.subnet,
		modified_at = excluded.modified_at
	returning
		id, uid, name, subnet, nusoft_id, firewall_id, created_at, modified_at
`
}
