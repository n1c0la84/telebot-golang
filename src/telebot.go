package main

import (
    "database/sql"; "strings"; "time" 

    "github.com/tucnak/telebot"
    _ "github.com/mattn/go-sqlite3"
)

var db *sql.DB; var bot *telebot.Bot

func main() {
	//get connection to telegram and sqlite3 
    bot_, err := telebot.NewBot("325327338:AAGBoD62nd6y1T3ZBcZcZXEcJPkN2uftQB4"); checkErr(err)
    db_, err := sql.Open("sqlite3", "./foo.db"); checkErr(err); db = db_; bot = bot_ 

    //create message channel
    messages := make(chan telebot.Message, 100) 
    bot.Listen(messages, 1*time.Second) 
    defer close(messages)

    reloadUsers() //for reminders

    for message := range messages {
    	//skip messages that are not text
        if message.Text == "" { continue } 

        //command is the message received
        command := strings.ToLower(message.Text)

        switch command {
        case "/start": 
            startCmd(message) 

        case "/help", "aiuto", "ciao":
            helpCmd(message) //WORKING!
            
        case "cancella":
            cancellaCmd(message) //WORKING!
            
        case "lista":
            listaCmd(message) //WORKING! 
            
        case "totale":
            totaleCmd(message) //WORKING!
            
        case "reset":
            resetCmd(message) //WORKING!

        case "noreminder":
            noreminderCmd(message) 

        case "reminderon":
            noreminderCmd(message) 
            
        default:
        	if strings.HasPrefix(command, "spesa") { 
        		spesaCmd(message) //WORKING!
        	} else {
            	bot.SendMessage(message.Chat, "Comando sconosciuto!", nil)
            }
        }
    }
}

func formatDate(datetime string) string {
    datepart := strings.Split(datetime, " ")[0]
    day := strings.Split(datepart, "-")[2];  month := strings.Split(datepart, "-")[1]

    timepart := strings.Split(datetime, " ")[1]
    hour := strings.Split(timepart, ":")[0]; minutes := strings.Split(timepart, ":")[1]
    
    return day + "/" + month + ", " + hour + ":" + minutes
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}