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
			if err := api.JSON(
				w,
				http.StatusBadRequest,
				map[string]string{"error": "Invalid input parameters"},
			); err != nil {
				logrus.WithError(err).Error("Failed to write error response")
			}
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
			if err := api.JSON(
				w,
				http.StatusInternalServerError,
				map[string]string{"error": "Failed to list sources"},
			); err != nil {
				logrus.WithError(err).Error("Failed to write error response")
			}
			return
		}

		sourceDTOs := mapperv1.ToSourceDTOs(sources, svc.Config)

		logrus.WithField("returned_count", len(sourceDTOs)).Debug("Returning sources")

		var nextCursor types.Cursor
		if len(sourceDTOs) > input.PageSize {
			lastSource := sourceDTOs[input.PageSize-1]
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
				return
			}

			logrus.WithError(err).Error("Failed to marshal response")
			if err := api.JSON(
				w,
				http.StatusInternalServerError,
				map[string]string{"error": "Failed to list sources"},
			); err != nil {
				logrus.WithError(err).Error("Failed to write error response")
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
			if err := api.JSON(
				w,
				http.StatusBadRequest,
				map[string]string{"error": "Invalid input parameters"},
			); err != nil {
				logrus.WithError(err).Error("Failed to write error response")
			}
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
				if err := api.JSON(
					w,
					http.StatusNotFound,
					map[string]string{"error": "Source not found"},
				); err != nil {
					logrus.WithError(err).Error("Failed to write error response")
				}
				return
			}
			logrus.WithField("source_id", input.SourceID).WithError(err).Error("Failed to get source")
			if err := api.JSON(
				w,
				http.StatusInternalServerError,
				map[string]string{"error": "Failed to get source"},
			); err != nil {
				logrus.WithError(err).Error("Failed to write error response")
			}
			return
		}

		sourceDTO := mapperv1.ToSourceDTO(source, svc.Config)

		if err := api.JSON(w, http.StatusOK, sourceDTO); err != nil {
			var writeErr *api.JSONWriteError
			if errors.As(err, &writeErr) {
				logrus.WithError(err).Error("Failed to write response")
				return
			}

			logrus.WithError(err).Error("Failed to marshal response")
			if err := api.JSON(
				w,
				http.StatusInternalServerError,
				map[string]string{"error": "Failed to get source"},
			); err != nil {
				logrus.WithError(err).Error("Failed to write error response")
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
			if err := api.JSON(
				w,
				http.StatusBadRequest,
				map[string]string{"error": "Invalid input parameters"},
			); err != nil {
				logrus.WithError(err).Error("Failed to write error response")
			}
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
			if err := api.JSON(
				w,
				http.StatusBadRequest,
				map[string]string{"error": "Invalid input parameters: " + err.Error()},
			); err != nil {
				logrus.WithError(err).Error("Failed to write error response")
			}
			return
		}

		source, err := svc.CreateSource(r.Context(), service.CreateSourceInput{
			EgressUrl:     input.EgressUrl,
			StaticHeaders: input.StaticHeaders,
			Description:   input.Description,
		})
		if err != nil {
			logrus.WithError(err).Error("Failed to create source")
			if err := api.JSON(
				w,
				http.StatusInternalServerError,
				map[string]string{"error": "Failed to create source"},
			); err != nil {
				logrus.WithError(err).Error("Failed to write error response")
			}
			return
		}
		logrus.WithField("source_id", source.ID).Info("Created new source")

		sourceDTO := mapperv1.ToSourceDTO(source, svc.Config)

		if err := api.JSON(w, http.StatusCreated, sourceDTO); err != nil {
			var writeErr *api.JSONWriteError
			if errors.As(err, &writeErr) {
				logrus.WithError(err).Error("Failed to write response")
				return
			}

			logrus.WithError(err).Error("Failed to marshal response")
			if err := api.JSON(
				w,
				http.StatusInternalServerError,
				map[string]string{"error": "Failed to create source"},
			); err != nil {
				logrus.WithError(err).Error("Failed to write error response")
			}
		}
	}
}

// updateSourceStatus handles PUT requests to update the status of a source.
func updateSourceStatus(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input, err := api.ParseRequestInput[requestsv1.UpdateSourceStatusInput](r)
		if err != nil {
			logrus.WithError(err).Error("Failed to parse input parameters")
			if err := api.JSON(
				w,
				http.StatusBadRequest,
				map[string]string{"error": "Invalid input parameters"},
			); err != nil {
				logrus.WithError(err).Error("Failed to write error response")
			}
			return
		}

		logrus.WithFields(logrus.Fields{
			"source_id":     input.SourceID,
			"status":        input.Status,
			"status_reason": input.StatusReason,
			"query":         r.URL.RawQuery,
		}).Debug("Received updateSourceStatus request")

		if err := requestsv1.ValidateUpdateSourceStatusInput(input); err != nil {
			logrus.WithError(err).Error("Input validation failed")
			if err := api.JSON(
				w,
				http.StatusBadRequest,
				map[string]string{"error": "Invalid input parameters: " + err.Error()},
			); err != nil {
				logrus.WithError(err).Error("Failed to write error response")
			}
			return
		}

		source, err := svc.UpdateSourceStatus(r.Context(), service.UpdateSourceStatusInput{
			SourceID:     input.SourceID,
			Status:       input.Status,
			StatusReason: input.StatusReason,
		})
		if err != nil {
			logrus.WithError(err).Error("Failed to create source")
			if err := api.JSON(
				w,
				http.StatusInternalServerError,
				map[string]string{"error": "Failed to create source"},
			); err != nil {
				logrus.WithError(err).Error("Failed to write error response")
			}
			return
		}
		logrus.WithFields(logrus.Fields{
			"source_id":     input.SourceID,
			"status":        input.Status,
			"status_reason": input.StatusReason,
		}).Info("Updated source status")

		sourceDTO := mapperv1.ToSourceDTO(source, svc.Config)

		if err := api.JSON(w, http.StatusOK, sourceDTO); err != nil {
			var writeErr *api.JSONWriteError
			if errors.As(err, &writeErr) {
				logrus.WithError(err).Error("Failed to write response")
				return
			}

			logrus.WithError(err).Error("Failed to marshal response")
			if err := api.JSON(
				w,
				http.StatusInternalServerError,
				map[string]string{"error": "Failed to update source status"},
			); err != nil {
				logrus.WithError(err).Error("Failed to write error response")
			}
		}
	}
}

// sourcesRouter sets up the router for sources-related endpoints.
func sourcesRouter(svc *service.Service) chi.Router {
	r := chi.NewRouter()
	r.Get("/", listSources(svc))
	r.Post("/", createSource(svc))
	r.Get("/{source_id}", getSourceByID(svc))
	r.Put("/{source_id}/status", updateSourceStatus(svc))
	r.Mount("/{source_id}/events", eventsRouter(svc))
	return r
}
