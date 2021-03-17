package firewall

// Sqls ...
type Sqls struct{}

// GetFirewalls ...
func (Sqls) GetFirewalls() string {
	return `
	select
		id, uid, name, host, username, password, created_at, modified_at
	from
		firewalls
`
}
