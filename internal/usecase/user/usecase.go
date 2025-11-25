package user

import (
	modeluser "AvitoProject/internal/models/user"
	modelpr "AvitoProject/internal/models/pr"
	"AvitoProject/internal/usecase/errors"
)

type repoUserI interface{
	SetActive(userID modeluser.ID, active bool) (*modeluser.User, error)
	GetActiveUsersInTeam(team string, excludeID modeluser.ID) ([]modeluser.ID, error)
	GetUserTeam(userID modeluser.ID) (string, error)
}

type repoPRI interface{
	GetPRsByReviewer(userID modeluser.ID) ([]modelpr.PullRequestShort, error)
}

type UseCase struct{
	repoUser repoUserI
	repoPR repoPRI
}

func New(repoUser repoUserI, repoPR repoPRI) *UseCase{
	return &UseCase{
		repoUser: repoUser,
		repoPR: repoPR,
	}
}

func (u *UseCase) SetUserActive(userID modeluser.ID, active bool) (*modeluser.User, error) {
	usr, err := u.repoUser.SetActive(userID, active)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return usr, nil
}

func (u *UseCase) GetReviewPRs(userID modeluser.ID) ([]modelpr.PullRequestShort, error) {
	prs, err := u.repoPR.GetPRsByReviewer(userID)
	if err != nil {
		return nil, err
	}
	return prs, nil
}