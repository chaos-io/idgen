package redis_idgen

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/chaos-io/redis"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func newTestRedis(t *testing.T) (redis.Cmdable, func(), error) {
	m := miniredis.NewMiniRedis()
	if err := m.Start(); err != nil {
		return nil, nil, err
	}

	opts := &goredis.Options{Addr: m.Addr()}
	cli := goredis.NewClient(opts)

	provider := redis.NewProvider(cli)

	cleanup := func() {
		m.Close()
	}

	return provider, cleanup, nil
}

func Test_generator_GenMultiIDs(t *testing.T) {
	ctx := context.Background()

	cli, cleanup, err := newTestRedis(t)
	assert.Nil(t, err)
	t.Cleanup(cleanup)

	idgen, err := NewIDGenerator(cli, []int64{0, 1, 2})
	assert.Nil(t, err)

	ids, err := idgen.GenMultiIDs(ctx, 10)
	assert.Nil(t, err)
	assert.Equal(t, 10, len(ids))

	id, err := idgen.GenID(ctx)
	assert.Nil(t, err)
	assert.True(t, id >= time.Now().UnixNano()-(id>>32))
}
