package generator

import (
	"testing"
	"time"

	"github.com/ArtemS18/ShortURL-API/internal/entity"
	"github.com/ArtemS18/ShortURL-API/pkg/showflake"
	"github.com/stretchr/testify/require"
)

func TestSlugGeneratorUseCase_GenerateSlug(t *testing.T) {
	t.Parallel()
	cfg := showflake.SnowflakeConfig{
		Epoch:         time.Now().Add(-1 * time.Hour),
		NodeID:        1,
		TimestampBits: 41,
		NodeBits:      8,
		SequenceBits:  10,
	}

	tests := []struct {
		name    string
		input   *entity.URL
		setup   func(t *testing.T) *showflake.Snowflake
		wantErr bool
	}{
		{
			name: "OK - Successfully generated slug",
			input: &entity.URL{
				Value: "https://ya.ru",
			},
			setup: func(t *testing.T) *showflake.Snowflake {
				sf, err := showflake.NewSnowflake(cfg)
				require.NoError(t, err)
				return sf
			},
			wantErr: false,
		},
		{
			name: "Error - Snowflake timestamp overflow",
			input: &entity.URL{
				Value: "https://google.com",
			},
			setup: func(t *testing.T) *showflake.Snowflake {
				brokenCfg := cfg
				brokenCfg.Epoch = time.Now().Add(100 * time.Hour)

				sf, err := showflake.NewSnowflake(brokenCfg)
				require.NoError(t, err)
				return sf
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			sf := test.setup(t)
			uc := NewSlugGeneratorUseCase(sf)

			res, err := uc.GenerateSlug(test.input)

			if test.wantErr {
				require.Error(t, err)
				require.Nil(t, res)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, res)

			require.Equal(t, test.input.Value, res.URL)

			require.Len(t, res.Slug, SlugLength)
			expectedSlug := sf.Int64ToBase63(res.ID, SlugLength)
			require.Equal(t, expectedSlug, res.Slug)
		})
	}
}
