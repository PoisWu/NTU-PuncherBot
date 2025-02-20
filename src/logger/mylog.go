package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
}

type MyLogger struct {
	d     telegram_dealer
	debug *log.Logger
	info  *log.Logger
	warn  *log.Logger
	error *log.Logger
	fatal *log.Logger
}

func NewMyLogger(chat_ID, chatbot_token string) *MyLogger {
	return &MyLogger{
		d:     *NewTelegramDealer(chat_ID, chatbot_token),
		debug: log.New(os.Stdout, "DEBUG: \t", log.LstdFlags),
		info:  log.New(os.Stdout, "INFO: \t", log.LstdFlags),
		warn:  log.New(os.Stdout, "WARN: \t", log.LstdFlags),
		error: log.New(os.Stderr, "ERROR: \t", log.LstdFlags),
		fatal: log.New(os.Stderr, "FATAL: \t", log.LstdFlags),
	}
}

func (l *MyLogger) Debug(v ...interface{}) {
	l.debug.Println(v...)
}

func (l *MyLogger) Info(v ...interface{}) {
	resp, _ := l.d.info(fmt.Sprintln(v...))
	resBody, _ := io.ReadAll(resp.Body)
	l.Debug(string(resBody))
	l.info.Println(v...)
}

func (l *MyLogger) Warn(v ...interface{}) {
	resp, _ := l.d.warn(fmt.Sprintln(v...))
	resBody, _ := io.ReadAll(resp.Body)
	l.Debug(string(resBody))
	l.warn.Println(v...)
}

func (l *MyLogger) Error(v ...interface{}) {
	resp, _ := l.d.error(fmt.Sprintln(v...))
	resBody, _ := io.ReadAll(resp.Body)
	l.Debug(string(resBody))
	l.error.Println(v...)
}

func (l *MyLogger) Fatal(v ...interface{}) {
	resp, _ := l.d.error(fmt.Sprintln(v...))
	resBody, _ := io.ReadAll(resp.Body)
	l.Debug(string(resBody))
	l.fatal.Println(v...)
}

type telegram_dealer struct {
	chat_id                string
	telegram_chatbot_token string
}

func NewTelegramDealer(chat_ID, chatbot_token string) *telegram_dealer {
	return &telegram_dealer{
		chat_id:                chat_ID,
		telegram_chatbot_token: chatbot_token,
	}
}

func (d *telegram_dealer) send(payload []byte) (*http.Response, error) {
	telegrambot_api := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", d.telegram_chatbot_token)
	resp, err := http.Post(
		telegrambot_api,
		"application/json",
		bytes.NewBuffer(payload),
	)
	return resp, err
}

func (d *telegram_dealer) info(message string) (*http.Response, error) {
	body, _ := json.Marshal(map[string]string{
		"chat_id":    d.chat_id,
		"text":       fmt.Sprintf("*[INFO!]\n*%s", message),
		"parse_mode": "Markdown",
	})
	return d.send(body)
}

func (d *telegram_dealer) warn(message string) (*http.Response, error) {
	body, _ := json.Marshal(map[string]string{
		"chat_id":    d.chat_id,
		"text":       fmt.Sprintf("*[❗Warn❗]*\n%s", message),
		"parse_mode": "Markdown",
	})
	return d.send(body)
}

func (d *telegram_dealer) error(message string) (*http.Response, error) {
	body, _ := json.Marshal(map[string]string{
		"chat_id":    d.chat_id,
		"text":       fmt.Sprintf("*[❗❗ERROR❗❗]*\n%s", message),
		"parse_mode": "Markdown",
	})
	return d.send(body)
}
