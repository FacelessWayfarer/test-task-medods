package handlers

import (
	"encoding/json"
	"errors"
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

	IP, req, err := decodeRefreshRequestBody(*r)
	if err != nil {
		h.logger.Printf("error checking request: %v", err)

		w.WriteHeader(http.StatusBadRequest)

		render.JSON(w, r, RefreshErrResponse{Error: "invalid request"})

		return
	}

	ctx := r.Context()

	request := RequestToRequest(*req)

	resp, err := h.service.UpdateTokens(ctx, request, IP)
	if err != nil {
		if errors.Is(err, models.ErrTokenExpired) {
			h.logger.Printf("error updating tokens: %v", err)

			w.WriteHeader(http.StatusInternalServerError)

			render.JSON(w, r, RefreshErrResponse{Error: ErrTokenExpired.Error()})

			return
		}

		h.logger.Printf("error updating tokens: %v", err)

		w.WriteHeader(http.StatusInternalServerError)

		render.JSON(w, r, RefreshErrResponse{Error: "internal server error"})

		return
	}

	h.logger.Println("successfully refreshed jwt access token")

	w.Header().Set("Content-Type", "application/json")

	render.JSON(w, r, ToRefreshResponse(*resp))
}

func decodeRefreshRequestBody(r http.Request) (string, *RefreshTokensRequest, error) {
	stringIP := strings.Split(r.RemoteAddr, ":")

	IP := stringIP[0]

	var req RefreshTokensRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return "", nil, fmt.Errorf("could not decode request body: %v", err)
	}

	return IP, &req, nil
}
