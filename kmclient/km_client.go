package kmclient

import (
	"encoding/gob"
	"errors"
	"github.com/rohanthewiz/serr"
	"kitty_mon_client/auth"
	"kitty_mon_client/config"
	"kitty_mon_client/km_db"
	"kitty_mon_client/message"
	"kitty_mon_client/node"
	"kitty_mon_client/reading"
	"kitty_mon_client/util"
	"net"
	"strconv"
)

func SynchAsClient(host string, serverSecret string) (err error) {
	const stage = "in synch as client"

	conn, err := net.Dial("tcp", host+":"+config.Opts.SynchPort)
	if err != nil {
		util.Lpl(NetworkConnErrorMsg+" "+stage, err)
		return serr.Wrap(err, NetworkConnErrorMsg, "stage", stage)
	}

	defer func() {
		conn.Close()
		if r := recover(); r != nil {
			util.Lpl("Recovered in synch_client", r)
		}
	}()

	msg := message.Message{} // init to empty struct
	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)
	defer message.SendMsg(enc, message.Message{Type: "Hangup"})

	// Send handshake - Client initiates
	message.SendMsg(enc, message.Message{
		Type: "WhoAreYou", Param: auth.WhoAmI(), Param2: serverSecret,
	})
	message.RcxMsg(dec, &msg) // Decode the response

	if msg.Type == "WhoIAm" {
		guid := msg.Param // retrieve the server's guid
		util.Pl("The server's guid is", util.Short_sha(guid))
		if len(guid) != 40 {
			util.Fpl("The server's id is invalid. Run the server once with the -setup_db option")
			return errors.New("The server's id is invalid")
		}
		// Is there an auth token for us?
		if len(msg.Param2) == 40 {
			node.SetNodeToken(guid, msg.Param2) // make sure to save new auth
			// i.e. Given a node (server) with id guid, our auth token on that server is msg.Param2
		}
		// Get the server's node info from our DB
		node, err := node.GetNodeByGuid(guid)
		if err != nil {
			util.Fpl("Error retrieving node object")
			return serr.Wrap(err)
		}
		msg.Param2 = "" // clear for next msg

		// Auth
		msg.Type = "AuthMe"
		msg.Param = node.Token // This is our token for communication with this node (server). It is set by one of two access granting mechanisms

		if config.Opts.NodeName != "" {
			node.Name = config.Opts.NodeName // The server will know this node as this name
			km_db.Db.Save(&node)             // Save it locally
			msg.Param2 = config.Opts.NodeName
		}
		message.SendMsg(enc, msg)
		message.RcxMsg(dec, &msg)
		if msg.Param != "Authorized" {
			util.Fpl("The server declined the authorization request")
			return errors.New("The server declined the authorization request")
		}

		// The Client will send one or more messages to the server
		readings := RetrieveUnsentReadings()
		util.Pf("%d unsent readings found\n", len(readings))
		if len(readings) > 0 {
			message.SendMsg(enc, message.Message{Type: "NumberOfReadings",
				Param: strconv.Itoa(len(readings))})
			message.RcxMsg(dec, &msg)
			if msg.Type == "SendReadings" {
				msg.Type = "Reading"
				msg.Param = ""

				for _, reading := range readings {
					reading.SourceGuid = auth.WhoAmI()
					msg.Reading = reading
					message.SendMsg(enc, msg)
					// Let's go ahead and delete here
					//reading.Sent = 1
					km_db.Db.Delete(&reading) //db.Save(&reading)
				}
			}
		}

	} else {
		util.Fpl("Node does not respond to request for database id")
		util.Fpl("Make sure both server and client databases have been properly setup(migrated) with the -setup_db option")
		util.Fpl("or make sure kitty_mon version is >= 0.9")
		return errors.New("Handshake error")
	}

	util.Lpl("Synch Operation complete")

	return nil
}

func RetrieveUnsentReadings() []reading.Reading {
	var readings []reading.Reading
	if config.Opts.Bogus == false {
		km_db.Db.Where("sent = ?", 0).Order("created_at desc").Limit(config.Opts.L).Find(&readings)
	} else {
		// Send some bogus readings for development
		for i := 0; i < 3; i++ {
			readings = append(readings, reading.BogusReading())
		}
	}
	return readings
}
