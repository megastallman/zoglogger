package main
import (
    "fmt"
    "regexp"
    "net/http"
    _ "github.com/mattn/go-sqlite3"
    "database/sql"
    "flag"
    "hash/fnv"
    "strconv"
    "strings"
    "github.com/goji/httpauth"
    "io/ioutil"
)

var dbPath string
var dbName string

func handler(w http.ResponseWriter, r *http.Request) {
	var tabline string

	// Beginning
	fmt.Fprintf(w, "<html>\n<head>\n<title>ZOG is logging you!</title>\n</head>\n<body link='#0000FF' vlink='#CC00FF' bgcolor ='#CC00FF'><big><p><b>Welcome to ZOG</b></p></big>\n%s", r.URL.Path[1:])

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

		// Normalize < sign
		serverName = strings.Replace(serverName, "<", "&lt", -1)
		roomName = strings.Replace(roomName, "<", "&lt", -1)
		userName = strings.Replace(userName, "<", "&lt", -1)
		currentTime = strings.Replace(currentTime, "<", "&lt", -1)
		ircMessage = strings.Replace(ircMessage, "<", "&lt", -1)
		// Normalize > sign
		serverName = strings.Replace(serverName, ">", "&gt", -1)
		roomName = strings.Replace(roomName, ">", "&gt", -1)
		userName = strings.Replace(userName, ">", "&gt", -1)
		currentTime = strings.Replace(currentTime, ">", "&gt", -1)
		ircMessage = strings.Replace(ircMessage, ">", "&gt", -1)

		// Highlight links
		ircMessage = highlight_links(ircMessage)

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

func highlight_links(ircMessage string) string {

	var link []string
	var extracted_link string
	var formatted_link string

	re := regexp.MustCompile("( |^)http(s?)://[^ ]+")
	for {
		link = re.FindStringSubmatch(ircMessage)
		if link != nil {
			extracted_link = strings.TrimSpace(string(link[0]))
			//fmt.Println(extracted_link)
			formatted_link = "<a href=" + extracted_link + ">" + extracted_link + "</a> "
			ircMessage = strings.Replace(ircMessage, link[0], formatted_link, -1)
		} else {
			break
		}
	}
	return ircMessage

}

func main() {

	var listenSocket string
	var webUser string
	var webAss string
	var webAssFile string
	var webCrtFile string
	var webKeyFile string
	flag.StringVar(&listenSocket, "listenSocket", ":9991", "a string var")
	flag.StringVar(&dbPath, "dbPath", "db.sqlite3", "a string var")
	flag.StringVar(&dbName, "dbName", "zogger", "a string var")
	flag.StringVar(&webUser, "webUser", "", "a string var")
	flag.StringVar(&webAss, "webAss", "", "a string var")
	flag.StringVar(&webAssFile, "webAssFile", "", "a string var")
	flag.StringVar(&webCrtFile, "webCrtFile", "", "a string var")
	flag.StringVar(&webKeyFile, "webKeyFile", "", "a string var")
	
	// Parse the flags
	flag.Parse()

	// Using HTTP-auth if needed
	if webUser != "" {
		// Read assword from file
		if webAssFile != "" {
			fmt.Printf("Reading room assword from: " + webAssFile + "\n")
			dat, err := ioutil.ReadFile(webAssFile)
			checkErr(err)
			webAss = strings.TrimSpace(string(dat))
		}
		http.Handle("/", httpauth.SimpleBasicAuth(webUser, webAss)(http.HandlerFunc(handler)))
	} else {
		http.HandleFunc("/", handler)
	}

	// Verify if HTTPS should be used
	if webCrtFile == "" {
		http.ListenAndServe(listenSocket, nil)
        } else {
		http.ListenAndServeTLS(listenSocket, webCrtFile, webKeyFile, nil)
	}
}

func checkErr(err error) {
        if err != nil {
	        panic(err)
        }
}
