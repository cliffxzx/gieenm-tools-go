package seeder

// FirewallSeeder ...
func FirewallSeeder() {
	db.MustExec(`
		insert into
			firewalls
			(name, host, username, password)
		values
		` +
		// "('1f', 'http://140.115.151.224:8080', 'network', '3465134651')," +
		// "('2f', 'http://140.115.151.230:8080', 'network', '3465134651')," +
		// "('3f', 'http://140.115.151.228:8080', 'network', '3465134651')," +
		// "('4f', 'http://140.115.151.223:8080', 'network', '3465134651')," +
		// "('205r', 'http://140.115.150.146:8080', 'network', '3465134651')," +
		"('wifi', 'http://140.115.151.237:8080', 'network', '3465134651')")
}
