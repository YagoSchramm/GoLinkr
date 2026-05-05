package rules

import (
	"net/url"

	"github.com/YagoSchramm/Golinkr/domain/entity/derr"
)

func ValidateURL(rawURL string) error {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return derr.InvalidURL
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return derr.InvalidURL
	}

	return nil
}
