package services

import (
	config "goods/internal/cfg"
	database "goods/internal/database/postgres"
	"goods/internal/models"

	"github.com/google/uuid"
)

type Srv struct {
	db database.InterfacePostgresDB
}

func NewSrv(cfg config.Config) *Srv {
	base := database.NewPostgres(cfg)
	return &Srv{
		db: base,
	}
}

func (srv *Srv) SrvCreateGoodCard(card models.GoodCard) (uuid.UUID, error) {
	return srv.db.CreateGoodCard(card)
}

func (srv *Srv) SrvDeleteGoodCard(cardId uuid.UUID) error {
	return srv.db.DeleteGoodCard(cardId)
}

func (srv *Srv) SrvUpdateGoodCard(id uuid.UUID, card models.GoodCard) error {
	return srv.db.UpdateGoodCard(id, card)
}

func (srv *Srv) SrvReadGood(id uuid.UUID) (models.Good, error) {
	return srv.db.ReadGood(id)
}

func (srv *Srv) SrvAddCountGood(id uuid.UUID, number int) (int, error) {
	return srv.db.AddCountGood(id, number)
}

func (srv *Srv) SrvDeleteCountGood(id uuid.UUID, number int) (int, error) {
	return srv.db.DeleteCountGood(id, number)
}

func (srv *Srv) SrvCreateGood(cardID uuid.UUID, quantity int) error {
	return srv.db.CreateGood(cardID, quantity)
}

func (srv *Srv) SrvDeleteGood(goodID uuid.UUID) error {
	return srv.db.DeleteGood(goodID)
}
