package utils

import (
	"errors"
	"net"
	"net/url"
	"reflect"
)

// InterfaceToSlice ...
func InterfaceToSlice(slice interface{}) ([]interface{}, error) {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		return nil, errors.New("given a non-slice type")
	}

	// Keep the distinction between nil and empty slice input
	if s.IsNil() {
		return nil, errors.New("given a nil value")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret, nil
}

// CheckErr ...
// TODO: Don't panic
func CheckErr(v interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}

	return v
}

// SetQuery ...
func SetQuery(url *url.URL, query map[string]string) *url.URL {
	URLQuery := url.Query()
	for key, val := range query {
		URLQuery.Set(key, val)
	}
	url.RawQuery = URLQuery.Encode()

	return url
}

// ParseSubmask ...
func ParseSubmask(submask string) net.IPMask {
	ip := net.ParseIP(submask)
	if ip == nil {
		return nil
	}

	addr := ip.To4()
	return net.IPv4Mask(addr[0], addr[1], addr[2], addr[3])
}
