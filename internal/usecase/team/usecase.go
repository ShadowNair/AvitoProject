package team

import (
	modelteam "AvitoProject/internal/models/team"
	modeluser "AvitoProject/internal/models/user"
	"AvitoProject/internal/usecase/errors"
	"database/sql"
)

type repoTeamI interface {
	CreateTeam(name string, members []modeluser.User) error
	GetTeam(name string) (*modelteam.Team, error)
}

type UseCase struct {
	teamRepo repoTeamI
}

func New(t repoTeamI) *UseCase {
	return &UseCase{
		teamRepo: t,
	}
}	

func (u *UseCase) CreateTeam(t *modelteam.Team) error {
	existing, err := u.teamRepo.GetTeam(t.Name)
	if err == nil && existing != nil {
		return errors.ErrTeamExists
	}

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	return u.teamRepo.CreateTeam(t.Name, t.Members)
}

func (u *UseCase) GetTeam(name string) (*modelteam.Team, error) {
	t, err := u.teamRepo.GetTeam(name)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return t, nil
}