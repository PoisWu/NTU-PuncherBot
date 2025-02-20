package cmd

import (
	"fmt"
	"log"
	"puncherbot/src/db"
	"time"

	"github.com/spf13/cobra"
)

const (
	dbGroupID = "Database-management"
)

// datestatusCmd represents the datestatus command
var datestatusCmd = &cobra.Command{
	Use:   "datestatus",
	Short: "Check if a given date is a day off",
	Long: `It will send a query to the calendar
database to know if a given date is a day off.`,
	GroupID: dbGroupID,
	Run: func(cmd *cobra.Command, args []string) {
		date, _ := cmd.Flags().GetString("date")
		if date == "" {
			date = time.Now().Format(db.DateFormat)
		}
		if !checkDate(date) {
			log.Fatalf("%s doesn't follow the YYYY-MM-DD format or invalid date\n", date)
		}
		db, err := db.NewDayOffDB(db.DBName)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		isdayoff, description, err := db.DateStatus(date)
		if err != nil {
			fmt.Println(err)
			return
		} else {
			fmt.Printf(`The date %s
dayoff: %v
Description: %s
`, date, isdayoff, description)
		}
	},
}

var loadCmd = &cobra.Command{
	Use:     "load <path>",
	Short:   "Load a calendar into the calendar database",
	Long:    `It loads the given calendar into the calendar database`,
	GroupID: dbGroupID,
	Args:    cobra.ExactArgs(1),
	Example: "  puncherbot load 2025.json",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := db.NewDayOffDB(db.DBName)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		calendarPath := args[0]
		err = db.LoadCalendar(calendarPath)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var workCmd = &cobra.Command{
	Use:     "work <-d YYYY-MM-DD> [-m Reason]",
	Short:   "Set a date as workday",
	Long:    `It will update the calendar database and set <YYYY-MM-DD> as work day`,
	GroupID: dbGroupID,
	Run: func(cmd *cobra.Command, args []string) {
		date, _ := cmd.Flags().GetString("date")
		descrption, _ := cmd.Flags().GetString("message")
		if !checkDate(date) {
			log.Fatalf("%s doesn't follow the YYYY-MM-DD format or invalid date\n", date)
		}
		// Open a DayOffDB
		db, err := db.NewDayOffDB(db.DBName)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		err = db.Work(date, descrption)

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Update done")
		}

	},
}
var leaveCmd = &cobra.Command{
	Use:     "leave <-d YYYY-MM-DD> [-m Reason]",
	Short:   "Set a date as day off",
	Long:    `It will update the calendar database and set <YYYY-MM-DD> as day off`,
	GroupID: dbGroupID,
	Run: func(cmd *cobra.Command, args []string) {
		date, _ := cmd.Flags().GetString("date")
		descrption, _ := cmd.Flags().GetString("message")
		if !checkDate(date) {
			log.Fatalf("%s doesn't follow the YYYY-MM-DD format or invalid date\n", date)
		}
		db, err := db.NewDayOffDB(db.DBName)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		err = db.Leave(date, descrption)

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Update done")
		}
	},
}

func init() {
	rootCmd.AddGroup(
		&cobra.Group{
			ID:    dbGroupID,
			Title: "Calendar database management",
		},
	)
	datestatusCmd.Flags().StringP("date", "d", "", "specify the date in YYYY-MM-DD format (default is today's date)")

	workCmd.Flags().StringP("date", "d", "", "specify the date in YYYY-MM-DD format")
	leaveCmd.Flags().StringP("date", "d", "", "specify the date in YYYY-MM-DD format")
	workCmd.MarkFlagRequired("date")
	leaveCmd.MarkFlagRequired("date")
	workCmd.Flags().StringP("message", "m", "working is good!", "specify the reason")
	leaveCmd.Flags().StringP("message", "m", "I want to rest...", "specify the reason")

	rootCmd.AddCommand(datestatusCmd, loadCmd)
	rootCmd.AddCommand(workCmd, leaveCmd)
}

func checkDate(date string) bool {
	layout := db.DateFormat
	_, err := time.Parse(layout, date)
	return err == nil
}
