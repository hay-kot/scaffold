package scaffold

import (
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestMetadata_IsCompatible(t *testing.T) {
	tests := []struct {
		name           string
		minimumVersion string
		currentVersion string
		expected       bool
		expectError    bool
	}{
		{
			name:           "Compatible with wildcard",
			minimumVersion: "*",
			currentVersion: "1.0.0",
			expected:       true,
			expectError:    false,
		},
		{
			name:           "Compatible with empty minimum version",
			minimumVersion: "",
			currentVersion: "1.0.0",
			expected:       true,
			expectError:    false,
		},
		{
			name:           "Compatible with dev version",
			minimumVersion: "1.0.0",
			currentVersion: "dev",
			expected:       true,
			expectError:    false,
		},
		{
			name:           "Compatible version",
			minimumVersion: "1.0.0",
			currentVersion: "1.0.0",
			expected:       true,
			expectError:    false,
		},
		{
			name:           "Incompatible version",
			minimumVersion: "2.0.0",
			currentVersion: "1.0.0",
			expected:       false,
			expectError:    false,
		},
		{
			name:           "Error parsing current version",
			minimumVersion: "1.0.0",
			currentVersion: "invalid",
			expected:       false,
			expectError:    true,
		},
		{
			name:           "Error parsing minimum version",
			minimumVersion: "invalid",
			currentVersion: "1.0.0",
			expected:       false,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Metadata{MinimumVersion: tt.minimumVersion}
			logger := log.With().Logger()
			result, err := m.IsCompatible(logger, tt.currentVersion)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}
