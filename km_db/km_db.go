package km_db

import (
	"github.com/jinzhu/gorm"
	"kitty_mon_client/util"
)

// Init db
var Db gorm.DB

func Migrate(what interface{}, node interface{}) {
	// Create or update the table structure as needed
	util.Pl("Migrating the DB...")
	Db.AutoMigrate(node)
	Db.AutoMigrate(what)
	//According to GORM: Feel free to change your struct, AutoMigrate will keep your database up-to-date.
	// Fyi, AutoMigrate will only *add new columns*, it won't update column's type or delete unused columns for safety
	// If the table is not existing, AutoMigrate will create the table automatically.
	Db.Model(what).AddUniqueIndex("idx_reading_guid", "guid")
	Db.Model(what).AddIndex("idx_reading_source_guid", "source_guid")
	Db.Model(what).AddIndex("idx_reading_created_at", "created_at")
	Db.Model(node).AddUniqueIndex("idx_node_guid", "guid")
	// This would disallow blanks //db.Model(&Node{}).AddUniqueIndex("idx_node_name", "name")

	util.Pl("Migration complete")
}
