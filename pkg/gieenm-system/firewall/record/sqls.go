package record

// Sqls ...
type Sqls struct{}

func (Sqls) Adds() string {
	return `
	insert into
		records
		%s
	values
		%s
	on conflict
		(nusoft_id, group_id)
	do nothing
	returning
		id, uid, name, nusoft_id, ip_addr, mac_addr, user_id, group_id, created_at, modified_at
`
}

func (Sqls) Sets() string {
	return `
	insert into
		records
		%s
	values
		%s
	on conflict
		(nusoft_id, firewall_id)
	do update set
		uid = excluded.uid,
		name = excluded.name,
		nusoft_id = excluded.nusoft_id,
		ip_addr = excluded.ip_addr,
		mac_addr = excluded.mac_addr,
		user_id = excluded.user_id,
		group_id = excluded.group_id,
		modified_at = excluded.modified_at
	returning
		id, uid, name, nusoft_id, ip_addr, mac_addr, user_id, group_id, created_at, modified_at
`
}

func (Sqls) Dels() string {
	return `
		delete
		from
			records
		where
			id in %s
		returning
			id, uid, name, nusoft_id, ip_addr, mac_addr, user_id, group_id, created_at, modified_at
`
}

func (Sqls) GetsByUserID() string {
	return `
	select
		id, uid, nusoft_id, name, ip_addr, mac_addr, user_id, group_id, created_at, modified_at
	from
		records
	where
		user_id = $1
`
}

func (Sqls) GetsByUID() string {
	return `
	select
		id, uid, nusoft_id, name, ip_addr, mac_addr, user_id, group_id, created_at, modified_at
	from
		records
	where
		uid = $1
`
}

func (Sqls) GetsByGroupID() string {
	return `
	select
		id, uid, nusoft_id, name, ip_addr, mac_addr, user_id, group_id, created_at, modified_at
	from
		records
	where
		group_id = $1
`
}
