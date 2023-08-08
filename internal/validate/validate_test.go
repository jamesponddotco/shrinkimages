package validate_test

import (
	"testing"

	"git.sr.ht/~jamesponddotco/shrinkimages/internal/validate"
)

func TestAPIKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		key  string
		want bool
	}{
		{
			name: "Valid Key",
			key:  "shrink-EfYoM8FwPF7uycxzehm9i3qu2oS1pRd1",
			want: true,
		},
		{
			name: "Invalid Key: Too Short",
			key:  "shrink-ABCD",
			want: false,
		},
		{
			name: "Invalid Key: No Prefix",
			key:  "ABCDEFGHIJKLMNOPQRSTUVWXYZabcd",
			want: false,
		},
		{
			name: "Invalid Key: Not base64",
			key:  "shrink-ABCD**",
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := validate.APIKey(tt.key); got != tt.want {
				t.Errorf("APIKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
