package routev1

import (
	"encoding/json"
	"net/http"
	"strconv"

	dtov1 "github.com/ArtemSoldatkin/webhook-inbox/internal/api/dto/v1"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

// listDeliveryAttempts handles GET requests to list all delivery attempts.
func listDeliveryAttempts(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventIDRaw := chi.URLParam(r, "eventID")
		eventID, err := strconv.ParseInt(eventIDRaw, 10, 64)
		if err != nil {
			logrus.WithError(err).Error("Invalid event ID")
			http.Error(w, "Invalid event ID", http.StatusBadRequest)
			return
		}
		deliveryAttempts, err := svc.ListDeliveryAttempts(r.Context(), eventID)
		if err != nil {
			logrus.WithError(err).Error("Failed to list delivery attempts")
			http.Error(w, "Failed to list delivery attempts", http.StatusInternalServerError)
			return
		}
		deliveryAttemptsDTO := make([]dtov1.DeliveryAttemptDTO, len(deliveryAttempts))
		for i, deliveryAttempt := range deliveryAttempts {
			deliveryAttemptsDTO[i] = dtov1.DeliveryAttemptDTO{
				ID: 		  	deliveryAttempt.ID,
				EventID:       	deliveryAttempt.EventID,
				AttemptNumber: 	deliveryAttempt.AttemptNumber,
				State:         	deliveryAttempt.State,
				StatusCode: 	utils.PtrIfValid(deliveryAttempt.StatusCode.Int32, deliveryAttempt.StatusCode.Valid),
				ErrorType:     	utils.PtrIfValid(deliveryAttempt.ErrorType.String, deliveryAttempt.ErrorType.Valid),
				ErrorMessage: 	utils.PtrIfValid(deliveryAttempt.ErrorMessage.String,deliveryAttempt.ErrorMessage.Valid),
				StartedAt:     	utils.PtrIfValid(deliveryAttempt.StartedAt.Time,deliveryAttempt.StartedAt.Valid),
				FinishedAt:    	utils.PtrIfValid(deliveryAttempt.FinishedAt.Time,deliveryAttempt.FinishedAt.Valid),
				CreatedAt:     	deliveryAttempt.CreatedAt.Time,
				NextAttemptAt: 	utils.PtrIfValid(deliveryAttempt.NextAttemptAt.Time,deliveryAttempt.NextAttemptAt.Valid),
			}
		}
		response, err := json.Marshal(deliveryAttemptsDTO)
		if err != nil {
			logrus.WithError(err).Error("Failed to marshal delivery attempts")
			http.Error(w, "Failed to list delivery attempts", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

// deliveryAttemptsRouter sets up the router for delivery attempts-related endpoints.
func deliveryAttemptsRouter(svc *service.Service) chi.Router {
	r := chi.NewRouter()
	r.Get("/", listDeliveryAttempts(svc))
	return r
}
