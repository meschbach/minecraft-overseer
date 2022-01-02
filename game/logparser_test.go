package game

import (
	"testing"
)

func TestLoadingMessageIgnored(t *testing.T) {
	input := "[22:33:46] [Server thread/INFO]: Default game type: SURVIVAL"
	entry := parseLogEntry(input)
	if entry.AsString() != "Default game type: SURVIVAL" {
		t.Fatalf("Parser failed to extract random message, got '%s'", entry.AsString())
	}
}

func TestConsumesStartingMessage(t *testing.T) {
	input := "[22:33:46] [Server thread/INFO]: Starting minecraft server version 1.12"
	entry := parseLogEntry(input)
	startingEntry := entry.(*StartingEntry)
	if startingEntry.Version != "1.12" {
		t.Fatalf("Expected '1.12', got '%s'", startingEntry.Version)
	}
}

func TestConsumeStartedMessage(t *testing.T) {
	input := "[14:37:43] [Server thread/INFO]: Done (2.095s)! For help, type \"help\" or \"?\""
	rawEntry := parseLogEntry(input)
	entry := rawEntry.(*StartedEntry)
	if entry.TimeTaken != "2.095s" {
		t.Fatalf("Expected time taken to be '2.095s', got '%s'", entry.TimeTaken)
	}
}

func TestStoppingMessage(t *testing.T) {
	input := "[14:58:12] [Server thread/INFO]: Stopping the server"
	rawEntry := parseLogEntry(input)
	_, ok := rawEntry.(*StoppingEntry)
	if !ok {
		t.Fatalf("Expected a stopping message, got not okay")
	}
}

func TestUserJoined(t *testing.T) {
	input := "[15:03:18] [Server thread/INFO]: drakgremlin joined the game"
	rawEntry := parseLogEntry(input)
	entry := rawEntry.(*UserJoinedEntry)
	if entry.User != "drakgremlin" {
		t.Fatalf("Expected user to be 'drakgremlin', got '%s'", entry.User)
	}
}

func TestUserJoinedWithNewline(t *testing.T) {
	input := "[15:03:18] [Server thread/INFO]: drakgremlin joined the game\n"
	rawEntry := parseLogEntry(input)
	entry := rawEntry.(*UserJoinedEntry)
	if entry.User != "drakgremlin" {
		t.Fatalf("Expected user to be 'drakgremlin', got '%s'", entry.User)
	}
}

func TestUserLeft(t *testing.T) {
	input := "[15:03:38] [Server thread/INFO]: drakgremlin left the game"
	rawEntry := parseLogEntry(input)
	entry := rawEntry.(*UserLeftEvent)
	if entry.User != "drakgremlin" {
		t.Fatalf("Expected user to be 'drakgremlin', got '%s'", entry.User)
	}
}
