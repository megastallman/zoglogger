package main
import (
    "fmt"
    "net/http"
    _ "github.com/mattn/go-sqlite3"
    "database/sql"
    "flag"
    "hash/fnv"
    "strconv"
    "strings"
)

var dbPath string
var dbName string

func handler(w http.ResponseWriter, r *http.Request) {
	var tabline string

	// Beginning
	fmt.Fprintf(w, "<html>\n<head>\n<title>ZOG is logging you!</title>\n</head>\n<body bgcolor ='#CC00FF'><big><p><b>Welcome to ZOG</b></p></big>\n%s", r.URL.Path[1:])

	// Open DB
	db, err := sql.Open("sqlite3", dbPath)
	checkErr(err)

	// Query DB
	//rows, err := db.Query("SELECT * FROM zogger")
	rows, err := db.Query("SELECT * FROM zogger ORDER BY currentTime DESC")
	checkErr(err)

	// Output to web
	for rows.Next() {
		var serverName string
		var roomName string
		var userName string
		var currentTime string
		var ircMessage string
		err = rows.Scan(&serverName, &roomName, &userName, &currentTime, &ircMessage)
		checkErr(err)

		// Normalize percent signs
		serverName = strings.Replace(serverName, "%", "%%", -1)
		roomName = strings.Replace(roomName, "%", "%%", -1)
		userName = strings.Replace(userName, "%", "%%", -1)
		currentTime = strings.Replace(currentTime, "%", "%%", -1)
		ircMessage = strings.Replace(ircMessage, "%", "%%", -1)

		// Getting userName's hash
		h := fnv.New32a()
		h.Write([]byte(userName))
		bgHash := h.Sum32()
		strBgHash := strconv.FormatUint(uint64(bgHash), 16)
		bgColor := strBgHash[len(strBgHash)-6:]
		//fmt.Printf("|" + bgColor + "|")

		tabline = "<br><span style='background-color: #" + bgColor + "'><b>" + userName + ":</b></span>" + "<span style='background-color: #99FFFF'>" + ircMessage + "</span>" + "<small><i>: at " + currentTime + "</i></small>%s"
		fmt.Fprintf(w, tabline, r.URL.Path[1:])
	}

	// Close DB
	db.Close()
	checkErr(err)

	// Ending
	fmt.Fprintf(w, "</body></html>%s", r.URL.Path[1:])

	}

func main() {

	var listenSocket string
	flag.StringVar(&listenSocket, "listenSocket", ":9991", "a string var")
	flag.StringVar(&dbPath, "dbPath", "db.sqlite3", "a string var")
	flag.StringVar(&dbName, "dbName", "zogger", "a string var")
	
	// Parse the flags
	flag.Parse()

	http.HandleFunc("/", handler)
	http.ListenAndServe(listenSocket, nil)
}

func checkErr(err error) {
    if err != nil {
	panic(err)
    }
}
