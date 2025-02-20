# NTU-PuncherBot
The PuncherBot, written in Go, helps you to clock-in and clock out automatically. 

## Overview
PuncherBot:
- performs auto clock-in and clock-out
- provides a frendly CLI interface 
- integrates  with calendar database to store the dates to clock-in and clock-out
- integrates with Telegram for real-time notification 

## Get startted

### Prerequisites
Before you start, you need to 
- Install the go package. [Install link](https://go.dev/doc/install)
- Obtain telegram chatbot token and chat ID
    - Chatbot token : https://tcsky.cc/tips-01-telegram-chatbot/
    - Chat ID https://gist.github.com/nafiesl/4ad622f344cd1dc3bb1ecbe468ff9f8a
- Download Taiwanese goverment calendar from [Ruyut's calendar API](https://www.ruyut.com/2022/08/Taiwan-calendar-api.html) 
```
wget https://cdn.jsdelivr.net/gh/ruyut/TaiwanCalendar/data/{year}.json
```

### Steps 

1. Install this project 
```
git clone https://github.com/PoisWu/NTU-PuncherBot.git
```
2. Buid the source code into binary
```
go build 
```
3. Load the calendar into PuncherBot
```
./puncherbot load <year>.json
```

4. Edit the `config.toml` to your own credentials
```toml
[account]
username = "YOUR USERNAME"
password = "YOUR PASSWORD"

[telegram]
chat_id = "CHATID"
chatbot_token = "CHATBOTTOKEN"
```
5. Start the PuncherBot
```
./puncherbot run --config config.toml
```
