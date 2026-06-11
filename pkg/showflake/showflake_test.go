package showflake

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestShowflakeIDNext_OK(t *testing.T) {
	t.Parallel()
	fkCfg := SnowflakeConfig{
		Epoch:         time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		NodeID:        1,
		TimestampBits: 41,
		NodeBits:      10,
		SequenceBits:  12,
	}
	fk, err := NewSnowflake(fkCfg)
	require.NoError(t, err)

	timestampShift := fkCfg.NodeBits + fkCfg.SequenceBits
	want := 1 << timestampShift

	id, err := fk.NextID()
	require.NoError(t, err)
	require.GreaterOrEqual(t, id, int64(want))
}

func TestShowflakeIDNext_Concurency(t *testing.T) {
	t.Parallel()
	var wg sync.WaitGroup
	var m sync.Mutex

	gorutinesCount := 10000

	fkCfg := SnowflakeConfig{
		Epoch:         time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		NodeID:        1,
		TimestampBits: 41,
		NodeBits:      10,
		SequenceBits:  12,
	}
	fk, err := NewSnowflake(fkCfg)
	require.NoError(t, err)

	res := make(map[int64]struct{}, gorutinesCount)

	for range gorutinesCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id, err := fk.NextID()
			require.NoError(t, err)

			m.Lock()
			defer m.Unlock()

			_, ok := res[id]
			require.False(t, ok)
			res[id] = struct{}{}

		}()
	}
	wg.Wait()
}

func TestShowflakeIDNext_TimestampOverall(t *testing.T) {
	t.Parallel()

	fkCfg := SnowflakeConfig{
		Epoch:         time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
		NodeID:        1,
		TimestampBits: 41,
		NodeBits:      10,
		SequenceBits:  12,
	}
	_, err := NewSnowflake(fkCfg)
	require.Error(t, err)
}

func TestShowflakeIDNext_MinSlugLenght(t *testing.T) {
	t.Parallel()
	wantLenght := 10
	now := time.Now()
	fkCfg := SnowflakeConfig{
		Epoch:         now,
		NodeID:        0,
		TimestampBits: 41,
		NodeBits:      10,
		SequenceBits:  12,
	}
	fk, err := NewSnowflake(fkCfg)
	require.NoError(t, err)
	id, err := fk.NextIDInBase63(wantLenght)
	require.NoError(t, err)
	require.Equal(t, len([]rune(id)), wantLenght)

}

func TestShowflakeIDNext_MaxSlugLenght(t *testing.T) {
	t.Parallel()
	wantLenght := 8
	fkCfg := SnowflakeConfig{
		Epoch:         time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		NodeID:        0,
		TimestampBits: 41,
		NodeBits:      10,
		SequenceBits:  12,
	}
	fk, err := NewSnowflake(fkCfg)
	require.NoError(t, err)
	id, err := fk.NextIDInBase63(wantLenght)
	require.NoError(t, err)
	require.Equal(t, len([]rune(id)), wantLenght)

}

func TestShowflakeIDNext_SlugLenghtTabel(t *testing.T) {
	t.Parallel()
	oldDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	newDate := time.Now()
	tests := []struct {
		name       string
		cfg        SnowflakeConfig
		wantLenght int
	}{
		{
			name: "lenght 10 with big timestampt",
			cfg: SnowflakeConfig{
				Epoch:         newDate,
				NodeID:        0,
				TimestampBits: 41,
				NodeBits:      10,
				SequenceBits:  12,
			},
			wantLenght: 10,
		},
		{
			name: "lenght 10 with little timestampt",
			cfg: SnowflakeConfig{
				Epoch:         newDate,
				NodeID:        0,
				TimestampBits: 41,
				NodeBits:      10,
				SequenceBits:  12,
			},
			wantLenght: 10,
		},
		{
			name: "lenght 8 with big timestampt",
			cfg: SnowflakeConfig{
				Epoch:         oldDate,
				NodeID:        0,
				TimestampBits: 41,
				NodeBits:      10,
				SequenceBits:  12,
			},
			wantLenght: 8,
		},
		{
			name: "lenght 8 with little timestampt",
			cfg: SnowflakeConfig{
				Epoch:         newDate,
				NodeID:        0,
				TimestampBits: 41,
				NodeBits:      10,
				SequenceBits:  12,
			},
			wantLenght: 8,
		},
		{
			name: "lenght 12 with big timestampt",
			cfg: SnowflakeConfig{
				Epoch:         oldDate,
				NodeID:        0,
				TimestampBits: 41,
				NodeBits:      10,
				SequenceBits:  12,
			},
			wantLenght: 12,
		},
		{
			name: "lenght 12 with little timestampt",
			cfg: SnowflakeConfig{
				Epoch:         newDate,
				NodeID:        0,
				TimestampBits: 41,
				NodeBits:      10,
				SequenceBits:  12,
			},
			wantLenght: 12,
		},
		{
			name: "lenght 20 with little timestampt",
			cfg: SnowflakeConfig{
				Epoch:         newDate,
				NodeID:        0,
				TimestampBits: 41,
				NodeBits:      10,
				SequenceBits:  12,
			},
			wantLenght: 20,
		},
		{
			name: "lenght 20 with big timestampt",
			cfg: SnowflakeConfig{
				Epoch:         oldDate,
				NodeID:        0,
				TimestampBits: 41,
				NodeBits:      10,
				SequenceBits:  12,
			},
			wantLenght: 20,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()
				fk, err := NewSnowflake(tt.cfg)
				require.NoError(t, err)
				id, err := fk.NextIDInBase63(tt.wantLenght)
				require.NoError(t, err)
				require.Equal(t, len([]rune(id)), tt.wantLenght)
			},
		)
	}

}
