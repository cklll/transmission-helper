package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
)

type TorrentState struct {
	Id   string
	Name string
	Done string
}

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

	lines = lines[1 : len(lines)-1]
	for _, line := range lines {
		parts := strings.Fields(line)

		id := parts[0]
		done := parts[1]
		name := line[namePosition:]

		states = append(states, TorrentState{Id: id, Name: name, Done: done})
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
	result := []TorrentState{}

	for _, state := range states {
		if state.Done == "100%" {
			result = append(result, state)
		}
	}

	log.Printf("%v of %v torrents are finished.", len(result), len(states))

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
		message += fmt.Sprintf("%v: %v \r\n", state.Done, state.Name)
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

func delete(torrentStates []TorrentState) {
	// transmission-remote --auth <user>:<password> -t<ID> -r
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
