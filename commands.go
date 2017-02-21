package main

import (
    "fmt"; "strings"; "time" 

    "github.com/tucnak/telebot"
    _ "github.com/mattn/go-sqlite3"
)

func startCmd(msg telebot.Message) {
	go func(msg telebot.Message) {
		rows, err := db.Query("SELECT * FROM users WHERE telegramID=?", msg.Chat.ID); checkErr(err)
	    if rows.Next() { //already present
	    	helpCmd(msg)
	    	return
	    } //not present 
		
		stmt, err := db.Prepare("INSERT INTO users(telegramID, reminderActive) values(?, 1)"); checkErr(err)
	    _, err = stmt.Exec(msg.Chat.ID); checkErr(err) 
	    go reminderRoutine(msg)

		helpCmd(msg); rows.Close() 
	}(msg)
}

func helpCmd(msg telebot.Message) {
	go func(msg telebot.Message) {
	    bot.SendMessage(msg.Chat, fmt.Sprintf("Benvenuto, %s! I comandi sono:", msg.Sender.FirstName), nil)
	    bot.SendMessage(msg.Chat, "SPESA <amount> <motivo>, per registrare una spesa", nil) 
	    bot.SendMessage(msg.Chat, "CANCELLA, per cancellare l'ultima spesa", nil)
	    bot.SendMessage(msg.Chat, "LISTA, per avere l'elenco di tutte le spese", nil) 
	    bot.SendMessage(msg.Chat, "TOTALE, per avere il totale speso", nil)
	}(msg)
}

func spesaCmd(msg telebot.Message) {
	go func(msg telebot.Message) {
	    s := strings.Split(msg.Text, " "); 
	    var amount, reason string

	    if len(s) == 1 {
	    	bot.SendMessage(msg.Chat, "Inserisci almeno l'amount!", nil)
	    	return 
	    } else if len(s) == 2 {
	    	amount, reason = strings.Replace(s[1], "€", "", -1), ""
	    } else {
	    	amount, reason = strings.Replace(s[1], "€", "", -1), s[2]
	    }

	    stmt, err := db.Prepare("INSERT INTO spese(telegramID, amount, reason, created) values(?, ?, ?, ?)"); checkErr(err)
	    _, err = stmt.Exec(msg.Chat.ID, amount, reason, time.Now()); checkErr(err) 

	    bot.SendMessage(msg.Chat, "Spesa registrata!", nil)
	}(msg)
}

func cancellaCmd(msg telebot.Message) {
	go func(msg telebot.Message) {
	    rows, err := db.Query("SELECT MAX(created) FROM spese WHERE telegramID=?", msg.Chat.ID); checkErr(err)
	    rows.Next(); var maxCreated string; err = rows.Scan(&maxCreated); checkErr(err); rows.Close()
	    stmt, err := db.Prepare("DELETE FROM spese WHERE telegramID=? and created=?"); checkErr(err)
	    _, err = stmt.Exec(msg.Chat.ID, maxCreated); checkErr(err)

	    bot.SendMessage(msg.Chat, "Ultima spesa cancellata!", nil)
	}(msg)
}

func listaCmd(msg telebot.Message) {
	go func(msg telebot.Message) {
	    rows, err := db.Query("SELECT created, amount, reason  FROM spese WHERE telegramID=?", msg.Chat.ID); checkErr(err)
	    var created string; var amount float64; var reason string

	    var empty bool = true
	    for rows.Next() {
	        err = rows.Scan(&created, &amount, &reason); checkErr(err); empty = false
	        bot.SendMessage(msg.Chat, fmt.Sprintf("%s: %.2f€  %s", formatDate(created), amount, reason), nil)
	    }
	    if empty { 
	        bot.SendMessage(msg.Chat, "Nessuna spesa registrata!", nil)
	    }
	    rows.Close() 
	}(msg)
}

func totaleCmd(msg telebot.Message) {
	go func(msg telebot.Message) {
	    rows, err := db.Query("SELECT SUM(amount) FROM spese WHERE telegramID=?", msg.Chat.ID); checkErr(err)
	    
	    defer func() {
	        if e := recover(); e != nil {
	            bot.SendMessage(msg.Chat, "Totale speso: 0€", nil)
	        }
		}() //exception handling

	    var totale float64
	    for rows.Next() {
	        err = rows.Scan(&totale); checkErr(err)
	        bot.SendMessage(msg.Chat, fmt.Sprintf("Totale speso: %.2f€", totale), nil)
	    }
	    rows.Close() 
	}(msg)
}

func resetCmd(msg telebot.Message) {
	go func(msg telebot.Message) {
	    stmt, err := db.Prepare("DELETE FROM spese WHERE telegramID=?"); checkErr(err)
	    _, err = stmt.Exec(msg.Chat.ID); checkErr(err)

	    bot.SendMessage(msg.Chat, "Spese resettate!", nil)
	}(msg)
}

func noreminderCmd(msg telebot.Message) {
	go func(msg telebot.Message) {
		stmt, err := db.Prepare("update users set reminderActive=0 where telegramID=?"); checkErr(err)
		_, err = stmt.Exec(msg.Chat.ID); checkErr(err)
	}(msg)
}

func reminderonCmd(msg telebot.Message) {
	go func(msg telebot.Message) {
		stmt, err := db.Prepare("update users set reminderActive=1 where telegramID=?"); checkErr(err)
		_, err = stmt.Exec(msg.Chat.ID); checkErr(err)
	}(msg)
}