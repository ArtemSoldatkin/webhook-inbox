package mapperv1

import (
	dtov1 "github.com/ArtemSoldatkin/webhook-inbox/internal/api/dto/v1"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/utils"
)

// ToDeliveryAttemptDTO converts a db.DeliveryAttempt to a dtov1.DeliveryAttemptDTO.
func ToDeliveryAttemptDTO(
	deliveryAttempt db.DeliveryAttempt,
) dtov1.DeliveryAttemptDTO {
	return dtov1.DeliveryAttemptDTO{
		ID:            deliveryAttempt.ID,
		EventID:       deliveryAttempt.EventID,
		AttemptNumber: deliveryAttempt.AttemptNumber,
		State:         deliveryAttempt.State,
		StatusCode: utils.PtrIfValid(
			deliveryAttempt.StatusCode.Int32,
			deliveryAttempt.StatusCode.Valid,
		),
		ErrorType: utils.PtrIfValid(
			deliveryAttempt.ErrorType.String,
			deliveryAttempt.ErrorType.Valid,
		),
		ErrorMessage: utils.PtrIfValid(
			deliveryAttempt.ErrorMessage.String,
			deliveryAttempt.ErrorMessage.Valid,
		),
		StartedAt: utils.PtrIfValid(
			deliveryAttempt.StartedAt.Time,
			deliveryAttempt.StartedAt.Valid,
		),
		FinishedAt: utils.PtrIfValid(
			deliveryAttempt.FinishedAt.Time,
			deliveryAttempt.FinishedAt.Valid,
		),
		CreatedAt: deliveryAttempt.CreatedAt.Time,
		NextAttemptAt: utils.PtrIfValid(
			deliveryAttempt.NextAttemptAt.Time,
			deliveryAttempt.NextAttemptAt.Valid,
		),
	}
}

// ToDeliveryAttemptDTOs converts a slice of db.DeliveryAttempt to a slice of dtov1.DeliveryAttemptDTO.
func ToDeliveryAttemptDTOs(
	deliveryAttempts []db.DeliveryAttempt,
) []dtov1.DeliveryAttemptDTO {
	deliveryAttemptDTOs := make([]dtov1.DeliveryAttemptDTO, len(deliveryAttempts))
	for i, deliveryAttempt := range deliveryAttempts {
		deliveryAttemptDTOs[i] = ToDeliveryAttemptDTO(deliveryAttempt)
	}
	return deliveryAttemptDTOs
}
