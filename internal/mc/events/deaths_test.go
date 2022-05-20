package events

import (
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserKilledByZombie(t *testing.T) {
	t.Parallel()

	input := "[23:27:27] [Server thread/INFO]: drakgremlin was slain by Zombie"
	entry := ParseLogEntry(input)
	require.IsType(t, &GenericDeathMessage{}, entry)
	sut := entry.(*GenericDeathMessage)
	assert.Equal(t, sut.Message, "drakgremlin was slain by Zombie")
}

func TestUserKilledByWitch(t *testing.T) {
	t.Parallel()

	input := "[23:28:14] [Server thread/INFO]: drakgremlin was killed by Witch using magic"
	entry := ParseLogEntry(input)
	require.IsType(t, &GenericDeathMessage{}, entry)
	sut := entry.(*GenericDeathMessage)
	assert.Equal(t, sut.Message, "drakgremlin was killed by Witch using magic")
}

func TestUserKilledBySpider(t *testing.T) {
	t.Parallel()

	input := "[23:42:53] [Server thread/INFO]: drakgremlin was slain by Spider"
	entry := ParseLogEntry(input)
	require.IsType(t, &GenericDeathMessage{}, entry)
	sut := entry.(*GenericDeathMessage)
	assert.Equal(t, sut.Message, "drakgremlin was slain by Spider")
}

func TestUserKilledByEnderman(t *testing.T) {
	t.Parallel()

	input := "[23:40:24] [Server thread/INFO]: drakgremlin was slain by Enderman"
	entry := ParseLogEntry(input)
	require.IsType(t, &GenericDeathMessage{}, entry)
	sut := entry.(*GenericDeathMessage)
	assert.Equal(t, sut.Message, "drakgremlin was slain by Enderman")
}

func TestUserShotBySkelaton(t *testing.T) {
	t.Parallel()

	input := "[23:40:52] [Server thread/INFO]: drakgremlin was shot by Skeleton"
	entry := ParseLogEntry(input)
	require.IsType(t, &GenericDeathMessage{}, entry)
	sut := entry.(*GenericDeathMessage)
	assert.Equal(t, sut.Message, "drakgremlin was shot by Skeleton")
}
