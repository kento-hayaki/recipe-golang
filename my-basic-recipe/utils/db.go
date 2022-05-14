package utils

import (
	"database/sql"
	"log"
)

// db接続周り
func connectDB() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/go_typing")
	if err != nil {
		panic(err)
	}
	defer db.Close()
}

func selectSQL(db *sql.DB) []Result {
	rows, err := db.Query("select * from results;")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var results []Result
	for rows.Next() {
		result := Result{}
		err := rows.Scan(&result.wpm, &result.miss, &result.playTimeMinute, &result.playDate)
		if err != nil {
			panic(err)
		}
		results = append(results, result)
	}

	return results
}

func insertSQLOnResults(db *sql.DB, sql []string) {
	ins, err := db.Prepare("INSERT INTO results(wpm,miss,playTimeMinute,playDate) VALUES(?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	ins.Exec("") //TODO引数で受けとったもの
}
