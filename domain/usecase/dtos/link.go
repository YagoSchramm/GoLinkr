package dtos

type CreateLinkDTO struct {
	OriginalURL string `json:"original_url"`
}

type LinkDTO struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	OriginalURL string `json:"original_url"`
	CreatedAt   string `json:"created_at"`
	ShortURL    string `json:"short_url,omitempty"`
}

type ErrorDTO struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
