package events

import (
	"testing"
)

func TestLoadingMessageIgnored(t *testing.T) {
	t.Parallel()

	input := "[22:33:46] [Server thread/INFO]: Default game type: SURVIVAL"
	entry := ParseLogEntry(input)
	if entry.String() != "Default game type: SURVIVAL" {
		t.Fatalf("Parser failed to extract random message, got '%s'", entry.String())
	}
}

func TestConsumesStartingMessage(t *testing.T) {
	t.Parallel()

	input := "[22:33:46] [Server thread/INFO]: Starting minecraft server version 1.12"
	entry := ParseLogEntry(input)
	startingEntry := entry.(*StartingEntry)
	if startingEntry.Version != "1.12" {
		t.Fatalf("Expected '1.12', got '%s'", startingEntry.Version)
	}
}

func TestConsumeStartedMessage(t *testing.T) {
	t.Parallel()

	input := "[14:37:43] [Server thread/INFO]: Done (2.095s)! For help, type \"help\" or \"?\""
	rawEntry := ParseLogEntry(input)
	entry := rawEntry.(*StartedEntry)
	if entry.TimeTaken != "2.095s" {
		t.Fatalf("Expected time taken to be '2.095s', got '%s'", entry.TimeTaken)
	}
}

func TestConsumeStartedMessage_17_1(t *testing.T) {
	t.Parallel()

	input := "[05:43:55] [Server thread/INFO]: Done (7.682s)! For help, type \"help\""
	rawEntry := ParseLogEntry(input)
	entry := rawEntry.(*StartedEntry)
	if entry.TimeTaken != "7.682s" {
		t.Fatalf("Expected time taken to be '7.682s', got '%s'", entry.TimeTaken)
	}
}

func TestStoppingMessage(t *testing.T) {
	t.Parallel()

	input := "[14:58:12] [Server thread/INFO]: Stopping the server"
	rawEntry := ParseLogEntry(input)
	_, ok := rawEntry.(*StoppingEntry)
	if !ok {
		t.Fatalf("Expected a stopping message, got not okay")
	}
}

func TestUserJoined(t *testing.T) {
	t.Parallel()

	input := "[15:03:18] [Server thread/INFO]: drakgremlin joined the game"
	rawEntry := ParseLogEntry(input)
	entry := rawEntry.(*UserJoinedEntry)
	if entry.User != "drakgremlin" {
		t.Fatalf("Expected user to be 'drakgremlin', got '%s'", entry.User)
	}
}

func TestUserJoinedWithNewline(t *testing.T) {
	t.Parallel()

	input := "[15:03:18] [Server thread/INFO]: drakgremlin joined the game\n"
	rawEntry := ParseLogEntry(input)
	entry := rawEntry.(*UserJoinedEntry)
	if entry.User != "drakgremlin" {
		t.Fatalf("Expected user to be 'drakgremlin', got '%s'", entry.User)
	}
}

func TestUserLeft(t *testing.T) {
	t.Parallel()

	input := "[15:03:38] [Server thread/INFO]: drakgremlin left the game"
	rawEntry := ParseLogEntry(input)
	entry := rawEntry.(*UserLeftEvent)
	if entry.User != "drakgremlin" {
		t.Fatalf("Expected user to be 'drakgremlin', got '%s'", entry.User)
	}
}
