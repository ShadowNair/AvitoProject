package user

import (
	"database/sql"
	modeluser "AvitoProject/internal/models/user"
)

const(
	sqlTextExecSetActive = `UPDATE users SET is_active = $1 WHERE id = $2`
	sqlTextQuerySetActive = `
		SELECT id, username, team_name, is_active FROM users WHERE id = $1
	`
	sqlTextGetActive = `
		SELECT id FROM users WHERE team_name = $1 AND is_active = true AND id != $2
	`
	sqlTextGetUser = `SELECT team_name FROM users WHERE id = $1`
)

type DB struct{
	sql *sql.DB
}

func New(sql *sql.DB) *DB{
	return &DB{
		sql: sql,
	}
}

func (s *DB) SetActive(userID modeluser.ID, active bool) (*modeluser.User, error) {
	_, err := s.sql.Exec(sqlTextExecSetActive, active, userID)
	if err != nil {
		return nil, err
	}
	var u modeluser.User
	err = s.sql.QueryRow(sqlTextQuerySetActive, userID).Scan(&u.ID, &u.Username, &u.TeamName, &u.IsActive)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *DB) GetActiveUsersInTeam(team string, excludeID modeluser.ID) ([]modeluser.ID, error) {
	rows, err := s.sql.Query(sqlTextGetActive, team, excludeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []modeluser.ID
	for rows.Next() {
		var id modeluser.ID
		_ = rows.Scan(&id)
		ids = append(ids, id)
	}
	return ids, nil
}

func (s *DB) GetUserTeam(userID modeluser.ID) (string, error) {
	var team string
	err := s.sql.QueryRow(sqlTextGetUser, userID).Scan(&team)
	return team, err
}