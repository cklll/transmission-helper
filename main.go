package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"os/exec"
	"strings"
)

type TorrentState struct {
	Id   string
	Name string
	Done string
}

var execCommand = exec.Command

// =========
// ID     Done       Have  ETA           Up    Down  Ratio  Status       Name
//   29    53%    3.42 GB  Unknown      0.0     0.0    0.0  Idle         test
//   30    n/a    4.21 GB  Done         0.0     0.0   None  Stopped      test 2
// Sum:           7.63 GB               0.0     0.0
// =========
// ID & Done are the 1st and 2nd parts of the line
// For seed name, we note that the header "Name" and the actual seed name start at the same character
// so we collect all characters after that position
func parseRawOutput(output string) []TorrentState {
	states := []TorrentState{}

	lines := strings.Split(output, "\n")

	headerLine := lines[0]
	namePosition := strings.Index(headerLine, "Name")

	// 2nd last line is "Sum: xxxx"
	// last line is empty string
	lines = lines[1 : len(lines)-2]
	for _, line := range lines {
		parts := strings.Fields(line)

		id := parts[0]
		done := parts[1]
		name := line[namePosition:]

		states = append(states, TorrentState{Id: id, Name: name, Done: done})
	}

	return states
}

func getTranmissionRemoteListOutput() string {
	username := os.Getenv("TH_REMOTE_USERNAME")
	password := os.Getenv("TH_REMOTE_PASSWORD")
	auth := fmt.Sprintf("%v:%v", username, password)

	// this assume we always need auth. It may not always be the case
	cmd := execCommand("transmission-remote", "--auth", auth, "-l")
	stdout, err := cmd.Output()

	if err != nil {
		log.Println(err.Error())
		panic("Cannot get output from transmission-remote. Please check the environment variables for TH_REMOTE_USERNAME and TH_REMOTE_PASSWORD.")
	}

	return string(stdout)
}

func filterFinishedTorrents(states []TorrentState) []TorrentState {
	result := []TorrentState{}

	for _, state := range states {
		if state.Done == "100%" {
			result = append(result, state)
		}
	}

	log.Printf("%v of %v torrents are finished.\n", len(result), len(states))

	return result
}

// TODO: TEST CASE!!!!
// Need refactoring to make this testable
func notify(finishedTorrentStates []TorrentState) {
	log.Printf("Sending emails with %v finished torrents.\n", len(finishedTorrentStates))

	recipientsString := os.Getenv("TH_NOTIFY_EMAILS")
	recipients := strings.Split(recipientsString, ",")

	subject := fmt.Sprintf("[transmission-helper] %v torrents completed.", len(finishedTorrentStates))
	message := ""
	for _, state := range finishedTorrentStates {
		message += fmt.Sprintf("%v: %v \r\n", state.Done, state.Name)
	}

	mailNotifier := MailNotifier{smtp.SendMail}
	mailConfig := mailNotifier.GetMailConfig()
	err := mailNotifier.Send(mailConfig, subject, message, recipients)
	if err != nil {
		log.Printf("Failed to send emails. Error: %v\n", err.Error())
	} else {
		log.Printf("Sent emails to %v recipients.\n", len(recipients))
	}
}

func delete(states []TorrentState) {
	username := os.Getenv("TH_REMOTE_USERNAME")
	password := os.Getenv("TH_REMOTE_PASSWORD")
	auth := fmt.Sprintf("%v:%v", username, password)

	for _, state := range states {
		deleteArg := fmt.Sprintf("-t%v", state.Id)
		// transmission-remote --auth <user>:<password> -t<ID> -r
		cmd := execCommand("transmission-remote", "--auth", auth, deleteArg, "-r")
		_, err := cmd.Output()

		if err != nil {
			log.Printf("Cannot delete torrent ID %v. Error: %v\n", state.Id, err.Error())
		} else {
			log.Printf("Delete torrent %v\n.", state.Name)
		}
	}
}

// TODO: TEST CASE!!!!
func main() {
	log.Println("tranmission-helper started.")

	defer func() {
		if r := recover(); r != nil {
			log.Println("tranmission-helper exited with error.", r)
		} else {
			log.Println("tranmission-helper completed successfully.")
		}
	}()

	output := getTranmissionRemoteListOutput()

	torrentStates := parseRawOutput(output)
	log.Printf("Found %v torrents\n", len(torrentStates))

	finishedTorrents := filterFinishedTorrents(torrentStates)

	if len(finishedTorrents) == 0 {
		return
	}

	notify(finishedTorrents)
	delete(finishedTorrents)
}
