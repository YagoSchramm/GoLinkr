package rules

import (
	"github.com/YagoSchramm/Golinkr/domain/entity"
	"github.com/YagoSchramm/Golinkr/domain/entity/derr"
)

func ValidateAnalytics(analytics entity.Analytics) error {
	if analytics.LinkID == analytics.ID {
		return derr.InvalidLinkId
	}
	return nil
}
