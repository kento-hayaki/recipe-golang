package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// 単純なhttp_server
// ** HTTPハンドラ
// - 引数にレスポンスを書き込む先とリクエストを取る
// - 第1引数はレスポンスを書き込む先
// 書き込みにはfmtパッケージの関数などが使える
// - 第2引数はクライアントからのリクエスト
func handler(w http.ResponseWriter, r *http.Request) { // (レスポンスを書き込むWriter, リクエスト)
	fmt.Fprint(w, "Hello, HTTPサーバ") //レスポンスの書き込み
}

// **http.Handlerインタフェース
// HTTPハンドラはインタフェースとして定義されている
// ServeHTTPメソッドを実装していればハンドラとして扱われる
type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// ** http.Handleでハンドラを登録
// パターンとhttp.Handlerを指定して登録する
// 第1引数としてパターンを指定する
// 第2引数としてhttp.Handlerを指定する
// http.DefaultServeMuxに登録される
// **ServeHTTPメソッドを持つ型がハンドラとして扱われる
func Handle(pattern string, handler http.Handler) { // 実際には、ServeHTTPメソッドを持つ型の具体的な値がくる

}

// ** http.HandleFuncでハンドラを登録・・パターンと関数を指定して登録する
// 第1引数としてパターンを指定する
// 第2引数として関数を指定する
// http.DefaultServeMuxに登録される
func HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) { // http.HandlerのServeHTTPメソッドと同じ引数の関数
}

// ** http.HandlerFuncとは */ ・・関数にhttp.Handlerを実装させている
type HandlerFunc func(http.ResponseWriter, *http.Request)

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

// ** http.HandleFuncのしくみ */ ・・・http.HandleFunc
// 引数で受け取った関数をhttp.HandlerFuncに変換する
// http.Handleでhttp.Handlerとして登録する
// ** Handlerは登録されるもの。Handleは登録する関数
// func HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {}

// ** http.ServeMux
// 複数のハンドラをまとめる
// パスによって使うハンドラを切り替える
// 自身もhttp.Handlerを実装している
// http.Handleとhttp.HandleFuncはデフォルトのhttp.ServeMuxであるhttp.DefaultServeMuxを使用している

// ** レスポンスとリクエスト
// ** http.ResponseWriterについて */ ・・http.ResponseWriterインタフェース
// - io.Writerと同じWriteメソッドをもつ
// ResponseWriterを満たすとio.Writerを満たす
// - io.Writerとしても振る舞える
// fmt.Fprint*の引数に取れる
// json.NewEncoderの引数に取れる
// インタフェースなのでモックも作りやすい＝テスト簡単

// ** エラーを返す */ ・・http.Error関数を使う
// エラーメッセージとステータスコードを指定する
// ステータスコードは定数としてhttpパッケージで定義されている
// http.StatusOKやhttp.StatusInternalServerError
func Error(w http.ResponseWriter, error string, code int) {}

