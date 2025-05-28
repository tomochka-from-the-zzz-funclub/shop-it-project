package services

import (
	config "market/internal/cfg"
	database "market/internal/database/postgres"
	"market/internal/models"

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
func (srv *Srv) Create(seller models.Seller) (uuid.UUID, error) {
	uuid, err := srv.db.CreateSeller(seller)
	return uuid, err
}

func (srv *Srv) Update(id uuid.UUID, newSellers models.Seller) error {
	err := srv.db.UpdateSeller(id, newSellers)
	return err
}

func (srv *Srv) Delete(id uuid.UUID) error {
	err := srv.db.DeleteSeller(id)
	return err
}
