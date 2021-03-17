package scalars

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
)

// Float64 ...
type Float64 float64

// MarshalFloat64 ...
func MarshalFloat64(f float64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, fmt.Sprintf("%g", f))
	})
}

// UnmarshalFloat64 ...
func UnmarshalFloat64(v interface{}) (float64, error) {
	switch v := v.(type) {
	case string:
		return strconv.ParseFloat(v, 32)
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case float64:
		return v, nil
	case json.Number:
		return strconv.ParseFloat(string(v), 32)
	default:
		return 0, fmt.Errorf("%T is not an float", v)
	}
}
