// ** データベース
// **database/sqlパッケージ
// RDBにアクセスするためのパッケージ
// 共通機能を提供
// クエリの発行
// トランザクション
// データベースの種類ごとにドライバが存在
// http://golang.org/s/sqldrivers

// **ドライバの登録
// - ドライバ
// 各種RDBに対応したドライバ
// MySQLやSQLiteなど
// - インポートするだけで登録される
// initで登録される
// パッケージ自体は直接使わない

// 例：SQLiteのドライバの登録（LinuxとWindowsの場合）
// import _ "modernc.org/sqlite"

// ** SQLite
// ファイルベースのデータベース
// 軽量なRDBで単一のファイルに記録される
// アプリケーションに組み込まれることが多い
// 他のRDBは通常はサーバとして動作する
// モバイルアプリでも使用されることが多い

package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/tenntenn/sqlite"
)

func main() {
	// ** データベースのオープン
	// - *sql.DBの特徴
	// 複数のゴルーチンから使用可能
	// コネクションプール機能
	// 一度開いたら使いまわす
	// Closeは滅多にしない
	// - Open関数を使用する
	// ドライバ名, 接続文字列
	db, err := sql.Open("sqlite", "database.db")
	if err != nil {
		fmt.Println("dot connect")
	}
	fmt.Println(db)

	// ** テーブルの作成*/
	// (*sql.DB).Execを使う
	const sql = `
		CREATE TABLE IF NOT EXISTS user (
			id   INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			age  INTEGER NOT NULL
		);
	`
	if _, err := db.Exec(sql); err != nil {
		// エラー処理
	}

	// ** レコードの挿入
	//AUTOINCREMENTのIDは*sql.Resultから取得できる
	type User struct {
		ID   int64
		Name string
		Age  int64
	}
	users := []*User{{Name: "tenntenn", Age: 32}, {Name: "Gopher", Age: 10}}
	for i := range users {
		const sql = "INSERT INTO user(name, age) values (?,?)"
		r, err := db.Exec(sql, users[i].Name, users[i].Age)
		if err != nil { /* エラー処理 */
		}
		id, err := r.LastInsertId()
		if err != nil { /* エラー処理 */
		}
		users[i].ID = id
		fmt.Println("INSERT", users[i])
	}

	// ** 複数レコードのスキャン */ ・・(*sql.DB).Queryと*sql.Rowsを使う
	rows, err := db.Query("SELECT * FROM user WHERE age = ?", 24)
	if err != nil { /* エラー処理 */
	}
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Age); err != nil {
			/* エラー処理 */
		}
		fmt.Println(u)
	}
	if err := rows.Err(); err != nil { /* エラー処理 */
	}

	// **レコードの更新
	// 更新したレコード数は*sql.Resultから取得
	r, err := db.Exec("UPDATE user SET age = age + 1 WHERE id = 1")
	if err != nil { /* エラー処理 */
	}
	cnt, err := r.RowsAffected()
	if err != nil { /* エラー処理 */
	}
	fmt.Println("Affected rows:", cnt)

	// ** Q. 電話帳を作ろう */
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}

	// ** 11.2. トランザクション
	// ** トランザクションの開始
	// ** トランザクションを使った例
	tx, err := db.Begin()
	if err != nil { /* エラー処理 */
	}
	row := tx.QueryRow("SELECT * FROM user WHERE id = 1")
	var u User
	if err := row.Scan(&u.ID, &u.Name, &u.Age); err != nil {
		tx.Rollback()
		/* エラー処理 */
	}
	const updateSQL = "UPDATE user SET age = ? WHERE id = 1"
	if _, err = tx.Exec(updateSQL, u.Age+1, u.ID); err != nil {
		tx.Rollback()
		/* エラー処理 */
	}
	if err := tx.Commit(); err != nil { /* エラー処理 */
	}

}

// ** SQLの実行 */ ・・*sql.DBのメソッドを使用
// INSERTやDELETEなど
//context.Contextを引数に取るバージョンもある
// func (db *sql.DB) Exec(query string, args ...interface{}) (Result, error) {}

// SELECTなどで複数レコードを取得する場合
// func (db *sql.DB) Query(query string, args ...interface{}) (*Rows, error)

// SELECTなどで1つのレコードを取得する場合
// func (db *sql.DB) QueryRow(query string, args ...interface{}) *Row

// ** 11.2. トランザクション
// ** トランザクションの開始
// - (*sql.DB).Beginを呼ぶ
// トランザクションを開始する
// func (db *sql.DB) Begin() (*Tx, error) {}
// - context.Contextを渡したい場合
// Contextを渡してトランザクションを開始する
// func (db *DB) BeginTx(context.Context, *TxOptions) (*Tx, error)

// ** トランザクションに対する処理
// *sql.Txのメソッドを使用
// INSERTやDELETEなど
// func (tx *Tx) Exec(query string, args ...interface{}) (Result, error)
// SELECTなどで複数レコードを取得する場合
// func (tx *Tx) Query(query string, args ...interface{}) (*Rows, error)
// SELECTなどで1つのレコードを取得する場合
// func (tx *Tx) QueryRow(query string, args ...interface{}) *Row
// コミット
// func (tx *Tx) Commit() error
// ロールバック
// func (tx *Tx) Rollback() error

// ** Q. 電話帳を作ろう */
type Record struct {
	ID    int64
	Name  string
	Phone string
}

func run() error {
	db, err := sql.Open(sqlite.DriverName, "addressbook.db")
	if err != nil {
		return err
	}

	if err := createTable(db); err != nil {
		return err
	}

	for {
		if err := showRecords(db); err != nil {
			return err
		}

		if err := inputRecord(db); err != nil {
			return err
		}
	}

	return nil
}

// テーブルの作成
func createTable(db *sql.DB) error {
	const sql = `
	CREATE TABLE IF NOT EXISTS addressbook (
			id    INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			name  TEXT NOT NULL,
			phone TEXT NOT NULL
	);`

	if _, err := db.Exec(sql); err != nil {
		return err
	}

	return nil
}

// テーブルの中身全件取得
func showRecords(db *sql.DB) error {
	fmt.Println("全件表示")
	rows, err := db.Query("SELECT * FROM addressbook")
	if err != nil {
		return err
	}
	for rows.Next() {
		var r Record
		if err := rows.Scan(&r.ID, &r.Name, &r.Phone); err != nil {
			return err
		}
		fmt.Printf("[%d] Name:%s TEL:%s\n", r.ID, r.Name, r.Phone)
	}
	fmt.Println("--------")

	return nil
}

// テーブルにデータ挿入
func inputRecord(db *sql.DB) error {
	var r Record

	fmt.Print("Name >")
	fmt.Scan(&r.Name)

	fmt.Print("TEL >")
	fmt.Scan(&r.Phone)

	const sql = "INSERT INTO addressbook(name, phone) values (?,?)"
	_, err := db.Exec(sql, r.Name, r.Phone)
	if err != nil {
		return err
	}

	return nil
}

// ** Q. 電話帳を作ろう ここまで */
