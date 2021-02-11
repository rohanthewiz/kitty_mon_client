package reading

import (
	"kitty_mon_client/auth"
	"kitty_mon_client/config"
	"kitty_mon_client/km_db"
	node2 "kitty_mon_client/node"
	"kitty_mon_client/util"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const redisPrefix = "reading:"

var CriticalTemp int = 75000

// TODO - move this out of init to a more predictable place
func init() {
	envCtemp, err := strconv.Atoi(os.Getenv("CRITICAL_TEMP"))
	if err != nil || envCtemp == 0 {
		// panic("CRITICAL_TEMP env var is not a valid positive integer - e.g. 70")
	} else {
		CriticalTemp = 1000 * envCtemp
	}
}

// This represents the payload sent to the server
// This is the equivalent of a NoteChange in GoNotes
type Reading struct {
	Id                   int64
	Guid                 string    `sql:"size:40" json:"guid"`       // random id for each message
	SourceGuid           string    `sql:"size:40" json:"sourceGUID"` // We will tag this with the client's db sig when reading sent
	IPs                  string    `sql:"size:254" json:"IPs"`
	Sent                 int       `json:"-"`                    // has the reading been sent // bool 0 - false, 1 - true
	Temp                 int       `json:"temp"`                 // temperature
	MeasurementTimestamp time.Time `json:"measurementTimestamp"` // True CreatedAt for the reading,
	// since GORM also updates CreatedAt when saved on the server
	CreatedAt time.Time `json:"createdAt"` // GORM automatically updates this field on save
}

type ReadingEnriched struct {
	Reading
	Name   string
	Status string
}

type ByCreatedAt []Reading

func (ncs ByCreatedAt) Len() int {
	return len(ncs)
}
func (ncs ByCreatedAt) Less(i, j int) bool {
	return ncs[i].CreatedAt.Before(ncs[j].CreatedAt)
}
func (ncs ByCreatedAt) Swap(i, j int) {
	ncs[i], ncs[j] = ncs[j], ncs[i]
}

func (r Reading) Save() bool {
	r.Id = 0 // Make sure the reading has a zero id for db creation
	// A non-zero Id will not be created
	km_db.Db.Create(&r)         // will auto create contained objects too and it's smart - 'nil' children will not be created :-)
	if !km_db.Db.NewRecord(r) { // was it saved?
		util.Pl("Reading saved:", util.Short_sha(r.Guid))
		return true
	}
	util.Fpl("Failed to save reading:", util.Short_sha(r.Guid))
	return false
}

func (r *Reading) Print() {
	util.Pf("%+v\n", r)
}

func BogusReading() Reading {
	return Reading{
		Guid:       util.Random_sha1(),
		SourceGuid: auth.WhoAmI(),
		Temp:       rand.Intn(100),
	}
}

func FindReadingById(id int64) Reading {
	var reading Reading
	km_db.Db.First(&reading, id)
	return reading
}

func PollTemp() {
	var reading Reading
	for {
		var wait time.Duration
		if config.Opts.Env == "dev" {
			wait = 8 * time.Second
		} else {
			wait = 2 * time.Minute
		}
		time.Sleep(wait)

		if strings.ToLower(os.Getenv("KM_SHUTDOWN")) == "true" { // graceful shutdown
			break
		}

		// Temperature
		if config.Opts.Bogus {
			reading = BogusReading()
		} else {
			reading = Reading{
				Guid:                 util.Random_sha1(),
				SourceGuid:           auth.WhoAmI(),
				IPs:                  util.IPs(true),
				Temp:                 CatTemp(),
				MeasurementTimestamp: time.Now(),
				Sent:                 0,
			}
		}
		km_db.Db.Save(&reading)
		// Cleanup
		DelOlderThanNDays(7) // TODO add config for this
	}
}

func CatTemp() int {
	cmdArgs := []string{"/sys/class/thermal/thermal_zone0/temp"}
	byteTemps, err := exec.Command("cat", cmdArgs...).Output()
	if err != nil {
		util.Lpl("Error acquiring temperature.")
		return -255
	}
	strTemp := strings.Trim(string(byteTemps), " \n\t") // clean up whitespace
	var temp int
	temp, err = strconv.Atoi(strTemp)
	if err != nil {
		util.Lpl("Error converting temperature.")
		return -255
	}
	return temp
}

func ReadingsPopulateNodeName(readings []Reading) (readingsWithNames []ReadingEnriched) {
	guidToName := make(map[string]string, 4)

	for _, r := range readings {
		rwn := ReadingEnriched{Status: ReadingStatus(r)}
		rwn.SourceGuid = r.SourceGuid
		rwn.Temp = r.Temp
		rwn.Guid = r.Guid
		rwn.Id = r.Id
		rwn.CreatedAt = r.CreatedAt
		rwn.IPs = r.IPs
		rwn.MeasurementTimestamp = r.MeasurementTimestamp

		if name, ok := guidToName[r.SourceGuid]; ok { // Is it in the cache?
			rwn.Name = name
		} else {
			node, err := node2.GetNodeByGuid(r.SourceGuid)
			if err != nil {
				util.Lpl(err)
			} else {
				rwn.Name = node.Name
				guidToName[r.SourceGuid] = node.Name // also load into cache
			}
		}

		readingsWithNames = append(readingsWithNames, rwn)
	}

	return
}

func ReadingStatus(reading Reading) (status string) {
	temp := reading.Temp / 1000
	switch {
	case temp > 69:
		status = "fail"
	case temp > 65:
		status = "hot"
	case temp > 60:
		status = "warm"
	default:
		status = "good"
	}
	return
}
