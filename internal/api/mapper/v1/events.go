package mapperv1

import (
	dtov1 "github.com/ArtemSoldatkin/webhook-inbox/internal/api/dto/v1"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/utils"
	"github.com/sirupsen/logrus"
)

// ToEventDTO converts a db.Event to a dtov1.EventDTO, handling JSONB fields and potential errors gracefully.
func ToEventDTO(event db.Event) dtov1.EventDTO {
	queryParams, err := utils.JSONBtoType[map[string][]string](event.QueryParams)
	if err != nil {
		logrus.WithError(err).Error("Failed to unmarshal query params")
		queryParams = map[string][]string{
			"__error": {"Webhook Inbox Error - Failed to parse"},
		}
	}

	rawHeaders, err := utils.JSONBtoType[map[string][]string](event.RawHeaders)
	if err != nil {
		logrus.WithError(err).Error("Failed to unmarshal raw headers")
		rawHeaders = map[string][]string{
			"__error": {"Webhook Inbox Error - Failed to parse"},
		}
	}

	var remoteAddress *string
	if event.RemoteAddress != nil {
		str := event.RemoteAddress.String()
		remoteAddress = &str
	}

	return dtov1.EventDTO{
		ID:              event.ID,
		SourceID:        event.SourceID,
		DedupHash:       event.DedupHash.String,
		Method:          event.Method,
		IngressPath:     event.IngressPath,
		RemoteAddress:   remoteAddress,
		QueryParams:     queryParams,
		RawHeaders:      rawHeaders,
		Body:            event.Body,
		BodyContentType: event.BodyContentType,
		ReceivedAt:      event.ReceivedAt.Time,
	}
}

// ToEventDTOs converts a slice of db.Event to a slice of dtov1.EventDTO using the ToEventDTO function for each event.
func ToEventDTOs(events []db.Event) []dtov1.EventDTO {
	eventDTOs := make([]dtov1.EventDTO, 0, len(events))
	for _, event := range events {
		eventDTOs = append(eventDTOs, ToEventDTO(event))
	}
	return eventDTOs
}
