package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	schema = `
    CREATE TABLE IF NOT EXISTS DayOffTable (
        Date TEXT PRIMARY KEY,
        IsDayOff INTEGER NOT NULL,
        Description TEXT
    )
    `
	DBName     = `dayofftable.db`
	DateFormat = "2006-01-02"
)

type calendarEntry []struct {
	Date        string `json:"date"`
	Week        string `json:"week"`
	IsHoliday   bool   `json:"isHoliday"`
	Description string `json:"description"`
}

// DayOffDB wrap the method to interact with the database: iedayoftable.db
type DayOffDB struct {
	dbPath string
	db     *sql.DB
}

func (dayOffDB *DayOffDB) Close() {
	dayOffDB.db.Close()
}

// Open Database and create dayOff table
func NewDayOffDB(dbName string) (*DayOffDB, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	dbPath := fmt.Sprintf("%s/.cache/%s", homedir, dbName)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err != nil {
			db.Close()
		}
	}()
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(schema)
	if err != nil {
		log.Fatal(err)
	}
	dayOffDB := &DayOffDB{
		dbPath: dbPath,
		db:     db,
	}
	return dayOffDB, nil
}

// LoadCalendar should load only one year calendar
func (dayOffDB *DayOffDB) LoadCalendar(calendarPath string) error {
	// Open the calendar
	jsonFile, err := os.Open(calendarPath)

	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	// Read from the calendar

	calendarValue, _ := io.ReadAll(jsonFile)

	// Parsing json into calendarEntry Strucs
	var calendarEntries calendarEntry

	if err := json.Unmarshal(calendarValue, &calendarEntries); err != nil {
		return err
	}
	for _, entry := range calendarEntries {
		// Parse the entry.Data into YYYY-MM-DD
		date := fmt.Sprintf("%s-%s-%s",
			entry.Date[0:4], entry.Date[4:6], entry.Date[6:8])
		var isDayOff int
		if entry.IsHoliday {
			isDayOff = 1
		} else {
			isDayOff = 0
		}
		dayOffDB.insertEntry(date, isDayOff, entry.Description)
	}

	return nil
}

// Given a date in the calendar, return whether isDayOff
func (dayOffDB *DayOffDB) IsDayOff(date string) (bool, error) {
	var dayOff bool
	var err error
	dayOff, _, err = dayOffDB.DateStatus(date)

	if err != nil {
		return false, err
	} else {
		return dayOff, nil
	}
}
func (dayOffDB *DayOffDB) IsDayOffToday() (bool, error) {
	t := time.Now()
	dateString := t.Format(DateFormat)
	return dayOffDB.IsDayOff(dateString)
}

func (dayOffDB *DayOffDB) DateStatus(date string) (bool, string, error) {
	var dayOff int
	var description string

	// Need to learn how to query a entry
	err := dayOffDB.dateStatus(date, &dayOff, &description)

	if err == sql.ErrNoRows {
		return false, "", err
	} else if err != nil {
		return false, "", err
	} else {
		return dayOff == 1, description, nil
	}
}

func (dayOffDB *DayOffDB) DateStatusToday() (bool, string, error) {
	t := time.Now()
	dateString := t.Format(DateFormat)
	return dayOffDB.DateStatus(dateString)
}

// Set `date` as workday
func (dayOffDB *DayOffDB) Work(date, description string) error {
	return dayOffDB.edit(date, description, 0)
}

// Set `date` as leave day
func (dayOffDB *DayOffDB) Leave(date, description string) error {
	return dayOffDB.edit(date, description, 1)
}

func (dayOffDB *DayOffDB) dateStatus(date string, dayOff *int, description *string) error {
	query := `SELECT IsDayOff, Description FROM DayOffTable WHERE Date = $1`
	err := dayOffDB.db.QueryRow(query, date).Scan(dayOff, description)
	return err
}

func (dayOffDB *DayOffDB) insertEntry(date string, isDayOff int, description string) error {
	query := `INSERT INTO DayOffTable (Date,  IsDayOff, Description)
        VALUES($1, $2, $3)
    `
	err := dayOffDB.db.QueryRow(query, date, isDayOff, description).Scan()
	return err
}

func (dayOffDB *DayOffDB) edit(date, description string, isDayOff int) error {

	stmt, err := dayOffDB.db.Prepare(`UPDATE DayOffTable SET IsDayOff = ?, Description = ? WHERE Date = ?`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(isDayOff, description, date)
	return err
}