func main() {
	http.HandleFunc("/", handler)

	// ** レスポンスとリクエスト
	// ** JSONを返す
	// encoding/jsonパッケージを使う
	// 機械的に処理しやすいJSONをレスポンスに用いる場合も多い
	// JSONエンコーダを使ってGoの値をJSONに変換する
	// 構造体をやスライスをJSONのオブジェクトや配列にできる
	type Person struct {
		Name string `json:"name"` // 構造体のタグでJSONのフィールド名を指定
		Age  int    `json:"age"`
	}

	p := &Person{Name: "tenntenn", Age: 31}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(p); err != nil {
		log.Fatal(err)
	}
	fmt.Println(buf.String()) //{"name":"tenntenn","age":31}

	// ** JSONデコード */ json.Decoder型を使う
	var p2 Person
	dec := json.NewDecoder(&buf)
	if err := dec.Decode(&p2); err != nil {
		log.Fatal(err)
	}
	fmt.Println(p2) // {tenntenn 31}

	// ** リクエストパラメタの取得 */ ・・(*http.Request).FormValueから名前を指定して取得
	// - パラメタ指定の例：http://localhost:8080?msg=Gophers
	// - 複数ある場合は&でつなぐ
	// http://localhost:8080?a=100&b=200
	http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello", r.FormValue("msg"))
	})
	//http://localhost:8080/query?msg=test // hello test

	// ** リクエストボディの取得 */・・(*http.Request).Bodyから取得する
	// io.ReadCloserを実装している
	http.HandleFunc("/body", func(w http.ResponseWriter, r *http.Request) {
		var p Person
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&p); err != nil {
			// ... エラー処理 …
		}
		fmt.Println(p) // { 0}
	})

	// ** リクエストヘッダーを取得 */・・RequestのHeaderフィールドを使う
	// Getメソッドを使うとヘッダー名を指定して取得できる
	// func handler3(w http.ResponseWriter, req *http.Request) {
	// 	contentType := req.Header.Get("Content-Type")
	// 	fmt.Fprintln(w, contentType)
	// }

	// ** テンプレートエンジンの利用 */
	// - html/templateを使う
	// Go標準のテンプレートエンジン
	// text/templateのHTML特化版
	// - テンプレートの生成
	template.Must(template.New("sign"). // テンプレート名
						Parse("<html><body>{{.}}</body></html>")) // HTML
	// テンプレートに埋め込む (io.Writer, リクエストから貰った値を埋め込む)
	// tmpl.Execute(w, r.FormValue("content"))

	// ** よく使うテンプレートの記法
	// その文脈でトップレベルのデータを埋め込む
	// {{.}}
	// フィールドやメソッド
	// {{.Filed}}
	// {{.Method arg1 arg2}}
	// 条件分岐
	// {{if .}}{{.Filed}}{{else}}NO{{end}}
	// 繰り返し rangeの中の{{.}}は要素になる
	// {{range .}}{{.}}{{end}

	// ** おみくじアプリを作ろう課題用 ** //
	// 現在時刻
	t := time.Now().UnixNano()
	rand.Seed(t)

	http.HandleFunc("/omikuji", handler3)

	// **10.4. HTTPクライアント **/
	// **HTTPリクエストを送る //・・http.DefaultClientを用いる
	// デフォルトのHTTPクライアント
	// http.Getやhttp.Postはhttp.DefaultClientのラッパー
	// resp, err := http.Get("http://example.com/")
	// resp2, err := http.Post("http://example.com/upload", "image/jpeg", &buf)

	// v := url.Values{"key": {"Value"}, "id": {"123"}}
	// resp3, err := http.PostForm("http://example.com/form", v)

	// ** レスポンスを読み取る*/ ・・(*http.Response).Bodyを使う
	// io.ReadCloserを実装している
	// 読み込んだらCloseメソッドを呼ぶ
	resp, err := http.Get("http://example.com/")
	if err != nil { /* エラー処理 */
	}
	defer resp.Body.Close()
	var p3 Person
	//fmt.Println(resp.Body)
	dec2 := json.NewDecoder(resp.Body)
	if err := dec2.Decode(&p3); err != nil {
		// ... エラー処理 …
	}
	fmt.Println("p=", p) // p = &{tenntenn 31}

	// **リクエストを指定する・・http.Client.Doを用いる
	// 引数に*http.Requestを渡すことができる
	req, err := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Add("If-None-Match", `W/"wyzzy"`)
	client := http.Client{}
	resp2, err := client.Do(req)
	fmt.Println("resp2=", resp2)

	// ** リクエストとコンテキスト */・・リクエストとコンテキスト
	// *http.Requestから取得する（サーバ）
	// ctx := req.Context()
	// コンテキストを更新する（クライアント）
	// 新しい*http.Requestが生成される
	// req = req.WithConntext(ctx)

	// ** http.Clientとhttp.Transport */・・http.Transport型
	// http.RoundTriperを実装した型
	// HTTP/HTTPS/HTTPプロキシに対応している
	// コネクションのキャッシュを行う
	// 実際にHTTP通信ところ
	// ** http.DefaultTransport
	// http.ClientのTransportフィールドがnilの時に使われる
	var DefaultTransport http.RoundTripper = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	fmt.Println("DefaultTransport = ", DefaultTransport)

	// ** http.RondTripper **/・・HTTPのトランザクションを行うインタフェース
	// - 実装している型
	// *http.Transport
	// http.NewFileTransportで取得できる値のdynamic type
	// - http.RoundTripperを実装する場面
	// リクエストの度に前処理や後処理を行いたい場合
	// モックサーバを作りたい場合
	type RoundTripper interface {
		RoundTrip(*http.Request) (*http.Response, error)
	}
	// ** http.RoundTriperを実装する場合のTIPS
	// - 注意点
	// レスポンスを返す場合はエラーはnilにすること
	// リクエストは変更しない
	// リクエストは他のゴールーチンから参照される可能性がある
	// - TIPS
	// 元になるhttp.RoundTripperをラップしておく
	// フィールドで設定できるようにしておく
	// HTTP通信の部分は親のRoundTripメソッドを呼ぶ
	// フィールドがnilの場合はhttp.DefaultTransportを使う

	// ** Q. http.GetでおみくじWebアプリにリクエストを送ってみよう
	resp4, err := http.Get("http://localhost:8080?p=Gopher")
	if err != nil { /* エラー処理 */
		fmt.Println("接続できませんでした")
	}
	fmt.Println(&resp4)

	// ** HTTPサーバの起動 */ ・・http.ListenAndServeを使う
	// - 第1引数でホスト名とポート番号を指定
	// ホスト名を省略した場合localhost
	// - 第2引数でHTTPハンドラを指定
	// nilで省略した場合はhttp.HandleFuncなどで登録したハンドラが使用される
	http.ListenAndServe(":8080", nil) //(ホスト名:ポート番号, HTTPハンドラ)

}

