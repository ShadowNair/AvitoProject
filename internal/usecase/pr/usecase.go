package pr

import (
	modelpr "AvitoProject/internal/models/pr"
	modeluser "AvitoProject/internal/models/user"
	"AvitoProject/internal/usecase/errors"
	"math/rand"
	"time"
)

type repoPRI interface{
	CreatePR(prID modelpr.ID, name string, authorID modeluser.ID, reviewers []modeluser.ID) error
	GetPR(prID modelpr.ID) (*modelpr.PullRequest, error)
	UpdateReviewers(prID modelpr.ID, reviewers []modeluser.ID) error
	MergePR(prID modelpr.ID) (*modelpr.PullRequest, error)
	GetPRsByReviewer(userID modeluser.ID) ([]modelpr.PullRequestShort, error)
	PRExists(prID modelpr.ID) (bool, error)
}

type repoUserI interface{
	SetActive(userID modeluser.ID, active bool) (*modeluser.User, error)
	GetActiveUsersInTeam(team string, excludeID modeluser.ID) ([]modeluser.ID, error)
	GetUserTeam(userID modeluser.ID) (string, error)
}

type Usecase struct {
	repoPR repoPRI
	repoUser repoUserI
	rand     *rand.Rand
}

func New(repoPR repoPRI, repoUser repoUserI) *Usecase{
	return &Usecase{
		repoPR: repoPR,
		repoUser: repoUser,
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (u *Usecase) CreatePR(prID modelpr.ID, name string, authorID modeluser.ID) (*modelpr.PullRequest, error) {
	exists, err := u.repoPR.PRExists(prID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.ErrPRExists
	}

	teamName, err := u.repoUser.GetUserTeam(authorID)
	if err != nil {
		return nil, errors.ErrNotFound
	}

	candidates, err := u.repoUser.GetActiveUsersInTeam(teamName, authorID)
	if err != nil {
		return nil, err
	}

	u.rand.Shuffle(len(candidates), func(i, j int) { candidates[i], candidates[j] = candidates[j], candidates[i] })
	if len(candidates) > 2 {
		candidates = candidates[:2]
	}

	if err := u.repoPR.CreatePR(prID, name, authorID, candidates); err != nil {
		return nil, err
	}

	return u.repoPR.GetPR(prID)
}

func (u *Usecase) MergePR(prID modelpr.ID) (*modelpr.PullRequest, error) {
	pr, err := u.repoPR.GetPR(prID)
	if err != nil {
		return nil, err
	}
	if pr.Status == "MERGED" {
		return pr, nil 
	}
	return u.repoPR.MergePR(prID)
}

func (u *Usecase) ReassignReviewer(prID modelpr.ID, oldUserID modeluser.ID) (*modelpr.PullRequest, modeluser.ID, error) {
	pr, err := u.repoPR.GetPR(prID)
	if err != nil {
		return nil, "", errors.ErrNotFound 
	}
	if pr.Status == "MERGED" {
		return nil, "", errors.ErrPRMerged
	}

	isAssigned := false
	for _, r := range pr.AssignedReviewers {
		if r == oldUserID {
			isAssigned = true
			break
		}
	}
	if !isAssigned {
		return nil, "", errors.ErrNotAssigned
	}

	teamName, err := u.repoUser.GetUserTeam(oldUserID)
	if err != nil {
		return nil, "", errors.ErrNotFound
	}

	candidates, err := u.repoUser.GetActiveUsersInTeam(teamName, oldUserID)
	if err != nil {
		return nil, "", err
	}
	if len(candidates) == 0 {
		return nil, "", errors.ErrNoCandidate
	}

	newReviewer := candidates[u.rand.Intn(len(candidates))]

	newList := make([]modeluser.ID, 0, len(pr.AssignedReviewers))
	for _, r := range pr.AssignedReviewers {
		if r == oldUserID {
			newList = append(newList, newReviewer)
		} else {
			newList = append(newList, r)
		}
	}

	if err := u.repoPR.UpdateReviewers(prID, newList); err != nil {
		return nil, "", err
	}

	prUpdated, _ := u.repoPR.GetPR(prID)
	return prUpdated, newReviewer, nil
}