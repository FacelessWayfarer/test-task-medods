package handlers

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

const urlPathArg = "user_id"

// ShowAccount godoc
// @Summary      GetTokens
// @Description  Generates access and refresh tokens
// @Tags         Tokens
// @Accept       json
// @Produce      json
// @Param        user_id path string false "user_id"
// @Success      200  {object}  GenResponse
// @Failure      400  {object}  GenErrResponse
// @Failure      500  {object}  GenErrResponse
// @Router       /tokens/{user_id} [GET]
func (h *Handler) GetTokens(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("Checking request")

	userID, ip, err := checkRequest(*r)
	if err != nil {
		h.logger.Printf("error checking request: %v", err)

		w.WriteHeader(http.StatusBadRequest)

		render.JSON(w, r, GenErrResponse{Error: "invalid request"})

		return
	}

	ctx := r.Context()

	resp, err := h.service.GenerateTokens(ctx, userID, ip)
	if err != nil {
		h.logger.Printf("error generating tokens: %v", err)

		w.WriteHeader(http.StatusInternalServerError)

		render.JSON(w, r, GenErrResponse{Error: "internal server error"})

		return
	}

	h.logger.Println("Successfully generated and saved jwt access and refresh tokens")

	w.Header().Set("Content-Type", "application/json")

	render.JSON(w, r, ToGenResponse(*resp))
}

func checkRequest(r http.Request) (string, string, error) {
	userID := chi.URLParam(&r, urlPathArg)
	if userID == "" {
		return "", "", ErrEmptyUserID
	}

	stringIP := strings.Split(r.RemoteAddr, ":")

	ip := stringIP[0]

	return userID, ip, nil
}
