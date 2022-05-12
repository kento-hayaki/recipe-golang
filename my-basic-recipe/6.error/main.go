package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"
)

func main() {
	// ** エラー */ ・・errorインタフェース
	// エラーを表す型
	// 最後の戻り値として返すことが多い
	type error interface {
		Error() string
	}

	// ** エラー処理*/・・error型で表現する
	// エラーがない場合はnilになる
	// 実行時エラーを表す
	// ** nilと比較してエラーが発生したかをチェックする
	// スコープを狭くするため、代入付きのifを用いる場合が多い
	//if err := f(); err != nil { //**こんなかんじ
	// エラー処理
	//}

	//** エラー処理のよくあるミス・・err変数を使い回すことによるハンドルミス
	// エラーが発生してもハンドルされずに次に進んでしまう
	// errcheckなどの静的解析ツールで回避できる
	_, err := os.Open("file.txt")
	if err != nil {
		// エラー処理
		fmt.Println("error") //file.txがないとき出力される
	}
	// ** 本来は err = doSomething(f) としたつもり.下記エラーは絶対に実行されない
	// doSomething(f)
	// if err != nil {
	// 	// エラー処理
	// }

	//** 文字列ベースで簡単なエラーの作り方
	//1. errors.Newを使う
	err = errors.New("Error")
	fmt.Println(err) // Error
	//2. fmt.Errorfを使う・・書式を指定してエラーを作る
	name := "error test"
	err = fmt.Errorf("%s is not found", name) //error test is not found
	fmt.Println(err)

	// Stringerインタフェースに変換する関数を実装し実行
	//s := S("test")
	s := 100
	if s1, err := ToStringer(s); err != nil {
		//s = 100のとき
		fmt.Fprintln(os.Stderr, "ERROR:", err) //ERROR: CastError
	} else {
		//s = S("test")のとき
		fmt.Println("s = ", s1.String()) //s1 =  test
	}

	// ** エラー処理をまとめる */ ・・bufio.Scannerの実装が参考になる
	// 途中でエラーが発生したらそれ以降の処理を飛ばす
	// すべての処理が終わったらまとめてエラーを処理
	// それ以降の処理を実行する必要ない場合に使う
	// エラー処理が1箇所になる
	r, err := os.Open("./test.txt")
	s2 := bufio.NewScanner(r)
	for s2.Scan() {
		fmt.Println(s2.Text())
	}
	if err := s2.Err(); err != nil {
		// エラー処理
		fmt.Println(err)
	}

	//コードポイント(rune)ずつ読み込むScannerを作る問題
	s3 := NewRuneScanner(strings.NewReader("Hello, 世界"))
	for {
		r, err := s3.Scan()
		// fmt.Println("r = ", r)//r =  101 etc
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%c\n", r)
	}

	// ** エラーをまとめる */ ・・https://github.com/uber-go/multierrを使う
	// 	成功したものは成功させたい
	// 失敗したものだけエラーとして報告したい
	// N番目エラーはどういうエラーなのか知れる
	// 	var rerr error
	// 	if err := step1(); err != nil {
	// 		rerr = multierr.Append(rerr, err)
	// 	}
	// 	if err := step2(); err != nil {
	// 		rerr = multierr.Append(rerr, err)
	// 	}
	// 	return rerr
	// 	for _, err := range multierr.Errors(rerr) {
	// 		fmt.Println(err)
	//  }

	// **エラーに文脈を持たせる */ ・・github.com/pkg/errorsを使う
	// 	エラーメッセージが「File Not Found」とかでは分かりづらい
	// 何をしようとした時にエラーが起きたか知りたい
	// どんなパラメータだったのか知りたい
	// errors.Wrapを使うとエラーをラップできる
	// errors.Causeを使うと元のエラーが取得できる
	// if err := f(s); err != nil {
	// 	return errors.Wrapf(err, "f() with %s", s)
	// }

	// ** エラーに文脈を持たせる（Go1.13） */・・fmt.Errorf関数の%wを使う
	// 引数で指定したエラーをラップしてエラーを作る
	// Unwrapメソッドを実装したエラーが作られる
	// errors.Unwrap関数で元のエラーが取得できる
	err2 := fmt.Errorf("bar: %w", errors.New("foo"))
	fmt.Println(err2)                // bar: foo
	fmt.Println(errors.Unwrap(err2)) // foo

	// ** 値によって分岐する（Go1.13） */ ・・errors.Is関数を使う
	// 第1引数のエラーが第2引数の値かどうか判定する
	// ==で比較できる場合は比較
	// Isメソッドを実装している場合はそれで比較
	// 判定不能の場合はerrors.Unwrap関数を呼んでアンラップ後に判定
	if errors.Is(err, os.ErrExist) {
		// os.ErrExistだった場合の処理
	}

	// ** エラーから情報を取り出す（Go1.13）*/ ・・errors.As関数を用いる
	// 	第1引数で指定したエラーを第2引数で指定したポインタが指す変数に代入する
	// キャスト不可な場合はerrors.Unwrap関数でアンラップ後に試みる
	var pathError *os.PathError
	if errors.As(err, &pathError) {
		fmt.Println("Failed at path:", pathError.Path)
	} else {
		fmt.Println(err) // <nil>
	}

	//** エラーとログ */
	// ** エラーメッセージを工夫する
	// ログに出すエラーメッセージに必要十分な情報を入れる
	// スタックトレースを付加する
	// pkg/errors.WithStackすると付加される
	// pkg/erros.Wrapでも可
	// xerrosでも付加されるがerrorsでは付加されない（Go1.13）
	// ** ログの出力がボトルネックにならないように
	// なんでもログに出せばよいという訳ではない
	// リフレクションを使うような処理をやりすぎない
	// ログ出力でサーバに負荷を与え過ぎないように
	// 例：桁数がめちゃくちゃでかいbig.Floatの出力など
	// zapなどの高速なログライブラリを使う
	// 高速なログライブラリ
	// https://pkg.go.dev/go.uber.org/zap

	// ** パニックとリカバー */
	// ** パニック
	// 回復不能だと判断された実行時のエラーを発生させる機構
	// 組み込み関数のpanicを呼び出すと発生する
	// ** リカバー
	// 発生したパニックを取得し、エラー処理を行う
	// recover関数をdeferで呼び出された関数内で実行する
	// 関数単位でしかリカバーできない
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r) // "ERROR"
		}
	}()
	// panic("ERROR") // 文字列じゃなくてもいい

	// ** エラーとパニック */ ・・エラーとパニックの使い分け
	// パニックは回避不可能な場合のみ使用する
	// 想定内のエラーはerror型で処理する

	// ** Must*関数 */・・パッケージ初期化時のエラー処理に用いる
	// エラーではなくパニックを発生させる
	// 実行直後にパニックが発生する
	// 正規表現やテンプレートエンジンの初期化関数に設けられている
	// パッケージの初期化時に行う
	var validID = regexp.MustCompile(`^[a-z]+\[[0-9]+\]$`)
	fmt.Println(validID.MatchString("adam[23]")) //true

	// 関数内で行う場合はエラー処理をする
	validID2, err := regexp.Compile(`^[a-z]+\[[0-9]+\]$`)
	if err != nil { /* エラー処理 */
	}
	fmt.Println(validID2.MatchString("adam[23]")) //true

	// ** 名前付き戻り値とパニック */ ・・パニックで渡された値を戻り値として返す
	// if err := f(); err != nil {
	// 	log.Fatal(err)//2022/05/01 19:59:49 error
	// }

	//** 大域脱出のテクニックとして使う */ ・・
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(escape); ok { //大域脱出されてきたか？
				println("Escaped") //Escaped
			} else {
				panic(r) //それ以外はそのままパニック

			}
		}
	}()
	fg() //Escaped
}

