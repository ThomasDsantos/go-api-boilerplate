package handlers

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"backend/internal/app/models"
	"backend/internal/database"
)

type HealthHandler struct {
	PgxPool *pgxpool.Pool
}

func NewHealthHandler(db *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{
		PgxPool: db,
	}
}

func (h *HealthHandler) GetHealth(_ context.Context, _ *models.HealthInput) (*models.HealthOutput, error) {
	health, err := database.Health(h.PgxPool)

	if err != nil {
		return nil, huma.Error500InternalServerError("can't get database health")
	}
	dbOk, ok := health["status"]
	if !ok {
		dbOk = "down"
	}

	resp := &models.HealthOutput{}

	resp.Body.Ok = dbOk == "up"
	resp.Body.Database = health
	return resp, nil
}
