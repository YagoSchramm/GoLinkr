package util

import (
	"net/http"
	"net/url"
	"time"

	"github.com/YagoSchramm/Golinkr/domain/entity"
	"github.com/YagoSchramm/Golinkr/domain/usecase/dtos"
)

func BuildShortURL(r *http.Request, code string) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	if forwardedProto := r.Header.Get("X-Forwarded-Proto"); forwardedProto != "" {
		scheme = forwardedProto
	}

	return (&url.URL{
		Scheme: scheme,
		Host:   r.Host,
		Path:   "/link/" + code,
	}).String()
}

func BuildLinkDTO(r *http.Request, link entity.Link) dtos.LinkDTO {
	return dtos.LinkDTO{
		ID:          link.ID.String(),
		Code:        link.Code,
		OriginalURL: link.OriginalURL,
		CreatedAt:   link.CreatedAt.Format(time.RFC3339),
		ShortURL:    BuildShortURL(r, link.Code),
	}
}
