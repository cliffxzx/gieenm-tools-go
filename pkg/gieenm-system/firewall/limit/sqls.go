package limit

// Sqls ...
type Sqls struct{}

func (Sqls) Adds() string {
	return `
	insert into
		user_record_limits
		%s
	values
		%s
	on conflict
		(user_id, group_id)
	do update set
		id = excluded.id,
		uid = excluded.uid,
		user_id = excluded.user_id,
		group_id = excluded.group_id,
		record_max_count = excluded.record_max_count,
		modified_at = excluded.modified_at
	returning
		id, uid, user_id, group_id, record_max_count, created_at, modified_at
`
}

func (Sqls) Get() string {
	return `
	select
		id, uid, user_id, group_id, record_max_count, created_at, modified_at
	from
		user_record_limits
	where
		user_id = $1
		and
		group_id = $2
`
}

func (Sqls) GetByUserID() string {
	return `
	select
		id, uid, user_id, group_id, record_max_count, created_at, modified_at
	from
		user_record_limits
	where
		user_id = $1
`
}

func (Sqls) IsFull() string {
	return `
	select
		count(id)
	from
		user_record_limits
	where
		user_id = $1
		and
		group_id = $2
`
}
