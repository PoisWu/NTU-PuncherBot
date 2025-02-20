package punchclock_test

import (
	"log"
	"os"
	"os/exec"
	"puncherbot/src/punchclock"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var cfg punchclock.Config

func TestMain(m *testing.M) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	binarygitpath, _ := cmd.Output()
	gitpath := string(binarygitpath)
	// remove the '\n' at the end
	viper.AddConfigPath(gitpath[:len(gitpath)-1])
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Using config file:", viper.ConfigFileUsed(), "  ", err)
	}
	viper.Unmarshal(&cfg)
	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestReadConfig(t *testing.T) {
	assert.Equal(t, cfg.Account.Username, "chengyenwu")
}

func TestAttend(t *testing.T) {
	puncher, err := punchclock.NewPuncher(cfg)
	if err != nil {
		t.Log(err)
	}
	puncher.Attend(false)
}

// Write a test to PunchinOut
func TestLeave(t *testing.T) {
	puncher, err := punchclock.NewPuncher(cfg)
	if err != nil {
		t.Log(err)
	}
	puncher.Leave(false)
}

func TestTodayStatus(t *testing.T) {
	puncher, err := punchclock.NewPuncher(cfg)
	if err != nil {
		t.Log(err)
	}
	puncher.TodayStatus()
}
