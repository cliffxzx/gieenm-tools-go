package scalars

import (
	"database/sql/driver"
	"fmt"
	"io"
	"net"

	"github.com/99designs/gqlgen/graphql"
)

// IPAddr ...
type IPAddr net.IPNet

// MarshalIPAddr ...
func MarshalIPAddr(f *IPAddr) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		if f.IP == nil || f.Mask == nil {
			w.Write([]byte("null"))
		} else {
			tmp := net.IPNet(*f)
			w.Write([]byte(fmt.Sprintf("\"%s\"", (&tmp).String())))
		}
	})
}

// UnmarshalIPAddr ...
func UnmarshalIPAddr(v interface{}) (*IPAddr, error) {
	switch v := v.(type) {
	case string:
		ip, ipNet, err := net.ParseCIDR(v)
		ipNet.IP = ip
		tmp := IPAddr(*ipNet)
		return &tmp, err
	default:
		return &IPAddr{}, fmt.Errorf("%T is not an IPNet", v)
	}
}

// ToStd ...
func (i *IPAddr) ToStd() *net.IPNet {
	if i == nil {
		return nil
	}

	ip := net.IPNet(*i)

	return &ip
}

// Scan ...
func (i *IPAddr) Scan(value interface{}) error {
	ip, ipNet, _ := net.ParseCIDR(string(value.([]byte)))
	ipNet.IP = ip
	*i = IPAddr(*ipNet)
	return nil
}

// Value ...
func (i IPAddr) Value() (driver.Value, error) {
	ones, _ := i.Mask.Size()
	return fmt.Sprintf(`%s/%d`, i.IP.String(), ones), nil
}
