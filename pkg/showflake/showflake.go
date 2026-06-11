package showflake

import (
	"fmt"
	"sync"
	"time"
)

const base63Chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-"

type SnowflakeConfig struct {
	Epoch         time.Time
	NodeID        int64
	TimestampBits int
	NodeBits      int
	SequenceBits  int
}

type Snowflake struct {
	config         SnowflakeConfig
	sequence       int64
	lastTimestamp  int64
	mu             sync.Mutex
	maxSequence    int64
	maxNodeID      int64
	timestampShift int
	nodeShift      int
}

func NewSnowflake(cfg SnowflakeConfig) (*Snowflake, error) {
	if cfg.NodeID < 0 || cfg.NodeID > (1<<cfg.NodeBits)-1 {
		return nil, fmt.Errorf("NodeID должен быть в диапазоне [0, %d]", (1<<cfg.NodeBits)-1)
	}
	if cfg.TimestampBits+cfg.NodeBits+cfg.SequenceBits > 63 {
		return nil, fmt.Errorf("сумма бит не должна превышать 63")
	}

	nowMs := time.Now().UnixMilli() - cfg.Epoch.UnixMilli()
	if nowMs > (1<<cfg.TimestampBits)-1 {
		return nil, fmt.Errorf("timestamp переполнен")
	}

	sf := &Snowflake{
		config:         cfg,
		sequence:       0,
		lastTimestamp:  0,
		maxSequence:    (1 << cfg.SequenceBits) - 1,
		maxNodeID:      (1 << cfg.NodeBits) - 1,
		timestampShift: cfg.NodeBits + cfg.SequenceBits,
		nodeShift:      cfg.SequenceBits,
	}
	return sf, nil
}

func (sf *Snowflake) NextID() (int64, error) {
	sf.mu.Lock()
	defer sf.mu.Unlock()

	nowMs := time.Now().UnixMilli() - sf.config.Epoch.UnixMilli()

	if nowMs < sf.lastTimestamp {
		return 0, fmt.Errorf("clock moved backwards or invalid epoch: current %d, last %d", nowMs, sf.lastTimestamp)
	}
	if nowMs == sf.lastTimestamp {
		if sf.sequence >= sf.maxSequence {
			for nowMs <= sf.lastTimestamp {
				nowMs = time.Now().UnixMilli() - sf.config.Epoch.UnixMilli()
			}
			sf.sequence = 0
		} else {
			sf.sequence++
		}
	} else {
		sf.sequence = 0
	}

	if nowMs > (1<<sf.config.TimestampBits)-1 {
		return 0, fmt.Errorf("timestamp overflow")
	}
	sf.lastTimestamp = nowMs
	id := (nowMs << sf.timestampShift) |
		(sf.config.NodeID << sf.nodeShift) |
		(sf.sequence)

	return id, nil
}

func (sf *Snowflake) Int64ToBase63(id int64, length int) string {
	result := make([]byte, length)
	for i := range result {
		result[i] = base63Chars[0]
	}

	idx := length - 1
	for id > 0 && idx >= 0 {
		result[idx] = base63Chars[id%63]
		id /= 63
		idx--
	}

	return string(result)
}

func (sf *Snowflake) NextIDInBase63(length int) (string, error) {
	id, err := sf.NextID()
	if err != nil {
		return "", fmt.Errorf("sf.NextID: %v", err)
	}
	return sf.Int64ToBase63(id, length), nil
}
