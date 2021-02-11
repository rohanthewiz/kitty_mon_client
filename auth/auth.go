package auth

import (
	"kitty_mon_client/node"
)

const AuthFailMsg = "Authentication failure. Generate authorization token with -synch_auth\nThen store in node entry on client with -store_synch_auth"

var LocalNode node.Node // cache the local node

// The low-down on auth.
// Each node will have a Node table
// Each record will store the GUID of a node, and an associated token
//  which is the token required to authenticate with the node,
// or the server's secret token if the node is the local machine
//  depending on the setting of the is_local field
//

// Get local DB signature
func WhoAmI() string {
	var node node.Node
	var err error
	if node, err = GetLocalNode(); err != nil {
		return ""
	}
	return node.Guid
}

// LocalNode should be populated at startup
// so we try to get it from the global variable
func GetLocalNode() (node.Node, error) {
	if LocalNode.Id > 0 { // it has been inited
		return LocalNode, nil
	}

	EnsureDBSig() // this will populate LocalNode

	return LocalNode, nil // TODO - better error handling here
}
