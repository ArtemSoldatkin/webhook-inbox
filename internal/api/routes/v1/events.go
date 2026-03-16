package routev1

import (
	"errors"
	"net/http"

	mapperv1 "github.com/ArtemSoldatkin/webhook-inbox/internal/api/mapper/v1"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/api/types"
	api "github.com/ArtemSoldatkin/webhook-inbox/internal/api/utils"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

// ListEventsInput defines the expected input parameters for listing events.
type ListEventsInput struct {
	SourceID      int64        `url_param:"source_id,required,min=1"`
	Search        string       `query_param:"search,max_length=512"`
	SortDirection string       `query_param:"sort,allowed=ASC|DESC,default=DESC"`
	PageSize      int          `query_param:"limit,min=1,max=100,default=20"`
	Cursor        types.Cursor `query_param:"cursor"`
}

// listEvents handles GET requests to list all events.
func listEvents(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input, err := api.ParseRequestInput[ListEventsInput](r)
		if err != nil {
			logrus.WithError(err).Error("Failed to parse input parameters")
			http.Error(w, "Invalid input parameters", http.StatusBadRequest)
			return
		}

		logrus.WithFields(logrus.Fields{
			"source_id":      input.SourceID,
			"page_size":      input.PageSize,
			"cursor":         input.Cursor,
			"search":         input.Search,
			"sort_direction": input.SortDirection,
			"query":          r.URL.RawQuery,
		}).Debug("Received listEvents request")

		events, err := svc.ListEvents(
			r.Context(),
			input.SourceID,
			input.Cursor,
			input.PageSize,
			input.Search,
			input.SortDirection,
		)
		if err != nil {
			logrus.WithError(err).Error("Failed to list events")
			http.Error(w, "Failed to list events", http.StatusInternalServerError)
			return
		}

		eventDTOs := mapperv1.ToEventDTOs(events)

		logrus.WithFields(logrus.Fields{
			"source_id":      input.SourceID,
			"returned_count": len(eventDTOs),
		}).Debug("Returning events")

		var nextCursor types.Cursor
		if len(eventDTOs) > input.PageSize {
			lastEvent := eventDTOs[len(eventDTOs)-1]
			nextCursor = types.NewCursor(
				&lastEvent.ReceivedAt,
				&lastEvent.ID,
			)
		}

		paginatedResponse := api.ToPaginatedResponse(
			eventDTOs,
			input.PageSize,
			nextCursor,
		)

		if err := api.JSON(w, http.StatusOK, paginatedResponse); err != nil {
			var writeErr *api.JSONWriteError
			if errors.As(err, &writeErr) {
				logrus.WithError(err).Error("Failed to write response")
			} else {
				logrus.WithError(err).Error("Failed to marshal response")
				http.Error(w, "Failed to list events", http.StatusInternalServerError)
			}
		}
	}
}

// GetEventInput defines the expected input parameters for retrieving a specific event.
type GetEventInput struct {
	EventID int64 `url_param:"event_id,required,min=1"`
}

// getEvent handles GET requests to retrieve an event by its ID.
func getEvent(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input, err := api.ParseRequestInput[GetEventInput](r)
		if err != nil {
			logrus.WithError(err).Error("Failed to parse input parameters")
			http.Error(w, "Invalid input parameters", http.StatusBadRequest)
			return
		}

		logrus.WithFields(logrus.Fields{
			"event_id": input.EventID,
			"query":    r.URL.RawQuery,
		}).Debug("Received getEvent request")

		event, err := svc.GetEventByID(r.Context(), input.EventID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				logrus.WithField("event_id", input.EventID).Info("Event not found")
				http.Error(w, "Event not found", http.StatusNotFound)
				return
			}
			logrus.WithField("event_id", input.EventID).WithError(err).Error("Failed to get event")
			http.Error(w, "Failed to get event", http.StatusInternalServerError)
			return
		}

		eventDTO := mapperv1.ToEventDTO(event)

		if err := api.JSON(w, http.StatusOK, eventDTO); err != nil {
			var writeErr *api.JSONWriteError
			if errors.As(err, &writeErr) {
				logrus.WithError(err).Error("Failed to write response")
			} else {
				logrus.WithError(err).Error("Failed to marshal response")
				http.Error(w, "Failed to get event", http.StatusInternalServerError)
			}
		}
	}
}

// eventsRouter sets up the router for events-related endpoints.
func eventsRouter(svc *service.Service) chi.Router {
	r := chi.NewRouter()
	r.Mount("/{event_id}/delivery-attempts", deliveryAttemptsRouter(svc))
	r.Get("/{event_id}", getEvent(svc))
	r.Get("/", listEvents(svc))
	return r
}
