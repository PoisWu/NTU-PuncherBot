package punchclock

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"puncherbot/src/db"
	"puncherbot/src/logger"
	"time"
)

var (
	readerAttend   = bytes.NewReader([]byte(`type=6&t=1`))
	readerLeave    = bytes.NewReader([]byte(`type=6&t=2`))
	readerGetToday = bytes.NewReader([]byte(`type=3`))
	readerGetPast  = bytes.NewReader([]byte(`type=4&day=7`))
)

const (
	t_attend   int = 1
	t_leave    int = 2
	t_getToday int = 3
	t_getPast  int = 4
)

type Tick struct {
	ticker    *time.Ticker
	executing bool
}

type Puncher struct {
	puncher  *http.Client
	db       *db.DayOffDB
	myLog    *logger.MyLogger
	username string
	password string
}

type Config struct {
	Account struct {
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"account"`
	Telegram struct {
		Telegram_chat_id       string `mapstructure:"chat_id"`
		Telegram_chatbot_token string `mapstructure:"chatbot_token"`
	} `mapstructure:"telegram"`
}

func NewPuncher(cfg Config) (*Puncher, error) {
	client, err := getClient(false)
	if err != nil {
		log.Fatalf("Create error run into problem - %s", err)
	}
	db, err := db.NewDayOffDB(db.DBName)
	if err != nil {
		log.Fatalf("Cannot due connect to the DayOffDB due to %s", err)
	}

	return &Puncher{
		puncher: client,
		db:      db,
		myLog: logger.NewMyLogger(cfg.Telegram.Telegram_chat_id,
			cfg.Telegram.Telegram_chatbot_token),
		username: cfg.Account.Username,
		password: cfg.Account.Password,
	}, nil
}

func (p *Puncher) login() {
	p.myLog.Debug("Loging in...")

	// Get the session cookies
	req1, err := http.NewRequest(http.MethodGet, "https://my.ntu.edu.tw/attend/ssi.aspx?type=login", nil)
	if err != nil {
		p.myLog.Fatal("Parsing MyNTU portail Url runs into error ")
	}
	set_request_header(req1)
	_, err = p.puncher.Do(req1)

	if err != nil {
		p.myLog.Fatal("Not being able to send request to MyNTU protail: ", err)
	}

	// Authentication
	payload_login := fmt.Sprintf("user=%s&pass=%s&Submit=登入", p.username, p.password)
	readerLogin := bytes.NewReader([]byte(payload_login))

	req, err := http.NewRequest(http.MethodPost, "https://web2.cc.ntu.edu.tw/p/s/login2/p1.php", readerLogin)

	if err != nil {
		p.myLog.Fatal("Error: Parsing MyNTU login Url runs into error - ", err)
	}
	set_request_header(req)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err = p.puncher.Do(req)
	if err != nil {
		p.myLog.Fatal("Error: Authentication fail due to error - ", err)
	}

	p.myLog.Debug("Login successfully")
}

func (p *Puncher) logout() {
	p.myLog.Debug("Loging out...")
	req, err := http.NewRequest(http.MethodGet, "https://my.ntu.edu.tw/attend/ssi.aspx?type=logout", nil)

	if err != nil {
		p.myLog.Fatal("Parsing MyNTU logout Url runs into error ")
	}

	set_request_header(req)
	_, err = p.puncher.Do(req)
	if err != nil {
		p.myLog.Fatal("Error : Logout run into error ", err)
	}

	// Empty the cookiejar.
	jar, _ := cookiejar.New(nil)
	p.puncher.Jar = jar

	log.Println("Log out sucessfully")
}

func (p *Puncher) Attend(need_wait bool) {
	p.login()
	defer p.logout()

	if need_wait {
		wait_a_while(5, 30)
	}
	p.myLog.Debug("Punching In")
	resHTTP, err := p.request(t_attend)
	if err != nil {
		p.myLog.Fatal("Punch in run into error - ", err)
	}
	type jsonResponse struct {
		T   string `json:"t"`
		Msg string `json:"msg"`
	}
	body, _ := io.ReadAll(resHTTP.Body)
	var res jsonResponse
	if err := json.Unmarshal(body, &res); err != nil {
		p.myLog.Warn("Can not unmarshal JSON due to ", err)
	}
	p.myLog.Info("Attend sent, get response", res.Msg)
}

func (p *Puncher) Leave(need_wait bool) {
	p.login()
	defer p.logout()
	// Wait a while
	if need_wait {
		wait_a_while(5, 30)
	}
	p.myLog.Debug("Punching out")

	// Send a leave request to the server
	resHTTP, err := p.request(t_leave)
	if err != nil {
		p.myLog.Fatal("Punching out run into error - ", err)
	}

	// Parsing the response to JSON from the server
	type jsonResponse struct {
		T   string `json:"t"`
		Msg string `json:"msg"`
	}
	body, _ := io.ReadAll(resHTTP.Body)
	var res jsonResponse
	if err := json.Unmarshal(body, &res); err != nil {
		p.myLog.Warn("Can not unmarshal JSON due to ", err)
	}
	p.myLog.Info("Leave request sent, get response", res.Msg)
}

func (p *Puncher) TodayStatus() {
	p.login()
	defer p.logout()
	p.myLog.Debug("Get Today Status")
	resHTTP, err := p.request(t_getToday)
	if err != nil {
		p.myLog.Fatal("Punch in run into error - ", err)
	}
	type jsonResponse struct {
		D string `json:"d"`
	}
	body, _ := io.ReadAll(resHTTP.Body)
	var res jsonResponse
	if err := json.Unmarshal(body, &res); err != nil {
		p.myLog.Warn("Can not unmarshal JSON due to ", err)
	}
	p.myLog.Info("Today's Punch In/Out status*", res.D, "*")
}

func (p *Puncher) run(tick *Tick, currentTime time.Time) {
	p.myLog.Debug("Starting puncher.... ")
	if tick.executing {
		return
	}

	tick.executing = true

	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		p.myLog.Fatal("Time Zone not correct")
	}
	currentTime = currentTime.In(loc)
	// Get the current time
	cur_hour, cur_min, _ := currentTime.Clock()

	if cur_hour == 8 && cur_min >= 0 && cur_min < 10 {
		dayOff, err := p.db.IsDayOffToday()
		if err != nil {
			p.myLog.Fatal("Call IsDayOffToday fails - ", err)
		}
		if dayOff {
			p.myLog.Info("Today is dayoff, isn't it? ")
		} else {
			p.myLog.Debug("Punching in")
			p.Attend(true)
		}
		time.Sleep(10 * time.Minute)
	}

	if cur_hour == 17 && cur_min >= 30 && cur_min < 40 {
		dayOff, err := p.db.IsDayOffToday()
		if err != nil {
			p.myLog.Fatal("Call IsDayOffToday fails - ", err)
		}
		if dayOff {
			p.myLog.Info("Today is dayoff, isn't it? ")
		} else {
			p.myLog.Debug("Punching Out")
			p.Leave(true)
		}
		time.Sleep(10 * time.Minute)
	}
	tick.executing = false
}

func (p *Puncher) Run() {
	tick := &Tick{
		ticker:    time.NewTicker(time.Minute),
		executing: false,
	}
	defer tick.ticker.Stop()
	for t := range tick.ticker.C {
		go p.run(tick, t)
	}
}

func (p *Puncher) request(request_type int) (*http.Response, error) {
	var reader *bytes.Reader
	var logMessage string
	switch request_type {
	case t_attend:
		reader = readerAttend
		logMessage = "Attend"
	case t_leave:
		reader = readerLeave
		logMessage = "Leave"
	case t_getToday:
		reader = readerGetToday
		logMessage = "Get today status"
	case t_getPast:
		reader = readerGetPast
		logMessage = "Get pass status"
	default:
		return nil, errors.New("request_type not fits")
	}

	p.myLog.Debug("Puncher is sending request to ", logMessage)

	// Construct the request
	req, err := http.NewRequest(http.MethodPost, "https://my.ntu.edu.tw/attend/ajax/signin.ashx", reader)
	if err != nil {
		p.myLog.Fatal("Parsing url fail")
	}
	set_request_header(req)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	res, err := p.puncher.Do(req)
	if err != nil {
		p.myLog.Error("Error : ", err)
		return nil, err
	}
	return res, nil
}
