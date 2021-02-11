package loaders

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"kitty_mon_client/auth"
	"kitty_mon_client/config"
	"kitty_mon_client/km_db"
	"kitty_mon_client/node"
	"kitty_mon_client/reading"
	"kitty_mon_client/util"
	"os"
)

func DBLoader() (err error) {
	km_db.Db, err = gorm.Open("sqlite3", config.Opts.DbPath)
	if err != nil {
		util.Lf("There was an error connecting to the DB.\nDBPath: " + config.Opts.DbPath)
		os.Exit(2)
	}

	//Do we need to migrate?
	if !km_db.Db.HasTable(&node.Node{}) || !km_db.Db.HasTable(&reading.Reading{}) {
		km_db.Migrate(&reading.Reading{}, &node.Node{})
		auth.EnsureDBSig() // Initialize local with a SHA1 signature if it doesn't already have one
	}

	km_db.Db.LogMode(config.Opts.Debug) // Set debug mode for Gorm db

	if config.Opts.Admin == "delete_tables" {
		fmt.Println("Are you sure you want to delete all data? (N/y)")
		var input string
		fmt.Scanln(&input) // Get keyboard input
		util.Pd("input", input)
		if input == "y" || input == "Y" {
			km_db.Db.DropTableIfExists(&reading.Reading{})
			km_db.Db.DropTableIfExists(&node.Node{})
			util.Pl("Readings tables deleted")
		}
		return
	}

	return
}
