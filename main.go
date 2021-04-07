package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
)

type TorrentState struct {
	Name   string
	Status string
}

// =========
// ID     Done       Have  ETA           Up    Down  Ratio  Status       Name
//   29    53%    3.42 GB  Unknown      0.0     0.0    0.0  Idle         test
//   30    n/a    4.21 GB  Done         0.0     0.0   None  Stopped      test 2
// Sum:           7.63 GB               0.0     0.0
// =========
// the status and name should start at the same character as the header
// so we get the start position of "Status" and "Name" and use it to retrieve the string
func parseRawOutput(output string) []TorrentState {
	states := []TorrentState{}

	lines := strings.Split(output, "\n")
	headerLine := lines[0]
	statusPosition := strings.Index(headerLine, "Status")
	namePosition := strings.Index(headerLine, "Name")

	lines = lines[1 : len(lines)-1]
	for _, line := range lines {
		status := line[statusPosition:namePosition]
		status = strings.Trim(status, " ")

		name := line[namePosition:]
		name = strings.Trim(name, " ")

		states = append(states, TorrentState{name, status})
	}

	return states
}

// TODO
// func getRawTorrentStates() string {
// 	auth_credentials := ""

// 	cmd := exec.Command("transmission-remote", "--auth", auth_credentials, "-l")
// 	stdout, err := cmd.Output()

// 	if err != nil {
// 		log.Println(err.Error())
// 		panic("Cannot get output from transmission-remote.")
// 	}

// 	return string(stdout)
// }

func filterFinishedTorrents(states []TorrentState) []TorrentState {
	stateMap := map[string][]TorrentState{}

	for _, state := range states {
		stateMap[state.Status] = append(stateMap[state.Status], state)
	}

	result := []TorrentState{}
	for status, states := range stateMap {
		// !!!!! FIXME !!!!
		// TODO: Need to fix, when in "seeding to 0 of 0 peers"
		// the output is
		// 30   100%    5.39 GB  Done         0.0     0.0    0.0  Idle
		// We MUST use "Done: 100%" to check

		// TODO: haven't manually verified the keywords
		// maybe instead check "Done" percentage
		// "Finished" - https://github.com/transmission/transmission/blob/8566df069899ce8923463cadeb0ff66d4544991a/utils/remote.c#L844
		// "Seeding" - https://github.com/transmission/transmission/blob/8566df069899ce8923463cadeb0ff66d4544991a/utils/remote.c#L898
		if status == "Finished" {
			result = append(result, states...)
		} else if status == "Seeding" {
			result = append(result, states...)
		}

		// TODO: want to test this but not sure how
		log.Printf("Found %v torrents with %v status.", len(states), status)
	}

	return result
}

// TODO: TEST CASE!!!!
// Need refactoring to make this testable
func notify(finishedTorrentStates []TorrentState) {
	log.Printf("Sending emails with %v finished torrents.", len(finishedTorrentStates))

	recipientsString := os.Getenv("TH_NOTIFY_EMAILS")
	recipients := strings.Split(recipientsString, ",")

	subject := fmt.Sprintf("[transmission-helper] %v torrents completed.", len(finishedTorrentStates))
	message := ""
	for _, state := range finishedTorrentStates {
		message += fmt.Sprintf("%v: %v \r\n", state.Status, state.Name)
	}

	mailNotifier := MailNotifier{smtp.SendMail}
	mailConfig := mailNotifier.GetMailConfig()
	err := mailNotifier.Send(mailConfig, subject, message, recipients)
	if err != nil {
		log.Println("Failed to send emails.")
		log.Println(err)
	} else {
		log.Printf("Sent emails to %v recipients.", len(recipients))
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

	// output := getRawTorrentStates()
	output := strings.Trim(`
ID     Done       Have  ETA           Up    Down  Ratio   Status       Name
   29    53%    3.42 GB  Unknown      0.0     0.0    0.0  Idle         test
   30    n/a    4.21 GB  Done         0.0     0.0   None  Stopped      test 2
   31    n/a    4.21 GB  Done         0.0     0.0   None  Finished     test 3
Sum:           7.63 GB               0.0     0.0
`, "\n")

	torrentStates := parseRawOutput(output)
	log.Printf("Found %v torrents", len(torrentStates))

	finishedTorrents := filterFinishedTorrents(torrentStates)

	if len(finishedTorrents) == 0 {
		return
	}

	notify(finishedTorrents)
	// delete(finishedTorrents)
}
