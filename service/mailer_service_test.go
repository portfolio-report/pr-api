package service

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMailerService(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		config         string
		recipientEmail string
		wantUser       string
		wantPass       string
		wantHost       string
		wantPort       string
		wantErr        bool
	}{
		{"", "", "", "", "", "", true},
		{"", "info@example.com", "", "", "", "", true},
		{"smtp://localhost", "", "", "", "", "", true},
		{"smtp://localhost", "info@example.com", "", "", "localhost", "25", false},
		{"smtp://localhost:587", "info@example.com", "", "", "localhost", "587", false},
		{"smtp://username:password@example.com:587", "info@example.com", "username", "password", "example.com", "587", false},
		{"smtp://username:password@example.com", "info@example.com", "username", "password", "example.com", "25", false},
	}
	for _, tc := range tests {
		t.Run(tc.config+" "+tc.recipientEmail, func(t *testing.T) {
			i, err := NewMailerService(tc.config, tc.recipientEmail, validate)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			require.Nil(t, err)
			require.NotNil(t, i)

			got := i.(*mailerService)
			assert.Equal(t, tc.wantUser, got.User)
			assert.Equal(t, tc.wantPass, got.Pass)
			assert.Equal(t, tc.wantHost, got.Host)
			assert.Equal(t, tc.wantPort, got.Port)
			assert.Equal(t, tc.recipientEmail, got.RecipientEmail)
		})
	}
}
