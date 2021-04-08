package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestParseRawOutput(t *testing.T) {
	output := strings.Trim(`
ID     Done       Have  ETA           Up    Down  Ratio  Status       Name
  29    53%    3.42 GB  Unknown      0.0     0.0    0.0  Idle         test
  30    n/a    4.21 GB  Done         0.0     0.0   None  Stopped      test 2
Sum:           7.63 GB               0.0     0.0
`, "\n")

	torrentStates := parseRawOutput(output)
	want := []TorrentState{
		{"29", "test", "53%"},
		{"30", "test 2", "n/a"},
	}

	assert.Equal(t, want, torrentStates)
}

func TestParseRawOutputNoTorrent(t *testing.T) {
	// TODO: this is not a real output, I will update later
	output := strings.Trim(`
ID     Done       Have  ETA           Up    Down  Ratio  Status       Name
Sum:           7.63 GB               0.0     0.0
`, "\n")

	torrentStates := parseRawOutput(output)
	want := []TorrentState{}

	assert.Equal(t, want, torrentStates)
}

func TestFilterFinishedTorrents(t *testing.T) {
	input := []TorrentState{
		{"1", "Seed 1", "53%"},
		{"2", "Seed 2", "n/a"},
		{"3", "Seed 3", "100%"},
		{"4", "Seed 4", "100%"},
		{"5", "Seed 5", "100%"},
	}
	want := []TorrentState{
		{"3", "Seed 3", "100%"},
		{"4", "Seed 4", "100%"},
		{"5", "Seed 5", "100%"},
	}

	result := filterFinishedTorrents(input)

	assert.ElementsMatch(t, want, result)
}

func TestFilterFinishedTorrentsNoneFinished(t *testing.T) {
	input := []TorrentState{
		{"1", "Seed 1", "53%"},
		{"2", "Seed 2", "n/a"},
	}
	want := []TorrentState{}

	result := filterFinishedTorrents(input)

	assert.Equal(t, want, result)
}

func TestNotify(t *testing.T) {
	// TODO
}

func TestMain(t *testing.T) {
	// TODO
}
