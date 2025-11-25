package teamrepo

import (
	"database/sql"
	modeluser "AvitoProject/internal/models/user"
	modelteam "AvitoProject/internal/models/team"
)

type DB struct{
	sql *sql.DB
}

func New(sql *sql.DB) *DB{
	return &DB{
		sql: sql,
	}
}

func (s *DB) CreateTeam(name string, members []modeluser.User) error {
	tx, _ := s.sql.Begin()
	_, _ = tx.Exec(`INSERT INTO teams (name) VALUES ($1) ON CONFLICT DO NOTHING`, name)
	for _, m := range members {
		_, _ = tx.Exec(`
			INSERT INTO users (id, username, team_name, is_active)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (id) DO UPDATE SET
				username = EXCLUDED.username,
				team_name = EXCLUDED.team_name,
				is_active = EXCLUDED.is_active
		`, m.ID, m.Username, name, m.IsActive)
	}
	return tx.Commit()
}

func (s *DB) GetTeam(name string) (*modelteam.Team, error) {
	rows, err := s.sql.Query(`SELECT id, username, is_active FROM users WHERE team_name = $1`, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	team := &modelteam.Team{Name: name}
	for rows.Next() {
		var u modeluser.User
		if err := rows.Scan(&u.ID, &u.Username, &u.IsActive); err != nil {
			return nil, err
		}
		u.TeamName = name
		team.Members = append(team.Members, u)
	}
	if len(team.Members) == 0 {
		return nil, sql.ErrNoRows
	}
	return team, nil
}