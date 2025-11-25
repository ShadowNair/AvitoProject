package team

import(
	"AvitoProject/internal/models/user"
)

type Team struct {
	Name    string `json:"team_name"`
	Members []user.User `json:"members"`
}