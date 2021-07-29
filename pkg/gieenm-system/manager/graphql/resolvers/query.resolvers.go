package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/announcement"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/group"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/record"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/manager"
	gql "github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/manager/graphql/generated"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/user"
)

func (r *queryResolver) Users(ctx context.Context) ([]user.User, error) {
	return user.Gets()
}

func (r *queryResolver) Records(ctx context.Context) ([]record.Record, error) {
	records, err := record.Gets()
	return *records, err
}

func (r *queryResolver) Groups(ctx context.Context) ([]group.Group, error) {
	groups, err := group.Gets()
	return *groups, err
}

func (r *queryResolver) Firewall(ctx context.Context) ([]manager.Firewall, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Announcement(ctx context.Context) ([]announcement.Announcement, error) {
	panic(fmt.Errorf("not implemented"))
}

// Query returns gql.QueryResolver implementation.
func (r *Resolver) Query() gql.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
