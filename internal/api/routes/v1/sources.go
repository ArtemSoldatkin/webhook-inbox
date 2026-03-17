package routev1

import (
	"errors"
	"net/http"

	mapperv1 "github.com/ArtemSoldatkin/webhook-inbox/internal/api/mapper/v1"
	requestsv1 "github.com/ArtemSoldatkin/webhook-inbox/internal/api/requests/v1"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/api/types"
	api "github.com/ArtemSoldatkin/webhook-inbox/internal/api/utils"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

// listSources handles GET requests to list all sources.
func listSources(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input, err := api.ParseRequestInput[requestsv1.ListSourcesInput](r)
		if err != nil {
			logrus.WithError(err).Error("Failed to parse input parameters")
			http.Error(w, "Invalid input parameters", http.StatusBadRequest)
			return
		}

		logrus.WithFields(logrus.Fields{
			"page_size":      input.PageSize,
			"cursor":         input.Cursor,
			"search":         input.Search,
			"filter_status":  input.Filter,
			"sort_direction": input.SortDirection,
			"query":          r.URL.RawQuery,
		}).Debug("Received listSources request")

		sources, err := svc.ListSources(
			r.Context(),
			input.Cursor,
			input.PageSize,
			input.Search,
			input.Filter,
			input.SortDirection,
		)
		if err != nil {
			logrus.WithError(err).Error("Failed to list sources")
			http.Error(w, "Failed to list sources", http.StatusInternalServerError)
			return
		}

		sourceDTOs := mapperv1.ToSourceDTOs(sources, svc.Config)

		logrus.WithField("returned_count", len(sourceDTOs)).Debug("Returning sources")

		var nextCursor types.Cursor
		if len(sourceDTOs) > input.PageSize {
			lastSource := sourceDTOs[len(sourceDTOs)-1]
			nextCursor = types.NewCursor(
				&lastSource.UpdatedAt,
				&lastSource.ID,
			)
		}

		paginatedResponse := api.ToPaginatedResponse(
			sourceDTOs,
			input.PageSize,
			nextCursor,
		)

		if err := api.JSON(w, http.StatusOK, paginatedResponse); err != nil {
			var writeErr *api.JSONWriteError
			if errors.As(err, &writeErr) {
				logrus.WithError(err).Error("Failed to write response")
			} else {
				logrus.WithError(err).Error("Failed to marshal response")
				http.Error(w, "Failed to list sources", http.StatusInternalServerError)
			}
		}
	}
}

// getSourceByID handles GET requests to retrieve a source by its ID.
func getSourceByID(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input, err := api.ParseRequestInput[requestsv1.GetSourceByIDInput](r)
		if err != nil {
			logrus.WithError(err).Error("Failed to parse input parameters")
			http.Error(w, "Invalid input parameters", http.StatusBadRequest)
			return
		}

		logrus.WithFields(logrus.Fields{
			"source_id": input.SourceID,
			"query":     r.URL.RawQuery,
		}).Debug("Received getSourceByID request")

		source, err := svc.GetSourceByID(r.Context(), input.SourceID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				logrus.WithField("source_id", input.SourceID).Info("Source not found")
				http.Error(w, "Source not found", http.StatusNotFound)
				return
			}
			logrus.WithField("source_id", input.SourceID).WithError(err).Error("Failed to get source")
			http.Error(w, "Failed to get source", http.StatusInternalServerError)
			return
		}

		sourceDTO := mapperv1.ToSourceDTO(source, svc.Config)

		if err := api.JSON(w, http.StatusOK, sourceDTO); err != nil {
			var writeErr *api.JSONWriteError
			if errors.As(err, &writeErr) {
				logrus.WithError(err).Error("Failed to write response")
			} else {
				logrus.WithError(err).Error("Failed to marshal response")
				http.Error(w, "Failed to get source", http.StatusInternalServerError)
			}
		}
	}
}

// createSource handles POST requests to create a new source.
func createSource(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logrus.WithFields(logrus.Fields{
			"query": r.URL.RawQuery,
		}).Debug("Received createSource request")

		input, err := api.ParseRequestInput[requestsv1.CreateSourceInput](r)
		if err != nil {
			logrus.WithError(err).Error("Failed to parse input parameters")
			http.Error(w, "Invalid input parameters", http.StatusBadRequest)
			return
		}

		// Avoid logging full static_headers as it may contain secrets (e.g., Authorization/API keys).
		headerNames := make([]string, 0, len(input.StaticHeaders))
		for k := range input.StaticHeaders {
			headerNames = append(headerNames, k)
		}

		logrus.WithFields(logrus.Fields{
			"egress_url":          input.EgressUrl,
			"static_header_names": headerNames,
			"static_header_count": len(input.StaticHeaders),
			"description":         input.Description,
		}).Debug("Create source request data")

		if err := requestsv1.ValidateCreateSourceInput(input); err != nil {
			logrus.WithError(err).Error("Input validation failed")
			http.Error(w, "Invalid input parameters: "+err.Error(), http.StatusBadRequest)
			return
		}

		source, err := svc.CreateSource(r.Context(), service.CreateSourceInput{
			EgressUrl:     input.EgressUrl,
			StaticHeaders: input.StaticHeaders,
			Description:   input.Description,
		})
		if err != nil {
			logrus.WithError(err).Error("Failed to create source")
			http.Error(w, "Failed to create source", http.StatusInternalServerError)
			return
		}
		logrus.WithField("source_id", source.ID).Info("Created new source")

		sourceDTO := mapperv1.ToSourceDTO(source, svc.Config)

		if err := api.JSON(w, http.StatusCreated, sourceDTO); err != nil {
			var writeErr *api.JSONWriteError
			if errors.As(err, &writeErr) {
				logrus.WithError(err).Error("Failed to write response")
			} else {
				logrus.WithError(err).Error("Failed to marshal response")
				http.Error(w, "Failed to create source", http.StatusInternalServerError)
			}
		}
	}
}

// sourcesRouter sets up the router for sources-related endpoints.
func sourcesRouter(svc *service.Service) chi.Router {
	r := chi.NewRouter()
	r.Get("/", listSources(svc))
	r.Get("/{source_id}", getSourceByID(svc))
	r.Post("/", createSource(svc))
	r.Mount("/{source_id}/events", eventsRouter(svc))
	return r
}
