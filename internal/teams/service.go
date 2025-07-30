package teams

import (
	"context"
	"github.com/Gurrag09847/weft/internal/database"
	"github.com/Gurrag09847/weft/models"
	"time"
)

var dbTimeout = 3 * time.Second

type Service interface {
	GetTeamsService() (models.TeamSlice, error)
}

type service struct {
	db database.Service
}

func NewService(db database.Service) Service {
	return &service{
		db: db,
	}
}

func (s *service) GetTeamsService() (models.TeamSlice, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	return models.Teams.Query().All(ctx, s.db)
}
