package rules

import (
	"net/url"
	"strings"

	"github.com/YagoSchramm/Golinkr/domain/entity/derr"
)

func ValidateURL(rawURL string) error {
	if strings.TrimSpace(rawURL) == "" {
		return derr.InvalidURLError
	}

	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}

	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return derr.InvalidURLError
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return derr.InvalidURLError
	}

	if !strings.Contains(parsedURL.Host, ".") {
		return derr.InvalidURLError
	}

	return nil
}
