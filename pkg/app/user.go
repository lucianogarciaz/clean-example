package app

import (
	"context"
	"fmt"
	"github.com/lucianogarciaz/kit/cqs"
	"github.com/lucianogarciaz/kit/vo"
	"github.com/lucianogarciaz/pulley-example/pkg/domain"
)

type CreateUserCommand struct {
	Name      string
	Email     string
	CompanyID string
}

const createUserCommandName = "create_user"

func (c CreateUserCommand) CommandName() string {
	return createUserCommandName
}

var _ cqs.CommandHandler[CreateUserCommand] = &CreateUser{}

type CreateUser struct {
	repo domain.UserRepository
}

func NewCreateUserCH(repo domain.UserRepository) *CreateUser {
	return &CreateUser{repo: repo}
}

func (c CreateUser) Handle(ctx context.Context, cmd CreateUserCommand) ([]cqs.Event, error) {
	cID, err := vo.ParseID(cmd.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("createUser: parseID: %w", err)
	}

	user, err := domain.New(cmd.Name, cmd.Email, cID)
	if err != nil {
		return nil, fmt.Errorf("createUser: newUser: %w", err)
	}

	err = c.repo.CreateUser(ctx, *user)
	if err != nil {
		return nil, fmt.Errorf("createUser: addUser: %w", err)
	}

	return nil, nil
}
