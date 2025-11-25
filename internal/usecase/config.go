package usecase

import(
	"AvitoProject/internal/usecase/pr"
	"AvitoProject/internal/usecase/team"
	"AvitoProject/internal/usecase/user"
	"AvitoProject/internal/repository/postgres"
)

type Config struct {
	prUsecase *pr.Usecase  
	teamUsecase *team.UseCase
	userUsecase *user.UseCase
}

func New(configRepo *postgres.Config) *Config {
	PRUsecase := pr.New(configRepo.PrRepo, configRepo.UserRepo)
	UserUsecase := user.New(configRepo.UserRepo, configRepo.PrRepo)
	TeamUsecase := team.New(configRepo.TeamRepo)
	return &Config{
		prUsecase:   PRUsecase,
		userUsecase: UserUsecase,
		teamUsecase: TeamUsecase,
	}
}