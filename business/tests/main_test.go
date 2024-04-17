package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/EnesDemirtas/medisync/business/data/dbtest"
	"github.com/EnesDemirtas/medisync/foundation/docker"
)

var c *docker.Container

func TestMain(m *testing.M) {
	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}

	os.Exit(code)
}

func run(m *testing.M) (int, error) {
	var err error

	c, err = dbtest.StartDB()
	if err != nil {
		return 1, err
	}
	defer dbtest.StopDB(c)

	return m.Run(), nil
}