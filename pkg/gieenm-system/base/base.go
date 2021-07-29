package gql

// Node ...
type Node interface {
	IsNode()
}

// PageInfo ...
type PageInfo struct {
	HasNextPage     bool    `gqlgen:"hasNextPage"`
	HasPreviousPage bool    `gqlgen:"hasPreviousPage"`
	StartCursor     *string `gqlgen:"startCursor"`
	EndCursor       *string `gqlgen:"endCursor"`
}

// PageInput ...
type PageInput struct {
	Cursor *string `gqlgen:"cursor"`
	Count  *int    `gqlgen:"count"`

	Above *bool `gqlgen:"above"`
}
