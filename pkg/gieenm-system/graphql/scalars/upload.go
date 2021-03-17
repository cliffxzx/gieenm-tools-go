package scalars

import (
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"
)

type Upload io.Reader

// MarshalUpload ...
func MarshalUpload(f Upload) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.Copy(w, f)
	})
}

// UnmarshalUpload ...
func UnmarshalUpload(v interface{}) (Upload, error) {
	upload, ok := v.(Upload)
	if !ok {
		return upload, fmt.Errorf("%T is not an Upload", v)
	}
	return upload, nil
}
