package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"os/exec"
	"os"
	"strconv"
	"fmt"
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

func TestGetTranmissionRemoteListOutput(t *testing.T) {
	oldExecCommand := execCommand
	defer func() { execCommand = oldExecCommand }()

	os.Clearenv()
	os.Setenv("TH_REMOTE_USERNAME", "test_user")
	os.Setenv("TH_REMOTE_PASSWORD", "test_password")

	mockedStdout = "mocked stdout"
	execCommand = fakeExecCommand
	gotStdout := getTranmissionRemoteListOutput()

	assert.Equal(t, "transmission-remote", gotExecCommandCommand)
	assert.Equal(t, []string{"--auth", "test_user:test_password", "-l"}, gotExecCommandArgs)
	assert.Equal(t, "mocked stdout", gotStdout)
}

func TestNotify(t *testing.T) {
	// TODO
}

func TestMain(t *testing.T) {
	// TODO
}
