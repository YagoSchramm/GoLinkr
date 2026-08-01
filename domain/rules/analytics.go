package rules

import (
	"github.com/YagoSchramm/Golinkr/domain/entity"
	"github.com/YagoSchramm/Golinkr/domain/entity/derr"
	"github.com/google/uuid"
)

func ValidateAnalytics(analytics entity.Analytics) error {
	if analytics.LinkID == uuid.Nil {
		return derr.InvalidLinkId
	}
	return nil
}
