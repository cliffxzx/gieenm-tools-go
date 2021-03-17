package utils

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
)

// SliceToJSON ...
func SliceToJSON(v interface{}) string {
	slice, err := InterfaceToSlice(v)
	if err != nil {
		return "null"
	}

	result := "[\n"

	for _, v := range slice {
		switch v.(type) {
		case fmt.Stringer:
			result += fmt.Sprintf("\t%s", v)
		default:
			result += MustToJSON(v)
		}

		result += ",\n"
	}

	result += "]"

	return result
}

// URLToJSON ...
func URLToJSON(url *url.URL) string {
	if url == nil {
		return "null"
	}

	return fmt.Sprintf(`"%s"`, url.String())
}

// IPAddrToJSON ...
func IPAddrToJSON(ip *net.IPNet) string {
	if ip == nil {
		return "null"
	}

	ones, _ := ip.Mask.Size()

	return fmt.Sprintf(`"%s/%d"`, ip.IP.String(), ones)
}

// MacAddrToJSON ...
func MacAddrToJSON(macAddr *net.HardwareAddr) string {
	if macAddr == nil {
		return "null"
	}

	return fmt.Sprintf(`"%s"`, macAddr.String())
}

// MustToJSON ...
func MustToJSON(v interface{}) string {
	json, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(fmt.Sprintln("can't marshal type 'firewall/Firewall': ", err))
	}

	return string(json)
}
