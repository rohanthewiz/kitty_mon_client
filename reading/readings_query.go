package reading

import (
	"kitty_mon_client/km_db"
	"time"
)

func DelOlderThanNDays(nDays int) {
	threshold := time.Now().Add(-time.Duration(24*nDays) * time.Hour)
	km_db.Db.Where("created_at < ?", threshold).Delete(Reading{})
}
