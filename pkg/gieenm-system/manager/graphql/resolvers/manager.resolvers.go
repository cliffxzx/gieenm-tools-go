package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/group"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/manager"
	gql "github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/manager/graphql/generated"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/user"
)

func (r *firewallResolver) Host(ctx context.Context, obj *manager.Firewall) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *firewallResolver) RecordCount(ctx context.Context, obj *manager.Firewall) (*int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *firewallResolver) PageRowCount(ctx context.Context, obj *manager.Firewall) (*int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *firewallResolver) Groups(ctx context.Context, obj *manager.Firewall) ([]group.Group, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *groupResolver) Users(ctx context.Context, obj *group.Group) ([]user.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *groupResolver) Firewall(ctx context.Context, obj *group.Group) (*manager.Firewall, error) {
	panic(fmt.Errorf("not implemented"))
}

// Firewall returns gql.FirewallResolver implementation.
func (r *Resolver) Firewall() gql.FirewallResolver { return &firewallResolver{r} }

// Group returns gql.GroupResolver implementation.
func (r *Resolver) Group() gql.GroupResolver { return &groupResolver{r} }

type firewallResolver struct{ *Resolver }
type groupResolver struct{ *Resolver }
