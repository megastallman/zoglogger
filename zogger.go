package main

import (
    "github.com/thoj/go-ircevent"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "flag"
    "fmt"
    "time"
    "io/ioutil"
)

func main() {

	var ircMessage = ""
	var currentTime = ""
	var ircUser = ""

	var serverName string
	flag.StringVar(&serverName, "serverName", "irc.freenode.net:6667", "a string var")
	var roomName string
	flag.StringVar(&roomName, "roomName", "#my-bot-test", "a string var")
	var userName string
	flag.StringVar(&userName, "userName", "zogbot", "a string var")
	var helloMsg string
	flag.StringVar(&helloMsg, "helloMsg", "ZOG is logging you", "a string var")
	var dbPath string
	flag.StringVar(&dbPath, "dbPath", "db.sqlite3", "a string var")
	var dbName string
	flag.StringVar(&dbName, "dbName", "zogger", "a string var")
	var roomAssword string
	flag.StringVar(&roomAssword, "roomAssword", "", "a string var")
	var assFile string
	flag.StringVar(&assFile, "assFile", "", "a string var")

	// Parse the flags
	flag.Parse()

	// Read assword from file
	if assFile != "" {
		fmt.Printf("Reading room assword from: " + assFile + "\n")
		dat, err := ioutil.ReadFile(assFile)
		checkErr(err)
		roomAssword = string(dat)
		//fmt.Print(roomAssword)
	} else {
		fmt.Printf("No assFile provided\n")
	}

	con := irc.IRC(userName, userName)
	err := con.Connect(serverName)
	checkErr(err)

	db, err := sql.Open("sqlite3", dbPath)
	checkErr(err)

	con.AddCallback("001", func (e *irc.Event) {
        	con.Join(roomName + " " + roomAssword)
	})

	con.AddCallback("JOIN", func (e *irc.Event) {
        	con.Privmsg(roomName, helloMsg)
	})
	con.AddCallback("PRIVMSG", func (e *irc.Event) {
		ircMessage = e.Message()
		ircUser = e.Nick
		currentTime = time.Now().String()
		// Insert to DB
		stmt, err := db.Prepare("INSERT INTO " + dbName + "(serverName, roomName, ircUser, currentTime, ircMessage) values(?,?,?,?,?)")
		checkErr(err)
		res, err := stmt.Exec(serverName, roomName, ircUser, currentTime, ircMessage)
		checkErr(err)
		id, err := res.LastInsertId()
		checkErr(err)
		fmt.Println(id)
	})

	con.Loop()
}


func checkErr(err error) {
    if err != nil {
            panic(err)
    }
}
