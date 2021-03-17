package scalars

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
)

// Float ...
type Float float32

// MarshalFloat ...
func MarshalFloat(f float32) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, fmt.Sprintf("%g", f))
	})
}

// UnmarshalFloat ...
func UnmarshalFloat(v interface{}) (float32, error) {
	switch v := v.(type) {
	case string:
		tmp, err := strconv.ParseFloat(v, 32)
		return float32(tmp), err
	case int:
		return float32(v), nil
	case int64:
		return float32(v), nil
	case float32:
		return v, nil
	case json.Number:
		tmp, err := strconv.ParseFloat(string(v), 32)
		return float32(tmp), err
	default:
		return 0, fmt.Errorf("%T is not an float", v)
	}
}
