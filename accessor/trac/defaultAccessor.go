package trac

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"stevejefferson.co.uk/trac2gitea/log"

	"github.com/go-ini/ini"
	_ "github.com/mattn/go-sqlite3" // sqlite database driver
)

// DefaultAccessor is the default implementation of the trac Accessor interface, accessing Trac via its database and filestore.
type DefaultAccessor struct {
	rootDir string
	db      *sql.DB
	config  *ini.File
}

// CreateDefaultAccessor creates a new Trac accessor.
func CreateDefaultAccessor(tracRootDir string) *DefaultAccessor {
	stat, err := os.Stat(tracRootDir)
	if err != nil {
		log.Fatal(err)
	}
	if stat.IsDir() != true {
		log.Fatalf("Trac root directory %s is not a directory\n", tracRootDir)
	}

	tracIniPath := fmt.Sprintf("%s/conf/trac.ini", tracRootDir)
	stat, err = os.Stat(tracIniPath)
	if err != nil {
		log.Fatal(err)
	}

	tracConfig, err := ini.Load(tracIniPath)
	if err != nil {
		log.Fatal(err)
	}

	accessor := DefaultAccessor{db: nil, rootDir: tracRootDir, config: tracConfig}

	// extract path to trac DB - currently sqlite-specific...
	tracDatabaseString := accessor.GetStringConfig("trac", "database")
	tracDatabaseSegments := strings.SplitN(tracDatabaseString, ":", 2)
	tracDatabasePath := tracDatabaseSegments[1]
	if !filepath.IsAbs(tracDatabasePath) {
		tracDatabasePath = filepath.Join(tracRootDir, tracDatabasePath)
	}

	log.Infof("Using trac database %s\n", tracDatabasePath)

	tracDb, err := sql.Open("sqlite3", tracDatabasePath)
	if err != nil {
		log.Fatal(err)
	}
	accessor.db = tracDb

	return &accessor
}