package events

import (
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserSaidEventParsing(t *testing.T) {
	input := "[23:14:48] [Server thread/INFO]: <drakgremlin> why no sleep"
	entry := ParseLogEntry(input)
	require.IsType(t, &UserSaidEvent{}, entry)
	sut := entry.(*UserSaidEvent)
	assert.Equal(t, "drakgremlin", sut.Speaker)
	assert.Equal(t, "why no sleep", sut.Message)
}

func TestUserSaidEventParsing_EmptyMessage(t *testing.T) {
	input := "[23:14:48] [Server thread/INFO]: <drakgremlin> "
	entry := ParseLogEntry(input)
	require.IsType(t, &UserSaidEvent{}, entry)
	sut := entry.(*UserSaidEvent)
	assert.Equal(t, "drakgremlin", sut.Speaker)
	assert.Equal(t, "", sut.Message)
}
