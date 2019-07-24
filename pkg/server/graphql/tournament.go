package graphql

import (
	"context"

	"github.com/graph-gophers/graphql-go"
	"github.com/illfate/social-tournaments-service/pkg/sts"
	"github.com/pkg/errors"
)

type tournamentArgs struct {
	ID graphql.ID
}

func (r *Resolver) Tournament(ctx context.Context, args tournamentArgs) (*tournamentResolver, error) {
	id, err := decodeID(args.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't decode id [%s]", args.ID)
	}
	t, err := r.s.GetTournament(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't get tournament [%d]", id)
	}
	return &tournamentResolver{
		tournament: *t,
	}, nil
}

type createTournamentsArgs struct {
	Name    string
	Deposit int32
}

func (r *Resolver) CreateTournament(ctx context.Context, args createTournamentsArgs) (*tournamentResolver, error) {
	id, err := r.s.AddTournament(ctx, args.Name, uint64(args.Deposit))
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't add tournament [%s]", args.Name)
	}
	return &tournamentResolver{sts.Tournament{
		ID:      id,
		Name:    args.Name,
		Deposit: uint64(args.Deposit),
	}}, nil
}

type joinTournamentArgs struct {
	ID     graphql.ID
	UserID graphql.ID
}

func (r *Resolver) JoinTournament(ctx context.Context, args joinTournamentArgs) (*tournamentResolver, error) {
	tID, err := decodeID(args.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't decode tournament id [%s]", args.ID)
	}
	userID, err := decodeID(args.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't decode user id [%s]", args.ID)
	}
	err = r.s.JoinTournament(ctx, tID, userID)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't join tournament [%d]", tID)
	}
	result, err := r.Tournament(ctx, tournamentArgs{
		ID: args.ID,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't get tournament [%s]", tID)
	}
	return result, nil
}

type tournamentResolver struct {
	tournament sts.Tournament
}

func (tr *tournamentResolver) ID() graphql.ID {
	return encodeID(tr.tournament.ID)
}

func (tr *tournamentResolver) Name() string {
	return tr.tournament.Name
}

func (tr *tournamentResolver) Deposit() int32 {
	return int32(tr.tournament.Deposit)
}

func (tr *tournamentResolver) Prize() int32 {
	return int32(tr.tournament.Prize)
}

func (tr *tournamentResolver) Winner() *graphql.ID {
	id := encodeID(tr.tournament.Winner)
	return &id
}

func (tr *tournamentResolver) Users() *[]*graphql.ID {
	var idSlice []*graphql.ID
	for _, id := range tr.tournament.Users {
		id := encodeID(id)
		idSlice = append(idSlice, &id)
	}
	return &idSlice
}
