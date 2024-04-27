package pgs

import (
	"errors"
	"fmt"
	"github.com/go-faster/city"
	"github.com/spaolacci/murmur3"
)

var (
	ErrShardIndexOutOfRange = errors.New("shard index is out of range")
)

type ShardKey string
type ShardIndex int

type ShardFn func(ShardKey) ShardIndex

type Manager struct {
	fn     ShardFn
	shards []*DB
}

func GetMurmur3ShardFn(shardsCnt int) ShardFn {
	hasher := murmur3.New32()
	return func(key ShardKey) ShardIndex {
		defer hasher.Reset()
		_, _ = hasher.Write([]byte(key))
		return ShardIndex(hasher.Sum32() % uint32(shardsCnt))
	}
}

func GetCityHashShardFn(shardsCnt int) ShardFn {
	return func(key ShardKey) ShardIndex {
		hash := city.Hash64([]byte(key))
		return ShardIndex(hash % uint64(shardsCnt))
	}
}

func NewShardManager(fn ShardFn, shards []*DB) *Manager {
	return &Manager{
		fn:     fn,
		shards: shards,
	}
}

func (m *Manager) GetShardIndex(key ShardKey) ShardIndex {
	return m.fn(key)
}

func (m *Manager) GetShardIndexFromID(id int64) ShardIndex {
	return ShardIndex(id % 2)
}

func (m *Manager) Pick(index ShardIndex) (*DB, error) {
	if int(index) < len(m.shards) {
		return m.shards[index], nil
	}

	return nil, fmt.Errorf("%w: given index=%d, len=%d", ErrShardIndexOutOfRange, index, len(m.shards))
}

func (m *Manager) AutoPickIndex(id int64) ShardIndex {
	return m.GetShardIndexFromID(id)
}

func (m *Manager) GerAllShards() []*DB {
	return m.shards
}

func (m *Manager) Close() {
	for _, shard := range m.shards {
		shard.Close()
	}
}
