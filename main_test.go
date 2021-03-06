package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

// https://stackoverflow.com/a/45803548/2691976
// mock exec.Command
var mockedExitStatus = 0
var mockedStdout string
var gotExecCommandCommand string
var gotExecCommandArgs []string

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	gotExecCommandCommand = command
	gotExecCommandArgs = args

	cs := []string{"-test.run=TestExecCommandHelper", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	es := strconv.Itoa(mockedExitStatus)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1",
		"STDOUT=" + mockedStdout,
		"EXIT_STATUS=" + es}
	return cmd
}
func TestExecCommandHelper(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	fmt.Fprintf(os.Stdout, os.Getenv("STDOUT"))
	i, _ := strconv.Atoi(os.Getenv("EXIT_STATUS"))
	os.Exit(i)
}

func TestParseRawOutput(t *testing.T) {
	assertTorrentStates := func(t testing.TB, want, got []TorrentState) {
		assert.Equal(t, want, got)
	}

	t.Run("when multiple torrents", func(t *testing.T) {
		output := strings.TrimPrefix(`
ID     Done       Have  ETA           Up    Down  Ratio  Status       Name
  29    53%    3.42 GB  Unknown      0.0     0.0    0.0  Idle         test
  30    n/a    4.21 GB  Done         0.0     0.0   None  Stopped      test 2
Sum:           7.63 GB               0.0     0.0
`, "\n")

		got := parseRawOutput(output)
		want := []TorrentState{
			{"29", "test", "53%"},
			{"30", "test 2", "n/a"},
		}
		assertTorrentStates(t, want, got)
	})

	t.Run("when no torrent", func(t *testing.T) {
		output := strings.TrimPrefix(`
ID     Done       Have  ETA           Up    Down  Ratio  Status       Name
Sum:              None               0.0     0.0
`, "\n")

		got := parseRawOutput(output)
		want := []TorrentState{}
		assertTorrentStates(t, want, got)
	})
}

func TestFilterFinishedTorrents(t *testing.T) {
	assertTorrentStates := func(t testing.TB, want, got []TorrentState) {
		assert.Equal(t, want, got)
	}

	t.Run("when have finished torrents", func(t *testing.T) {
		input := []TorrentState{
			{"1", "Seed 1", "53%"},
			{"2", "Seed 2", "n/a"},
			{"3", "Seed 3", "100%"},
			{"4", "Seed 4", "100%"},
			{"5", "Seed 5", "100%"},
		}
		got := filterFinishedTorrents(input)
		want := []TorrentState{
			{"3", "Seed 3", "100%"},
			{"4", "Seed 4", "100%"},
			{"5", "Seed 5", "100%"},
		}

		assertTorrentStates(t, want, got)
	})

	t.Run("when don't have finished torrents", func(t *testing.T) {
		input := []TorrentState{
			{"1", "Seed 1", "53%"},
			{"2", "Seed 2", "n/a"},
		}
		got := filterFinishedTorrents(input)
		want := []TorrentState{}

		assertTorrentStates(t, want, got)
	})
}

func TestGetTranmissionRemoteListOutput(t *testing.T) {
	oldExecCommand := execCommand
	defer func() { execCommand = oldExecCommand }()

	appConfig := getApplicationConfig("./testdata/config/example.yaml")
	mockedStdout = "mocked stdout"
	execCommand = fakeExecCommand
	gotStdout := getTranmissionRemoteListOutput(appConfig)

	assert.Equal(t, "transmission-remote", gotExecCommandCommand)
	assert.Equal(t, []string{"--auth", "test_tr_user:test_tr_password", "-l"}, gotExecCommandArgs)
	assert.Equal(t, "mocked stdout", gotStdout)
}

func TestNotify(t *testing.T) {
	// TODO
}

func TestDelete(t *testing.T) {
	oldExecCommand := fakeExecCommand
	defer func() { execCommand = oldExecCommand }()

	appConfig := getApplicationConfig("./testdata/config/example.yaml")
	states := []TorrentState{
		{"1", "Seed 1", "100%"},
		{"2", "Seed 2", "100%"},
	}

	execCommand = fakeExecCommand
	delete(appConfig, states)

	// TODO: not very nice. We only asserted the 2nd cmd call
	assert.Equal(t, "transmission-remote", gotExecCommandCommand)
	assert.Equal(t, []string{"--auth", "test_tr_user:test_tr_password", "-t2", "-r"}, gotExecCommandArgs)
}
