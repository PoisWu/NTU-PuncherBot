package db_test

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"puncherbot/src/db"
	"testing"
)

var dayOffDB *db.DayOffDB
var err error

func TestMain(m *testing.M) {
	dayOffDB, err = db.NewDayOffDB("test.db")
	if err != nil {
		panic(err)
	}
	defer dayOffDB.Close()
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestLoadCalendar(t *testing.T) {

	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	binarygitpath, _ := cmd.Output()
	gitpath := string(binarygitpath)
	// remove the '\n' at the end
	err := dayOffDB.LoadCalendar(fmt.Sprintf("%s/%s", gitpath[:len(gitpath)-1], "2025.json"))
	if err != nil {
		log.Fatal(err)
	}
}

func TestDateStatus(t *testing.T) {
	dayoff, description, err := dayOffDB.DateStatusToday()
	t.Logf("\nError: %s\nDay off?: %v\ndescription: %s\n", err, dayoff, description)
}

func TestWork(t *testing.T) {
	err := dayOffDB.Work("2025-02-22", "Work good")
	if err != nil {
		log.Fatal(err)
	}
	dayoff, description, err := dayOffDB.DateStatus("2025-02-22")
	t.Logf("\nError: %s\nDay off?: %v\ndescription: %s\n", err, dayoff, description)

}

func TestLeave(t *testing.T) {
	err := dayOffDB.Leave("2025-02-20", "Leave")
	if err != nil {
		log.Fatal(err)
	}
	dayoff, description, err := dayOffDB.DateStatus("2025-02-20")
	t.Logf("\nError: %s\nDay off?: %v\ndescription: %s\n", err, dayoff, description)
}
