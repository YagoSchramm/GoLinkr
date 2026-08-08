package rules_test

import (
	"errors"
	"testing"

	"github.com/YagoSchramm/Golinkr/domain/entity"
	"github.com/YagoSchramm/Golinkr/domain/entity/derr"
	"github.com/YagoSchramm/Golinkr/domain/rules"
	"github.com/google/uuid"
)

func TestValidateRegister(t *testing.T) {
	tests := []struct {
		name    string
		user    entity.User
		wantErr error
	}{
		{
			name: "valid user",
			user: entity.User{
				Name:     "Yago",
				Email:    "yago@example.com",
				Password: "secret123",
			},
		},
		{
			name: "blank name",
			user: entity.User{
				Name:     " ",
				Email:    "yago@example.com",
				Password: "secret123",
			},
			wantErr: derr.NameIsTooShort,
		},
		{
			name: "blank email",
			user: entity.User{
				Name:     "Yago",
				Email:    " ",
				Password: "secret123",
			},
			wantErr: derr.EmailRequired,
		},
		{
			name: "email without at sign",
			user: entity.User{
				Name:     "Yago",
				Email:    "yago.example.com",
				Password: "secret123",
			},
			wantErr: derr.InvalidEmail,
		},
		{
			name: "email without local part",
			user: entity.User{
				Name:     "Yago",
				Email:    "@example.com",
				Password: "secret123",
			},
			wantErr: derr.InvalidEmail,
		},
		{
			name: "email without domain dot",
			user: entity.User{
				Name:     "Yago",
				Email:    "yago@example",
				Password: "secret123",
			},
			wantErr: derr.InvalidEmail,
		},
		{
			name: "email with domain starting dot",
			user: entity.User{
				Name:     "Yago",
				Email:    "yago@.example.com",
				Password: "secret123",
			},
			wantErr: derr.InvalidEmail,
		},
		{
			name: "email with domain ending dot",
			user: entity.User{
				Name:     "Yago",
				Email:    "yago@example.",
				Password: "secret123",
			},
			wantErr: derr.InvalidEmail,
		},
		{
			name: "blank password",
			user: entity.User{
				Name:     "Yago",
				Email:    "yago@example.com",
				Password: " ",
			},
			wantErr: derr.PasswordRequired,
		},
		{
			name: "short password",
			user: entity.User{
				Name:     "Yago",
				Email:    "yago@example.com",
				Password: "sec123",
			},
			wantErr: derr.WeakPassword,
		},
		{
			name: "password without digit",
			user: entity.User{
				Name:     "Yago",
				Email:    "yago@example.com",
				Password: "secretabc",
			},
			wantErr: derr.WeakPassword,
		},
		{
			name: "password without letter",
			user: entity.User{
				Name:     "Yago",
				Email:    "yago@example.com",
				Password: "12345678",
			},
			wantErr: derr.WeakPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rules.ValidateRegister(tt.user)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestValidateLogin(t *testing.T) {
	tests := []struct {
		name        string
		credentials entity.UserCredentials
		wantErr     error
	}{
		{
			name: "valid credentials",
			credentials: entity.UserCredentials{
				Email:    "yago@example.com",
				Password: "secret123",
			},
		},
		{
			name: "blank email",
			credentials: entity.UserCredentials{
				Email:    " ",
				Password: "secret123",
			},
			wantErr: derr.EmailRequired,
		},
		{
			name: "invalid email",
			credentials: entity.UserCredentials{
				Email:    "yago@example",
				Password: "secret123",
			},
			wantErr: derr.InvalidEmail,
		},
		{
			name: "blank password",
			credentials: entity.UserCredentials{
				Email:    "yago@example.com",
				Password: " ",
			},
			wantErr: derr.PasswordRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rules.ValidateLogin(tt.credentials)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		rawURL  string
		wantErr error
	}{
		{
			name:   "valid https url",
			rawURL: "https://example.com/docs",
		},
		{
			name:   "valid http url",
			rawURL: "http://example.com/docs",
		},
		{
			name:   "valid url without scheme",
			rawURL: "example.com/docs",
		},
		{
			name:    "blank url",
			rawURL:  " ",
			wantErr: derr.InvalidURLError,
		},
		{
			name:    "invalid syntax",
			rawURL:  "https://exa mple.com",
			wantErr: derr.InvalidURLError,
		},
		{
			name:    "host without dot",
			rawURL:  "https://localhost",
			wantErr: derr.InvalidURLError,
		},
		{
			name:    "missing host",
			rawURL:  "https:///docs",
			wantErr: derr.InvalidURLError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rules.ValidateURL(tt.rawURL)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestValidateAnalytics(t *testing.T) {
	linkID := uuid.New()

	tests := []struct {
		name      string
		analytics entity.Analytics
		wantErr   error
	}{
		{
			name: "valid analytics",
			analytics: entity.Analytics{
				LinkID: linkID,
			},
		},
		{
			name:      "nil link id",
			analytics: entity.Analytics{},
			wantErr:   derr.InvalidLinkId,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rules.ValidateAnalytics(tt.analytics)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestValidateClickAverageRules(t *testing.T) {
	linkID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name    string
		linkID  uuid.UUID
		userID  uuid.UUID
		rule    func(uuid.UUID, uuid.UUID) error
		wantErr error
	}{
		{
			name:   "hourly click average valid ids",
			linkID: linkID,
			userID: userID,
			rule:   rules.ValidateHourlyClickAverage,
		},
		{
			name:    "hourly click average nil link id",
			userID:  userID,
			rule:    rules.ValidateHourlyClickAverage,
			wantErr: derr.InvalidLinkId,
		},
		{
			name:    "hourly click average nil user id",
			linkID:  linkID,
			rule:    rules.ValidateHourlyClickAverage,
			wantErr: derr.UnauthorizedError,
		},
		{
			name:   "weekday click average valid ids",
			linkID: linkID,
			userID: userID,
			rule:   rules.ValidateWeekdayClickAverage,
		},
		{
			name:    "weekday click average nil link id",
			userID:  userID,
			rule:    rules.ValidateWeekdayClickAverage,
			wantErr: derr.InvalidLinkId,
		},
		{
			name:    "weekday click average nil user id",
			linkID:  linkID,
			rule:    rules.ValidateWeekdayClickAverage,
			wantErr: derr.UnauthorizedError,
		},
		{
			name:   "monthly week click average valid ids",
			linkID: linkID,
			userID: userID,
			rule:   rules.ValidateMonthlyWeekClickAverage,
		},
		{
			name:    "monthly week click average nil link id",
			userID:  userID,
			rule:    rules.ValidateMonthlyWeekClickAverage,
			wantErr: derr.InvalidLinkId,
		},
		{
			name:    "monthly week click average nil user id",
			linkID:  linkID,
			rule:    rules.ValidateMonthlyWeekClickAverage,
			wantErr: derr.UnauthorizedError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rule(tt.linkID, tt.userID)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}
