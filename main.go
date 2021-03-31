package main

// import "os/exec"
import "log"
import "strings"
// import "bufio"

type TorrentState struct {
	Name string
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

	lines = lines[1:len(lines)-1]
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
ID     Done       Have  ETA           Up    Down  Ratio  Status       Name
	 29    53%    3.42 GB  Unknown      0.0     0.0    0.0  Idle         test
	 30    n/a    4.21 GB  Done         0.0     0.0   None  Stopped      test 2
Sum:           7.63 GB               0.0     0.0
`, "\n")

	torrentStates := parseRawOutput(output)
	log.Printf("Found %v torrents", len(torrentStates))

	// finishedTorrents = filterFinishedTorrents(torrentStates)

	// notify(finishedTorrents)
	// logNotify(finishedTorrents)
	// delete(finishedTorrents)
	// logDelete(finishedTorrents)
}
