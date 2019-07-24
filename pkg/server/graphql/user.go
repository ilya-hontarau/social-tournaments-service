package graphql

import (
	"context"

	"github.com/graph-gophers/graphql-go"
	"github.com/illfate/social-tournaments-service/pkg/sts"
	"github.com/pkg/errors"
)

type userArgs struct {
	ID graphql.ID
}

func (r *Resolver) User(ctx context.Context, args userArgs) (*userResolver, error) {
	id, err := decodeID(args.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't decode id [%s]", args.ID)
	}
	user, err := r.s.GetUser(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't get user [%d]", id)
	}
	return &userResolver{
		user: *user,
	}, nil
}

type createUserArgs struct {
	Name string
}

func (r *Resolver) CreateUser(ctx context.Context, args createUserArgs) (*userResolver, error) {
	id, err := r.s.AddUser(ctx, args.Name)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't add user [%s]", args.Name)
	}
	return &userResolver{sts.User{
		ID:      id,
		Name:    args.Name,
		Balance: 0,
	}}, nil
}

func (r *Resolver) DeleteUser(ctx context.Context, args userArgs) (*graphql.ID, error) {
	id, err := decodeID(args.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't decode id [%s]", args.ID)
	}
	err = r.s.DeleteUser(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't delete user [%d]", id)
	}
	return &args.ID, nil
}

type userPointsArgs struct {
	ID     graphql.ID
	Points int32
}

func (r *Resolver) TakeUserPoints(ctx context.Context, args userPointsArgs) (*userResolver, error) {
	id, err := decodeID(args.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't decode id [%s]", args.ID)
	}
	err = r.s.AddPoints(ctx, id, int64(-args.Points))
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't take points from user [%s]", id)
	}

	result, err := r.User(ctx, userArgs{
		ID: args.ID,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't get user [%d]", id)
	}
	return result, nil
}

func (r *Resolver) AddUserPoints(ctx context.Context, args userPointsArgs) (*userResolver, error) {
	id, err := decodeID(args.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't decode id [%s]", args.ID)
	}
	err = r.s.AddPoints(ctx, id, int64(args.Points))
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't add points to user [%s]", id)
	}

	result, err := r.User(ctx, userArgs{
		ID: args.ID,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't get user [%d]", id)
	}
	return result, nil
}

type userResolver struct {
	user sts.User
}

func (ur *userResolver) ID() graphql.ID {
	return encodeID(ur.user.ID)
}

func (ur *userResolver) Name() string {
	return ur.user.Name
}

func (ur *userResolver) Balance() int32 {
	return int32(ur.user.Balance)
}
