package scalars

import (
	"database/sql/driver"
	"fmt"
	"io"
	"net"

	"github.com/99designs/gqlgen/graphql"
)

// MacAddr ...
type MacAddr net.HardwareAddr

// MarshalMacAddr ...
func MarshalMacAddr(f *MacAddr) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		if f == nil {
			w.Write([]byte("null"))
		} else {
			w.Write([]byte(fmt.Sprintf("\"%s\"", net.HardwareAddr(*f).String())))
		}
	})
}

// UnmarshalMacAddr ...
func UnmarshalMacAddr(v interface{}) (*MacAddr, error) {
	switch v := v.(type) {
	case string:
		mac, err := net.ParseMAC(v)
		tmp := MacAddr(mac)
		return &tmp, err
	default:
		return &MacAddr{}, fmt.Errorf("%T is not an MacAddr", v)
	}
}

// ToStd ...
func (m *MacAddr) ToStd() *net.HardwareAddr {
	if m == nil {
		return nil
	}

	mac := net.HardwareAddr(*m)

	return &mac
}

// Scan ...
func (m *MacAddr) Scan(value interface{}) error {
	mac, _ := net.ParseMAC(string(value.([]byte)))
	*m = MacAddr(mac)
	return nil
}

// Value ...
func (m MacAddr) Value() (driver.Value, error) {
	return net.HardwareAddr(m).String(), nil
}
