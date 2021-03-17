package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	gql "github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/graphql"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/viewer"
)

func (r *queryResolver) Viewer(ctx context.Context) (*viewer.Viewer, error) {
	return viewer.Controller(ctx)
}

// Query returns gql.QueryResolver implementation.
func (r *Resolver) Query() gql.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