// ** エラー */ ・・errorインタフェース
type error interface {
	Error() string
}

// ** エラー型の定義 */ ・・Errorメソッドを実装している型を定義する
//そのエラー特有の情報を保持する
type PathError struct {
	Op   string
	Path string
	Err  error
}

func (e *PathError) Error() string {
	return e.Op + " " + e.Path + ": " + e.Err.Error()
}

// ** 振る舞いでエラーを処理する */ ・・エラー処理は具象型に依存させない
// エラーの種類で処理を分けたい場合がある
// インタフェースを使い振る舞いで処理する
// 一時的なエラーかどうかを判定する関数
func IsTemporary(err error) bool { //error型であることが重要
	te, ok := err.(interface {
		Temporary() bool //Temporaryメソッドを持っているか？
	})
	return ok && te.Temporary()
}

// ** 名前付き戻り値とパニック */ ・・パニックで渡された値を戻り値として返す
func f() (rerr error) {
	defer func() {
		if r := recover(); r != nil {
			if err, isErr := r.(error); isErr {
				rerr = err //2022/05/01 19:59:49 error
			} else {
				rerr = fmt.Errorf("%v", r) //2022/05/01 19:59:49 error
			}
		}
	}()
	panic(errors.New("error"))
	return nil
}

type Stringer interface {
	String() string
}
type S string

func (s S) String() string {
	return string(s)
}

//ユーザー定義型のエラーインターフェース
type MyError string

func (e MyError) Error() string {
	return string(e)
}

//任意の値をStringer型に変換する関数
func ToStringer(x interface{}) (Stringer, error) {
	if s, ok := x.(Stringer); ok {
		return s, nil
	}
	return nil, MyError("CastError")
}

//コードポイント(rune)ずつ読み込むScannerを作る問題
type RuneScanner struct {
	r   io.Reader
	buf [16]byte
}

func NewRuneScanner(r io.Reader) *RuneScanner {
	return &RuneScanner{r: r}
}

func (s *RuneScanner) Scan() (rune, error) {
	n, err := s.r.Read(s.buf[:]) //与えられたバイトスライス p []byte を先頭から埋めていく。埋まったバイト数 n と、埋める過程で発生したエラー err を返す。
	// fmt.Println("n=",n)
	if err != nil {
		return 0, err
	}

	r, size := utf8.DecodeRune(s.buf[:n])
	// fmt.Println("r, size = ", r, size)
	if r == utf8.RuneError {
		return 0, errors.New("RuneError")
	}

	s.r = io.MultiReader(bytes.NewReader(s.buf[size:n]), s.r)
	// fmt.Println("s.r = ",s.r)
	// fmt.Println("r = ",r)
	return r, nil
}

// 大域脱出のテクニックとして使う
type escape struct{} //パッケージ内の型にする

func fg() { g() }
func g()  { panic(escape{}) }
