package postgres

import (
	"AvitoProject/internal/repository/postgres/pr"
	"AvitoProject/internal/repository/postgres/teamrepo"
	"AvitoProject/internal/repository/postgres/user"
	"AvitoProject/internal/connections"
)

type Config struct {
	PrRepo		*pr.DB
	TeamRepo	*teamrepo.DB
	UserRepo	*user.DB
}

func New(connCFG *connections.Config) *Config {
	prRepo := pr.New(connCFG.PostgresSQL)
	teamRepo := teamrepo.New(connCFG.PostgresSQL)
	userRepo := user.New(connCFG.PostgresSQL)
	return &Config{
		PrRepo:	prRepo,
		TeamRepo:	teamRepo,
		UserRepo:	userRepo,
	}
}