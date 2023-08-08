package cryptoutil_test

import (
	"bytes"
	"strings"
	"testing"

	"git.sr.ht/~jamesponddotco/shrinkimages/internal/cryptoutil"
)

func TestGenerateAPIKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		random  []byte
		wantErr bool
	}{
		{
			name:    "normal case",
			random:  bytes.Repeat([]byte{1}, 24),
			wantErr: false,
		},
		{
			name:    "error case",
			random:  []byte{}, // insufficient bytes
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			reader := bytes.NewBuffer(tt.random)

			got, err := cryptoutil.GenerateAPIKey(reader)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GenerateAPIKey() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if err == nil && !strings.HasPrefix(got, cryptoutil.KeyPrefix) {
				t.Fatalf("GenerateAPIKey() got = %v, want prefix %v", got, cryptoutil.KeyPrefix)
			}
		})
	}
}
