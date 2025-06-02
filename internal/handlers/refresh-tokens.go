package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/render"

	"github.com/FacelessWayfarer/test-task-medods/internal/service/models"
)

// ShowAccount godoc
// @Summary      RefreshTokens
// @Description  Refreshes access and refresh tokens
// @Tags         Tokens
// @Accept       json
// @Produce      json
// @Param        tokens body RefreshTokensRequest true "access and refresh tokens"
// @Success      200  {object}  RefreshResponse
// @Failure      400  {object}  RefreshErrResponse
// @Failure      500  {object}  RefreshErrResponse
// @Router       /tokens/ [POST]
func (h *Handler) RefreshTokens(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("refreshing tokens")

	ctx, IP, req, err := decodeRefreshRequestBody(*r)
	if err != nil {
		h.logger.Printf("error checking request: %v", err)

		w.WriteHeader(http.StatusBadRequest)

		render.JSON(w, r, RefreshErrResponse{Error: "invalid request"})

		return
	}

	resp, err := h.service.UpdateTokens(ctx, models.TokensToRefresh{
		AccessToken:        req.AccessToken,
		Base64RefreshToken: req.Base64RefreshToken,
	}, IP)
	if err != nil {
		h.logger.Printf("error updating tokens: %v", err)

		w.WriteHeader(http.StatusInternalServerError)

		render.JSON(w, r, RefreshErrResponse{Error: "internal server error"})

		return
	}

	h.logger.Println("successfully refreshed jwt access token")

	w.Header().Set("Content-Type", "application/json")

	render.JSON(w, r, ToRefreshResponse(*resp))
}

func decodeRefreshRequestBody(r http.Request) (context.Context, string, *RefreshTokensRequest, error) {
	ctx := r.Context()

	stringIP := strings.Split(r.RemoteAddr, ":")

	IP := stringIP[0]

	var req RefreshTokensRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return context.Background(), "", nil, fmt.Errorf("could not decode request body: %v", err)
	}

	return ctx, IP, &req, nil
}
