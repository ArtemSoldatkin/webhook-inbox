package requestsv1

import (
	"errors"
	"regexp"
)

// This file defines the input parameters and validation logic for the ingestEvent API endpoint.
var uuidRegexp = regexp.MustCompile(
	`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`,
)

// IngestEventInput defines the expected input parameters for ingesting a new event.
type IngestEventInput struct {
	PublicID string `url_param:"public_id,required"`
}

// validatePublicID checks if the provided public ID matches the expected UUID format.
func validatePublicID(publicID string) bool {
	return uuidRegexp.MatchString(publicID)
}

// ValidateIngestEventInput validates the input parameters for ingesting a new event.
func ValidateIngestEventInput(input *IngestEventInput) error {
	if !validatePublicID(input.PublicID) {
		return errors.New("invalid public_id format")
	}
	return nil
}
