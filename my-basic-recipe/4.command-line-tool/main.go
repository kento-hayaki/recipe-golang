package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//** flagパッケージ
// 設定される変数のポインタを取得
var msg = flag.String("msg", "デフォルト値", "説明")
var n int

func init() {
	// ポインタを指定して設定を予約
	flag.IntVar(&n, "n", 1, "回数")
}

func main() {
	//** プログラム引数の取得 */ os.Args
	// プログラム引数が入った文字列型のスライス
	// 要素のひとつめはプログラム名
	fmt.Println(os.Args) // go run main.go hello > [/var/folders/_q/t_p_01392kzd_wdkjv6tm9dc0000gp/T/go-build805512576/b001/exe/main hello]

	//** フラグ（オプション）を便利に扱うパッケージ */ flagパッケージ
	// ここで実際に設定される
	flag.Parse()
	fmt.Println(strings.Repeat(*msg, n))
	//go run main.go > 'デフォルト値'
	//go run main.go -msg=こんにちは > こんにちは
	//go run main.go -msg=こんにちは -n=2 >  こんにちはこんにちは

	//flagパッケージとプログラム引数 >> flag.Args関数を用いる
	// os.Argsだとフラグも含まれる
	// flag.Args関数はフラグの分は除外される
	fmt.Println(flag.Args()) //go run main.go -msg=こんにちは -n=2 やぁ > [やぁ]

	// ** 標準入力と標準出力 */ osパッケージで提供されている*os.File型の変数
	// さまざまな関数やメソッドの引数として渡せる
	// エラーを出力する場合は標準エラー出力に出力する
	//fmt.Fprintln関数 ・・・出力先を指定して出力する
	// 末尾に改行をつけて表示する
	// os.Stdoutやos.Stderrに出力できる
	// ファイルにも出力できる
	fmt.Fprintln(os.Stderr, "エラー")   // 標準エラー出力に出力 //エラー
	fmt.Fprintln(os.Stderr, "Hello") // 標準出力に出力 //Hello

	// ** プログラムの終了 */ os.Exit(code int)
	// 終了コードを指定してプログラムを終了
	// プログラムの呼び出し元に終了状態を伝えられる
	// 0: 成功（デフォルト）
	fmt.Fprintln(os.Stderr, "エラー") //エラー
	//os.Exit(1) //プログラム終了//exit status 1
	// **プログラムの終了（エラー）
	// log.Fatal
	// 標準エラー出力（os.Stderr）にエラーメッセージを表示
	// os.Exit(1)で異常終了させる
	// 終了コードがコントロールできないためあまり多用しない
	// if err := f(); err != nil {
	// 	log.Fatal(err)
	// }

	// ** ファイルを扱う */ //osパッケージを用いる
	// 読み込み用にファイルを開く

	sf, err := os.Open("./a.txt")
	fmt.Println(sf) //&{0xc0000bc180}
	if err != nil {
		//もしsf, err := os.Open("./b.txt")として存在しないファイルを開こうとしたら
		fmt.Println(err) //open ./b.txt: no such file or directory
		return
		//return err
	}
	// 関数終了時に閉じる
	defer sf.Close()

	// 書き込み用にファイルを開く
	df, err := os.Create("./b.txt") //./b.txtがつくられる
	if err != nil {
		fmt.Println(err)
		return
		//return err
	}
	// 関数終了時に閉じる
	defer func() {
		if err := df.Close(); err != nil {
			return
			//rerr = err
		}
	}()

	// ** 関数の遅延実行 */ defer
	// 	関数終了時に実行される
	// 引数の評価はdefer呼び出し時
	// スタック形式で実行される（最後に呼び出したものが最初に実行）

	msg := "!!!"
	defer fmt.Println(msg)
	msg = "world"
	defer fmt.Println(msg)
	fmt.Println("hello")
	// hello world !!!

	// ⚠forの中でdeferは避ける
	// 	予約した関数呼び出しはreturn時に実行される
	// forの中を関数に分ければよい

	// ** 入出力関連 */
	//1行ずつ読み込む > bufio.Scannerを使用する
	// 標準入力から読み込む
	scanner := bufio.NewScanner(os.Stdin)
	// 1行ずつ読み込んで繰り返す
	// for scanner.Scan() {
	// 	if err := scanner.Err(); err != nil {
	// 		fmt.Fprintln(os.Stderr, err)
	// 	}
	//1行分を出力する
	// 	fmt.Println(scanner.Text()) // 入力した文字をstringでうけとる。この書き方だとずっと受け取りを待ち続けるからreturn
	// 	return
	// }
	// まとめてエラー処理をする
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "読み込みに失敗しました:", err)
		return
	}

	//** ファイルパスを扱う */ ・・path/filepathパッケージを使う
	//OSに寄らないファイルパスの処理が行える
	// パスを結合する
	path := filepath.Join("dir", "main.go")
	// 拡張子を取る
	fmt.Println(filepath.Ext(path)) //.go
	// ファイル名を取得
	fmt.Println(filepath.Base(path)) //main.go
	// ディレクトリ名を取得
	fmt.Println(filepath.Dir(path)) //dir

	// ** ディレクトリをウォークする */ ・・filepath.Walk関数を使う
	// Goファイルを探し出す
	err2 := filepath.Walk(".",
		func(path string, info os.FileInfo, err2 error) error {
			//fmt.Println(path,info,err2)//dir, <nil>, lstat dir: no such file or directory
			if filepath.Ext(path) == ".go" {
				fmt.Println(path) //main.go
			}
			return nil
		})
	if err2 != nil {
		fmt.Print(err2)
		return
	}

}
