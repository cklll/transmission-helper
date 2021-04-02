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
		{"idle seed", "Idle"},
		{"stopped 1", "Stopped"},
		{"stopped 2", "Stopped"},
		{"Seeding 1", "Seeding"},
		{"Seeding 2", "Seeding"},
		{"Finished 1", "Finished"},
	}
	want := []TorrentState{
		{"Seeding 1", "Seeding"},
		{"Seeding 2", "Seeding"},
		{"Finished 1", "Finished"},
	}

	result := filterFinishedTorrents(input)

	assert.Equal(t, want, result)
}

func TestFilterFinishedTorrentsNoneFinished(t *testing.T) {
	input := []TorrentState{
		{"idle seed", "Idle"},
		{"stopped 1", "Stopped"},
		{"stopped 2", "Stopped"},
	}
	want := []TorrentState{}

	result := filterFinishedTorrents(input)

	assert.Equal(t, want, result)
}

func TestFilterFinishedTorrentsLogAllStates(t *testing.T) {
	// TODO
}
