package nusoft

import (
	"net"
	"time"
)

// AnFlow ...
type AnFlow struct {
	Name    string    ``
	IPAddr  net.IPNet ``
	MacAddr string    ``
	Date    time.Time ``
}

// AllAnFlow ...
func (*Firewall) AllAnFlow() {
}

// DropAnFlow ...
func (*Firewall) DropAnFlow() {
}
