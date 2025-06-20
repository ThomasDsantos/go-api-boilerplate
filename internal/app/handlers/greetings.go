package handlers

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/rs/zerolog/log"

	"backend/internal/app/middleware"
	"backend/internal/app/models"
	"backend/internal/store"
)

type GreetingHandler struct {
	Queries *store.Queries
}

func NewGreetingHandler(queries *store.Queries) *GreetingHandler {
	return &GreetingHandler{
		Queries: queries,
	}
}

func (h *GreetingHandler) GetGreeting(
	ctx context.Context,
	input *models.GreetingInput,
) (*models.GreetingOutput, error) {
	if input.Name == "bob" {
		return nil, huma.Error404NotFound("no greeting for bob")
	}

	request, ok := ctx.Value(middleware.CtxRequest{}).(*http.Request)
	if !ok {
		log.Error().
			Msgf("Error: Could not get http.Request from context, key=%v context.t=%v", middleware.CtxRequest{}, request)
		return nil, huma.Error500InternalServerError("Can't retrieve user Ip")
	}
	host, _, err := net.SplitHostPort(request.RemoteAddr)
	if err != nil {
		host = request.RemoteAddr
	}
	res, err2 := h.Queries.InsertVisit(
		ctx,
		store.InsertVisitParams{
			Ip:   host,
			Name: input.Name,
		},
	)
	if err2 != nil {
		log.Error().Msgf("Can't insert visit: %v\n", err2)
		return nil, huma.Error500InternalServerError("Can't insert visit")
	}

	resp := &models.GreetingOutput{}
	resp.Body.Message = fmt.Sprintf("Hello, %s!", res.Name)
	return resp, nil
}
