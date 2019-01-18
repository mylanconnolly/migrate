package migrate

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

// Migration is used to represent a database migration.
type Migration struct {
	down    string
	up      string
	version int64
}

func migrationFromFile(upFilename string) (Migration, error) {
	downFilename := strings.Replace(upFilename, "up.sql", "down.sql", 1)
	up, err := ioutil.ReadFile(upFilename)

	if err != nil {
		return Migration{}, fmt.Errorf("could not read up migration '%s': %s", upFilename, err.Error())
	}
	down, err := ioutil.ReadFile(downFilename)

	if err != nil {
		return Migration{}, fmt.Errorf("could not read down migration '%s': %s", downFilename, err.Error())
	}
	base := filepath.Base(upFilename)
	vstr := strings.Split(base, "_")[0]
	version, err := strconv.ParseInt(vstr, 10, 64)

	if err != nil {
		return Migration{}, err
	}
	return Migration{up: string(up), down: string(down), version: version}, nil
}
