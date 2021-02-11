package config

import (
	"flag"
	"os"
	"strings"
	"time"
)

const App_name = "Kitty Monitor Client"
const Version string = "0.2.0"

var ReadingsProdPollRate = 4 * time.Minute
var ReadingsDevPollRate = 8 * time.Second

var Opts *Option //Cmdline options and flags

type Option struct {
	Q               string
	Port            string
	SynchPort       string
	Admin           string
	Db              string
	DbPath          string
	Exp             string
	Imp             string
	SynchClient     string // client asking to synch
	GetTokenForNode string
	SaveNodeToken   string
	ServerSecret    string
	Env             string
	NodeName        string
	L               int
	Ql              bool
	V               bool
	WhoAmI          bool
	Del             bool
	Svr             bool
	GetServerSecret bool
	SynchServer     bool
	SetupDb         bool
	Sep             string
	Verbose         bool
	Debug           bool
	Bogus           bool
}

//Setup commandline options and other configuration for Go Notes
func NewOpts() *Option {
	opts := new(Option)

	// flag.{String|Bool|Int|Float...}( the_option, default_value, description )
	qPtr := flag.String("q", "", "Query for notes based on a LIKE search. \"all\" will return all notes")
	pPtr := flag.String("port", "9080", "Specify webserver port")
	synchPortPtr := flag.String("synch_port", "1470", "Specify webserver port")
	adminPtr := flag.String("admin", "", "Privileged actions like 'delete_table'")
	dbPtr := flag.String("db", "", "Sqlite DB path")
	expPtr := flag.String("exp", "", "Export the notes queried to the format of the file given")
	impPtr := flag.String("imp", "", "Import the notes queried from the file given")
	synchClientPtr := flag.String("synch_client", "", "Synch client mode")
	getTokenForNodePtr := flag.String("get_node_token", "", "Get a token for interacting with this as server")
	saveNodeTokenPtr := flag.String("save_node_token", "", "Save a token for interacting with this as server")
	envPtr := flag.String("env", "dev", "App Environment (dev|prod)")
	nodeNamePtr := flag.String("node_name", "", "Upsert node name on server")

	lPtr := flag.Int("l", 64, "Limit the number of readings returned")
	qlPtr := flag.Bool("ql", false, "Query for the last note updated")
	vPtr := flag.Bool("v", false, "Show version")
	whoamiPtr := flag.Bool("whoami", false, "Show Client GUID")
	setupDBPtr := flag.Bool("setup_db", false, "Setup the Database")
	delPtr := flag.Bool("del", false, "Delete the notes queried")
	svrPtr := flag.Bool("svr", false, "Web server mode")
	verbosePtr := flag.Bool("verbose", true, "verbose mode")                   // Todo - turn off for production
	debugPtr := flag.Bool("debug", true, "debug mode")                         // Todo - turn off for production
	bogusPtr := flag.Bool("bogus", false, "generate bogus data in production") // Todo - turn off for production

	flag.Parse()

	// Store options in a couple of maps
	opts.Q = *qPtr
	opts.Port = *pPtr
	opts.SynchPort = *synchPortPtr
	opts.Admin = *adminPtr
	opts.Db = *dbPtr
	opts.Exp = *expPtr
	opts.Imp = *impPtr
	opts.SynchClient = *synchClientPtr
	opts.GetTokenForNode = *getTokenForNodePtr
	opts.SaveNodeToken = *saveNodeTokenPtr
	opts.Env = *envPtr
	opts.NodeName = *nodeNamePtr
	opts.L = *lPtr
	opts.Ql = *qlPtr
	opts.V = *vPtr
	opts.WhoAmI = *whoamiPtr
	opts.Del = *delPtr
	opts.Svr = *svrPtr
	opts.SetupDb = *setupDBPtr
	opts.Verbose = *verbosePtr
	opts.Debug = *debugPtr
	opts.Bogus = *bogusPtr

	_ = ReadEnvFile() // Pick up some configs from env vars

	separator := "/"
	if strings.Contains(strings.ToUpper(os.Getenv("OS")), "WINDOWS") {
		separator = "\\"
	}
	opts.Sep = separator

	db_file := "kitty_mon_client.sqlite"

	var db_folder string
	var db_full_path string

	if len(*dbPtr) == 0 {
		if len(os.Getenv("HOME")) > 0 {
			db_folder = os.Getenv("HOME")
		} else if len(os.Getenv("HOMEDRIVE")) > 0 && len(os.Getenv("HOMEPATH")) > 0 {
			db_folder = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		} else {
			db_folder = separator /// last resort
		}
		db_full_path = db_folder + separator + db_file
	} else {
		db_full_path = *dbPtr
	}
	opts.DbPath = db_full_path

	return opts
}
