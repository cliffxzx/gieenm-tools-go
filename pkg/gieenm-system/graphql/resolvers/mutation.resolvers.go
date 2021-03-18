package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/authentication"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/group"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/record"
	gql "github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/graphql"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/user"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/validator"
	"github.com/cliffxzx/gieenm-tools/pkg/utils"
)

func (r *mutationResolver) Login(ctx context.Context, input user.LoginInput) (*user.AuthPayload, error) {
	validator := validator.Get(ctx)

	err := validator.Struct(input)
	if err != nil {
		return nil, err
	}

	login, token, err := input.ToUser().Login()
	if err != nil {
		return nil, err
	}

	tokenStr, err := token.Value()
	if err != nil {
		return nil, err
	}

	return &user.AuthPayload{User: login, Token: &tokenStr}, nil
}

func (r *mutationResolver) Logout(ctx context.Context) (*bool, error) {
	gCtx, err := utils.GetGinContext(ctx)
	if err != nil {
		return nil, nil
	}

	token, err := authentication.VerifyToken(gCtx)
	if err != nil {
		return nil, err
	}

	err = authentication.DelAuth(token)
	if err != nil {
		return nil, err
	}

	success := true

	return &success, nil
}

func (r *mutationResolver) Register(ctx context.Context, input user.RegisterInput) (*user.AuthPayload, error) {
	validator := validator.Get(ctx)

	err := validator.Struct(input)
	if err != nil {
		return nil, err
	}

	register := input.ToUser()
	register.Role = &user.STUDENT

	register, token, err := register.Register()
	if err != nil {
		return nil, err
	}

	tokenStr, err := token.Value()
	if err != nil {
		return nil, err
	}

	return &user.AuthPayload{User: register, Token: &tokenStr}, nil
}

func (r *mutationResolver) AddRecords(ctx context.Context, input gql.AddRecordInput) (*record.Record, error) {
	gCtx, err := utils.GetGinContext(ctx)
	if err != nil {
		return nil, err
	}

	u, err := user.GetByHeaderController(gCtx)
	if err != nil {
		return nil, err
	}

	rd := record.Record{
		Name:    &input.Name,
		MacAddr: input.MacAddr,
		Group:   &group.Group{UID: &input.GroupID},
		User:    u,
	}

	return firewall.AddRecordsController(&rd)
}

func (r *mutationResolver) SetRecords(ctx context.Context, input gql.SetRecordInput) (*record.Record, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DelRecords(ctx context.Context, input gql.DelRecordInput) (*record.Record, error) {
	gCtx, err := utils.GetGinContext(ctx)
	if err != nil {
		return nil, err
	}

	u, err := user.GetByHeaderController(gCtx)
	if err != nil {
		return nil, err
	}

	rd := record.Record{
		UID:  &input.ID,
		User: u,
	}

	return firewall.DelRecordsController(&rd)
}

// Mutation returns gql.MutationResolver implementation.
func (r *Resolver) Mutation() gql.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
