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

func (srv *Srv) SrvCreateGoodCard(card models.GoodCard, sellerID uuid.UUID) (uuid.UUID, error) {
	return srv.db.CreateGoodCard(card, sellerID)
}

func (srv *Srv) SrvDeleteGoodCard(cardId uuid.UUID, sellerID uuid.UUID) error {
	return srv.db.DeleteGoodCard(cardId, sellerID)
}

func (srv *Srv) SrvUpdateGoodCard(id uuid.UUID, card models.GoodCard, sellerID uuid.UUID) error {
	return srv.db.UpdateGoodCard(id, card, sellerID)
}

func (srv *Srv) SrvReadGood(id uuid.UUID, sellerID uuid.UUID) (models.Good, error) {
	return srv.db.ReadGood(id)
}

func (srv *Srv) SrvAddCountGood(id uuid.UUID, number int, sellerID uuid.UUID) (int, error) {
	return srv.db.AddCountGood(id, number, sellerID)
}

func (srv *Srv) SrvDeleteCountGood(id uuid.UUID, number int, sellerID uuid.UUID) (int, error) {
	return srv.db.DeleteCountGood(id, number, sellerID)
}

func (srv *Srv) SrvCreateGood(cardID uuid.UUID, quantity int, sellerID uuid.UUID) error {
	return srv.db.CreateGood(cardID, quantity, sellerID)
}

func (srv *Srv) SrvDeleteGood(goodID uuid.UUID, sellerID uuid.UUID) error {
	return srv.db.DeleteGood(goodID, sellerID)
}

func (s *Srv) SrvSearchGoods(req models.SearchRequest) (*models.SearchResponse, error) {
	return s.db.SearchGoods(req)
}
