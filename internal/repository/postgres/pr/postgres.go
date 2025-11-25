package pr

import (
	"database/sql"
	modelpr "AvitoProject/internal/models/pr"
	modeluser "AvitoProject/internal/models/user"
	"time"
)

type DB struct{
	sql *sql.DB
}

func New(sql *sql.DB) *DB{
	return &DB{
		sql: sql,
	}
}

func (s *DB) CreatePR(prID modelpr.ID, name string, authorID modeluser.ID, reviewers []modeluser.ID) error {
	_, err := s.sql.Exec(`
		INSERT INTO pull_requests (id, name, author_id, status, assigned_reviewers)
		VALUES ($1, $2, $3, 'OPEN', $4)
	`, prID, name, authorID, reviewers)
	return err
}

func (s *DB) GetPR(prID modelpr.ID) (*modelpr.PullRequest, error) {
	var pr modelpr.PullRequest
	var createdAt, mergedAt sql.NullTime
	err := s.sql.QueryRow(`
		SELECT id, name, author_id, status, assigned_reviewers, created_at, merged_at
		FROM pull_requests WHERE id = $1
	`, prID).Scan(
		&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status,
		&pr.AssignedReviewers, &createdAt, &mergedAt,
	)
	if err != nil {
		return nil, err
	}
	if createdAt.Valid { pr.CreatedAt = &createdAt.Time }
	if mergedAt.Valid { pr.MergedAt = &mergedAt.Time }
	return &pr, nil
}

func (s *DB) UpdateReviewers(prID modelpr.ID, reviewers []modeluser.ID) error {
	_, err := s.sql.Exec(`UPDATE pull_requests SET assigned_reviewers = $1 WHERE id = $2`, reviewers, prID)
	return err
}

func (s *DB) MergePR(prID modelpr.ID) (*modelpr.PullRequest, error) {
	now := time.Now()
	_, err := s.sql.Exec(`UPDATE pull_requests SET status = 'MERGED', merged_at = $1 WHERE id = $2`, now, prID)
	if err != nil {
		return nil, err
	}
	return s.GetPR(prID)
}

func (s *DB) GetPRsByReviewer(userID modeluser.ID) ([]modelpr.PullRequestShort, error) {
	rows, err := s.sql.Query(`
		SELECT id, name, author_id, status FROM pull_requests
		WHERE $1 = ANY(assigned_reviewers)
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []modelpr.PullRequestShort
	for rows.Next() {
		var p modelpr.PullRequestShort
		_ = rows.Scan(&p.ID, &p.Name, &p.AuthorID, &p.Status)
		list = append(list, p)
	}
	return list, nil
}

func (s *DB) PRExists(prID modelpr.ID) (bool, error) {
	var exists bool
	err := s.sql.QueryRow(`SELECT EXISTS(SELECT 1 FROM pull_requests WHERE id = $1)`, prID).Scan(&exists)
	return exists, err
}