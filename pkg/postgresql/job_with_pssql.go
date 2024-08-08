package postgresql

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

func InsertDb(date string, summa int) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	urlToDataBase := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", Cfg.PGuser, Cfg.PGpassword, Cfg.PGaddress, Cfg.PGPort, Cfg.PGdbname)
	conn, err := pgx.Connect(context.Background(), urlToDataBase)
	if err != nil {
		log.Println(err.Error())
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), `INSERT INTO user_of_site (date_of_purchase, summa_of_buy) VALUES ($1, $2)`, date, summa)
	if err != nil {
		log.Println(err.Error())
	}

}

func GetDateFromDb(date string) int {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	urlToDataBase := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", Cfg.PGuser, Cfg.PGpassword, Cfg.PGaddress, Cfg.PGPort, Cfg.PGdbname)
	conn, err := pgx.Connect(context.Background(), urlToDataBase)
	if err != nil {
		log.Println(err.Error())
	}
	defer conn.Close(context.Background())
	//rows, _ := conn.Query(context.Background(), `SELECT summa_of_buy FROM user_of_site WHERE date_of_purchase = $1`, date)
	////get from date data
	var rows pgx.Rows        //заглушка
	var error error          //заглушка
	var info_from_db int = 0 //return value
	if date != "" {          // если чел не сделал get для суммы
		parts := strings.Split(date, "-") // получаю данные из get запроса
		year, _ := strconv.Atoi(parts[0])
		month, _ := strconv.Atoi(parts[1])
		day, _ := strconv.Atoi(parts[2])
		_ = day //заглушка пока
		fmt.Printf("DATE: %v, YEAR:%v, MONTH:%v, Day: %v\n", date, year, month, strconv.Itoa(day))
		//условие по сути на день все равно, выдает кол-во потраченных денег за выбранный месяц с годом
		rows, _ = conn.Query(context.Background(), `SELECT summa_of_buy FROM user_of_site 
		WHERE EXTRACT(YEAR FROM date_of_purchase) = $1 AND EXTRACT(MONTH FROM date_of_purchase) = $2`, strconv.Itoa(year), strconv.Itoa(month))
		_ = error

		//get data from db
		for rows.Next() {
			var temp int
			rows.Scan(&temp)
			info_from_db += temp
		}
		defer rows.Close() // страно но до этого момента rows = nil
		return info_from_db
	}
	return info_from_db
}

func init() {
	file, err := os.Open("config.cfg")
	if err != nil {
		fmt.Println("Error open .cfg", err)
		panic("Can't open the file \"setting.cfg\"")
	}
	defer file.Close()

	fileInfo, _ := file.Stat()                   // получаю стату файла для его размера
	readSetting := make([]byte, fileInfo.Size()) // делаю такого же размера переменную
	_, err = file.Read(readSetting)
	if err != nil {
		panic("can't read file")
	}
	// fmt.Println(string(readSetting))  работает

	err = json.Unmarshal(readSetting, &Cfg) //unmarshal и json в обьект marshal из object in json
	if err != nil {
		panic("json err")
	}
}

type setting struct { // должен повторять структуру json
	PGaddress  string
	PGpassword string
	PGuser     string
	PGdbname   string
	PGPort     string
}

var (
	Cfg setting // for use in main for open db
)