// ** レスポンスヘッダーを設定する */・・ResponseWriterのHeaderメソッドを使う
// WriteやWriteHeaderを呼び出した後に設定しても効果がない
func handler2(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	v := struct {
		Msg string `json:"msg"`
	}{
		Msg: "hello",
	}
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Println("Error:", err)
	}
}

// ** リダイレクト */ ・・http.Redirect関数を使う
// 第3引数に遷移したいパスを指定する
// 第4引数に3xx系のステータスコードを指定する
func handler4(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusFound)
}

// ** ミドルウェアを作る */ ・・ハンドラより前に行う共通処理
// ライブラリを使ってもOK
type MiddleWare interface {
	ServeNext(h http.Handler) http.Handler
}
type MiddleWareFunc func(h http.Handler) http.Handler

func (f MiddleWareFunc) ServeNext(h http.Handler) http.Handler {
	return f(h)
}
func With(h http.Handler, ms ...MiddleWare) http.Handler {
	for _, m := range ms {
		h = m.ServeNext(h)
	}
	return h
}

// ** HTTPハンドラのテスト */・・net/http/httptestを使う
// ハンドラのテストのための機能など提供
// httptest.ResponseRecorder
// http.ResponseWriterインタフェースを実装している
// NewRequestメソッド(1.7以上)
// 簡単にテスト用のリクエストが作れる
// **ハンドラのテストの例
func TestSample(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	handler(w, r)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatal("unexpected status code")
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal("unexpected error")
	}
	const expected = "Hello, HTTPサーバ"
	if s := string(b); s != expected {
		t.Fatalf("unexpected response: %s", s)
	}
}

// ** おみくじアプリを作ろう課題用 ** //
var tmpl = template.Must(template.New("msg").
	Parse("<html><body>{{.Name}}さんの運勢は「<b>{{.Omikuji}}</b>」です</body></html>"))

type Result struct {
	Name    string
	Omikuji string
}

func handler3(w http.ResponseWriter, r *http.Request) {
	result := Result{
		Name:    r.FormValue("p"),
		Omikuji: omikuji(),
	}
	tmpl.Execute(w, result)
}

func omikuji() string {
	n := rand.Intn(6) // 0-5
	switch n + 1 {
	case 6:
		return "大吉"
	case 5, 4:
		return "中吉"
	case 3, 2:
		return "小吉"
	default:
		return "凶"
	}
}

// ** おみくじアプリを作ろう課題用 ここまで ** //
