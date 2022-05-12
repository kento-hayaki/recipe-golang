package main

import (
	"fmt"
	"time"

	greeting "github.com/tenntenn/greeting"
	greeting2 "github.com/tenntenn/greeting/v2"
)

func main() {
	fmt.Println(greeting.Do())
	fmt.Println(greeting2.Do(time.Now()))
	// ...

	// ** スコープ */
	// 	識別子（変数名、関数名など）を参照できる範囲
	// 参照元によって所属するスコープが違う
	// 親子関係があり親のスコープの識別子は参照できる

	// ** パッケージの初期化 */
	// 	依存パッケージの初期化
	// importしているパッケージリストを出す
	// 依存関係を解決して依存されてないパッケージから初期化していく
	// 各パッケージの初期化
	// パッケージ変数の初期化する
	// **init関数の実行を行う



}


