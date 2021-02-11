package main

import (
	"kitty_mon_client/auth"
	"kitty_mon_client/config"
	"kitty_mon_client/km_db"
	"kitty_mon_client/kmclient"
	"kitty_mon_client/loaders"
	"kitty_mon_client/node"
	"kitty_mon_client/reading"
	"kitty_mon_client/unloaders"
	"kitty_mon_client/util"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rohanthewiz/serr"
)

func main() {

	loaders.ConfigLoader()
	_ = loaders.DBLoader()

	// ----- UTILITY OPTIONS -----

	// Return our db signature
	if config.Opts.WhoAmI {
		util.Fpl(auth.WhoAmI())
		return
	}

	// TODO let's modify this approach as it involves stopping the server
	// TODO Rather, let's use some kind of admin login approach
	// TODO - First let's put this in a function
	// Server - Generate an auth token for a client
	// The format of the generated token is: server_id-auth_token_for_the_client
	if config.Opts.GetTokenForNode != "" {
		pt, err := node.GetNodeToken(config.Opts.GetTokenForNode)
		if err != nil {
			util.Fpl("Error retrieving token")
			return
		}
		util.Fpf("Node token is: %s-%s\nYou will now need to run the client with \n'kitty_mon -save_node_token the_token'\n",
			auth.WhoAmI(), pt)
		return
	}

	// Client - Save a token generated for us by a server
	if config.Opts.SaveNodeToken != "" {
		node.SaveNodeToken(config.Opts.SaveNodeToken)
		return
	}

	if config.Opts.SetupDb { // Migrate the DB
		km_db.Migrate(&reading.Reading{}, &node.Node{})

		auth.EnsureDBSig() // Initialize local with a SHA1 signature if it doesn't already have one
		return
	}

	go reading.PollTemp() // save temp, whether real or bogus to local db

	// go snapshots.RunSnapshotLoop(snapshotsStopChan, snapshotsDoneChan)

	wait := config.ReadingsProdPollRate
	if config.Opts.Env == "dev" {
		wait = config.ReadingsDevPollRate
	}

	util.Lpl("I will periodically send data to server...")

	networkErrCount := 0

	for {
		if networkErrCount > 3 {
			_ = os.Setenv("KM_SHUTDOWN", "true") // let everyone know we are shutting down
			_ = unloaders.Reboot()
			break
		}

		// The app behavior can be dynamically changed via env vars
		if strings.ToLower(os.Getenv("KM_SHUTDOWN")) == "true" {
			break
		}

		if strRate := strings.ToLower(os.Getenv("KM_READINGS_POLLRATE")); strRate != "" {
			rate, err := strconv.Atoi(strRate)
			if err != nil {
				util.Lpl("Error converting readings pollrate from env var (KM_READINGS_POLLRATE) " + err.Error())
			} else {
				wait = time.Duration(rate) * time.Second
			}
		}

		time.Sleep(wait)

		err := kmclient.SynchAsClient(config.Opts.SynchClient, config.Opts.ServerSecret)
		if ser, ok := err.(serr.SErr); ok {
			mp := ser.FieldsMap()
			if str, ok := mp["msg"]; ok && strings.Contains(str, kmclient.NetworkConnErrorMsg) {
				networkErrCount++
			}
		} else {
			networkErrCount--
			if networkErrCount < 0 {
				networkErrCount = 0
			}
		}
	}

	// Disable image snapshots for now
	// snapshotsStopChan := make(chan bool)
	// snapshotsDoneChan := make(chan bool)
}
