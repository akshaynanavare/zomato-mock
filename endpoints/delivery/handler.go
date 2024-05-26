package delivery

import (
	"errors"
	"net/http"
	"time"

	api "github.com/akshaynanavare/shortest-time/api"

	"github.com/julienschmidt/httprouter"
	"github.com/spf13/cast"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetMinimumTimeToDeliverAll(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	deliveryPartnerID := urlParams.ByName("deliveryPartner")
	if deliveryPartnerID == "" {
		api.Error(w, r, errors.New("invalid delivery partner id"), http.StatusBadRequest)
		return
	}

	ts, path, err := h.service.GetShortestPathOfActiveOrders(deliveryPartnerID)
	if err != nil {
		api.Error(w, r, err, http.StatusInternalServerError)
		return
	}

	api.SuccessJSON(w, r, map[string]interface{}{
		"approx_completion_time": ts,
		"duration":               cast.ToString(time.Until(ts).Minutes()) + " minutes",
		"path":                   path,
	})
}
