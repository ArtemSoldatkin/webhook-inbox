// Package mapperv1 provides functions to map between DTO v1 and internal data structures.
package mapperv1

import (
	"time"

	dtov1 "github.com/ArtemSoldatkin/webhook-inbox/internal/api/dto/v1"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/config"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/utils"
	"github.com/sirupsen/logrus"
)

// ToSourceDTO converts a db.Source to a dtov1.SourceDTO.
func ToSourceDTO(
	source db.Source,
	config *config.Config,
) dtov1.SourceDTO {
	staticHeaders, err := utils.JSONBtoType[map[string]string](source.StaticHeaders)
	if err != nil {
		logrus.WithError(err).Error("Failed to unmarshal static headers")
		staticHeaders = map[string]string{
			"__error": "Webhook Inbox Error - Failed to parse",
		}
	}

	var disableAt *time.Time
	if source.DisableAt.Valid {
		disableAt = &source.DisableAt.Time
	} else {
		disableAt = nil
	}

	return dtov1.SourceDTO{
		ID:       source.ID,
		PublicID: source.PublicID.String(),
		IngressUrl: utils.GenerateIngressURL(
			config.APIProtocol,
			config.APIHost,
			config.APIPort,
			source.PublicID.String(),
		),
		EgressUrl:     source.EgressUrl,
		StaticHeaders: staticHeaders,
		Status:        source.Status,
		StatusReason: utils.PtrIfValid(
			source.StatusReason.String,
			source.StatusReason.Valid,
		),
		Description: utils.PtrIfValid(
			source.Description.String,
			source.Description.Valid,
		),
		CreatedAt: source.CreatedAt.Time,
		UpdatedAt: source.UpdatedAt.Time,
		DisableAt: disableAt,
	}
}

// ToSourceDTOs converts a slice of db.Source to a slice of dtov1.SourceDTO.
func ToSourceDTOs(
	sources []db.Source,
	config *config.Config,
) []dtov1.SourceDTO {
	sourceDTOs := make([]dtov1.SourceDTO, 0, len(sources))
	for _, source := range sources {
		sourceDTOs = append(sourceDTOs, ToSourceDTO(source, config))
	}
	return sourceDTOs
}
