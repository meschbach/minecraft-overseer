package events

import (
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserKilledByZombie(t *testing.T) {
	input := "[23:27:27] [Server thread/INFO]: drakgremlin was slain by Zombie"
	entry := ParseLogEntry(input)
	require.IsType(t, &GenericDeathMessage{}, entry)
	sut := entry.(*GenericDeathMessage)
	assert.Equal(t, "drakgremlin was slain by Zombie", sut.Message)
}

func TestUserKilledByWitch(t *testing.T) {
	input := "[23:28:14] [Server thread/INFO]: drakgremlin was killed by Witch using magic"
	entry := ParseLogEntry(input)
	require.IsType(t, &GenericDeathMessage{}, entry)
	sut := entry.(*GenericDeathMessage)
	assert.Equal(t, "drakgremlin was killed by Witch using magic", sut.Message)
}
