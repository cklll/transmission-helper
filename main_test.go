package main

import (
    "testing"
    "strings"
		"github.com/stretchr/testify/assert"
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
		{"test", "Idle"},
		{"test 2", "Stopped"},
	}

	assert.Equal(t, torrentStates, want)
}

func TestParseRawOutputNoTorrent(t *testing.T) {
// TODO: this is not a real output, I will update later
	output := strings.Trim(`
ID     Done       Have  ETA           Up    Down  Ratio  Status       Name
Sum:           7.63 GB               0.0     0.0
`, "\n")

	torrentStates := parseRawOutput(output)
	want := []TorrentState{}

	assert.Equal(t, torrentStates, want)
}
