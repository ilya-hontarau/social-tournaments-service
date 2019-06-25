package mockdb

import (
	"context"

	"github.com/illfate/social-tournaments-service/pkg/sts"
	"github.com/stretchr/testify/mock"
)

type Connector struct {
	mock.Mock
}

func (c *Connector) AddUser(ctx context.Context, name string) (int64, error) {
	args := c.Called(name)
	return args.Get(0).(int64), args.Error(1)
}

func (c *Connector) GetUser(ctx context.Context, id int64) (*sts.User, error) {
	args := c.Called(id)
	return args.Get(0).(*sts.User), args.Error(1)
}

func (c *Connector) DeleteUser(ctx context.Context, id int64) error {
	args := c.Called(id)
	return args.Error(0)
}

func (c *Connector) AddPoints(ctx context.Context, id, points int64) error {
	args := c.Called(id, points)
	return args.Error(0)
}

func (c *Connector) AddTournament(ctx context.Context, name string, deposit uint64) (int64, error) {
	args := c.Called(name, deposit)
	return args.Get(0).(int64), args.Error(1)
}

func (c *Connector) GetTournament(ctx context.Context, id int64) (*sts.Tournament, error) {
	args := c.Called(id)
	return args.Get(0).(*sts.Tournament), args.Error(1)
}
