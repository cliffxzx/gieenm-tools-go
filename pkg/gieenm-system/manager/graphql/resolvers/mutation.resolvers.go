package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/record"
	gql "github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/manager/graphql/generated"
)

func (r *mutationResolver) AddRecords(ctx context.Context, input gql.AddRecordInput) (*record.Record, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SetRecords(ctx context.Context, input gql.SetRecordInput) (*record.Record, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DelRecords(ctx context.Context, input gql.DelRecordInput) (*record.Record, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns gql.MutationResolver implementation.
func (r *Resolver) Mutation() gql.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
