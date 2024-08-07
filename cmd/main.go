// TODO
// Для личных финансов день скоко потратил 3 категории еда, хуйня, важное
// хранит в базе дни а потом ключ считывает в 1 месяц в другой таблице сумма по дням
package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"mymod/pkg/postgresql"
)

// TODO clear main
// dependency injection
type application struct {
	port     string
	infoLog  *log.Logger
	errorLog *log.Logger
}

func (app application) mainPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		summa := 0
		date := r.FormValue("date")
		checkSumma, err := strconv.Atoi(r.FormValue("summa"))
		if err != nil {
			http.NotFound(w, r)
			return
		} else {
			summa = int(checkSumma)
		}

		postgresql.InsertDb(date, summa)

	}

	var date_summa int = 0
	if r.Method == http.MethodGet {
		date_out := r.FormValue("date-output")
		date_summa = postgresql.GetDateFromDb(date_out)
	}

	//есть еще tmpl := template.Must(template.Parsefile(...)) он без проверки + panic
	tmp, err := template.ParseFiles("./ui/html/index.html")
	if err != nil {
		app.errorLog.Println(err.Error())
	}

	//for if user use custom port
	setting_html := struct {
		Port            string
		Summa_from_date int
	}{
		Port:            app.port,
		Summa_from_date: date_summa,
	}

	err = tmp.Execute(w, setting_html)
	if err != nil {
		app.errorLog.Println(err.Error())
	}
}

func (app application) easterEgg(w http.ResponseWriter, r *http.Request) {
	// для гиф с тянкой оленем спратать на лого например
	tmp, err := template.ParseFiles("./ui/html/easterEgg.html")
	if err != nil {
		app.errorLog.Println(err.Error())
	}
	//for if user use custom port
	myport := struct {
		Port string
	}{
		Port: app.port,
	}

	err = tmp.Execute(w, myport)
	if err != nil {
		app.errorLog.Println(err.Error())
	}
}

func main() {
	//flag -addr=
	addr := flag.String("addr", ":8080", "Port for server")
	flag.Parse()

	file, _ := os.OpenFile("./infolog.log", os.O_APPEND, 0600)
	defer file.Close()
	infolog := log.New(file, "INFO:\t", log.LstdFlags|log.Lshortfile)
	errorlog := log.New(file, "ERROR:\t", log.LstdFlags|log.Lshortfile)
	app := &application{
		port:     *addr,
		infoLog:  infolog,
		errorLog: errorlog,
	}

	//fileserver
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", app.mainPage)
	mux.HandleFunc("/easterEgg/", app.easterEgg)

	log.Printf("server launch on localhost%v", *addr) // for localhost
	app.infoLog.Printf("server launch on localhost%v", *addr)
	//launch server
	err := http.ListenAndServe(*addr, mux)
	if err != nil {
		app.errorLog.Println(err.Error())
	}
}
