package cmd

import (
	"fmt"
	"log"
	"puncherbot/src/punchclock"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

const (
	punchclockGroupID = "Punchclock "
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:     "run [--config config.toml]",
	Short:   "Start the auto punch in/out process",
	Long:    `It starts the auto punch in/out process`,
	GroupID: punchclockGroupID,
	Run: func(cmd *cobra.Command, args []string) {
		var cfg punchclock.Config
		viper.Unmarshal(&cfg)
		p, err := punchclock.NewPuncher(cfg)
		if err != nil {
			log.Fatalf("Not able to initialized punchclock due to %s", err)
		}
		p.Run()
	},
}
var statusCmd = &cobra.Command{
	Use:     "status [--config config.toml]",
	Short:   "Present the clock-in record ",
	Long:    `It present today's clock-in and clock-out record`,
	GroupID: punchclockGroupID,
	Run: func(cmd *cobra.Command, args []string) {
		var cfg punchclock.Config
		fmt.Println(viper.AllSettings())
		viper.Unmarshal(&cfg)
		p, err := punchclock.NewPuncher(cfg)

		if err != nil {
			log.Fatalf("Not able to initialized punchclock due to %s", err)
		}
		p.TodayStatus()
	},
}

func init() {
	rootCmd.AddGroup(
		&cobra.Group{
			ID:    punchclockGroupID,
			Title: "Punchclock function",
		},
	)
	cobra.OnInitialize(initConfig)

	runCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.toml", "config file")
	viper.BindPFlag("config", runCmd.PersistentFlags().Lookup("config"))
	rootCmd.AddCommand(runCmd)

	statusCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.toml", "config file")
	viper.BindPFlag("config", statusCmd.PersistentFlags().Lookup("config"))
	rootCmd.AddCommand(statusCmd)

}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Using config file:", viper.ConfigFileUsed(), "  ", err)
	}
}
