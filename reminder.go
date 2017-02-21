package main

import ( 
	"fmt"

	"github.com/robfig/cron"
    "github.com/tucnak/telebot"
    _ "github.com/mattn/go-sqlite3"
)

var msg telebot.Message

func reminderRoutine(msg_ telebot.Message) {
	c := cron.New(); msg = msg_ 
	c.AddFunc("0 20 * * *", sendReminder) //@every 1m
	c.Start() //run at 8pm
}

func sendReminder() {
	bot.SendMessage(msg.Chat, "Ciao! Ti sei ricordato di segnare le tue spese di oggi?!", nil)
	fmt.Printf("Reminder sent to %d\n", msg.Chat.ID)
}

func reloadUsers() {
	rows, err := db.Query("SELECT telegramID FROM users WHERE reminderActive=1"); checkErr(err)
	var telegramID int64; var msg telebot.Message;

	for rows.Next() {
	    err = rows.Scan(&telegramID); checkErr(err)
	    msg.Chat.ID = telegramID; go reminderRoutine(msg)
	}
}