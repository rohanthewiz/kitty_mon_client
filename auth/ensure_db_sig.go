package auth

import (
	"kitty_mon_client/config"
	"kitty_mon_client/km_db"
	"kitty_mon_client/node"
	"kitty_mon_client/util"
)

const serverKey = "km.server"

// Ensure this node has an identity (DB Signature)
// (This is a node entry representing this node (servers act as one node btw) with is_local = 1)
// Servers will share the DBSig -- however server will only use Redis for storage
func EnsureDBSig() { // TODO - let's return an error if something goes wrong here
	// Local node info is cached
	if LocalNode.Id > 0 && len(LocalNode.Token) == 40 {
		return /* all is good */
	}

	var nde node.Node

	if config.Opts.SynchClient != "" { // we are a client
		if km_db.Db.Where("is_local = 1").First(&nde); nde.Id < 1 { // create the signature
			km_db.Db.Create(&node.Node{Guid: util.Random_sha1(), Token: util.Random_sha1(), IsLocal: 1})
			if km_db.Db.Where("is_local = 1").First(&nde); nde.Id > 0 && len(nde.Token) == 40 { // was it saved?
				util.Pl("Local signature created")
			}
		} /*else {
			util.Pl("Local db signature already exists")
		}*/
	}

	// Cache local copy
	LocalNode = nde
}
