package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	// ** 1行ずつ読み込む
	// bufio.Scannerを使用する
	// 標準入力から読み込む
	scanner := bufio.NewScanner(os.Stdin)
	// 1行ずつ読み込んで繰り返す
	// for scanner.Scan() {
	// 	1行分を出力する
	// 	fmt.Println(scanner.Text())
	// }
	// まとめてエラー処理をする
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "読み込みに失敗しました:", err)
	}

	// ** SplitFunc型 */ ・・分割するアルゴリズムを表す型
	// *bufio.Scanner型のSplitメソッドで設定する
	type SplitFunc func(data []byte, atEOF bool) (
		/* 戻り値 */ advance int, token []byte, err error)

	scanner2 := bufio.NewScanner(os.Stdin)
	scanner2.Split(bufio.ScanBytes) // 1バイトごと
	scanner2.Split(bufio.ScanLines) // 1行ごと（デフォルト）
	scanner2.Split(bufio.ScanRunes) // コードポイントごと
	scanner2.Split(bufio.ScanWords) // 1単語ごと

	// ** strconvパッケージ */・・文字列と他の型の変換を行うパッケージ
	// 文字列をint型に変換: 100 <nil>
	fmt.Println(strconv.Atoi("100"))

	// int型を文字列に変換: 100円
	fmt.Println(strconv.Itoa(100) + "円")

	// 100を16進数で文字列にする: 64
	fmt.Println(strconv.FormatInt(100, 16))

	// 文字列をbool型にする: true <nil>
	fmt.Println(strconv.ParseBool("true"))

	// ** 数値へ変換時の注意点 */ ・・strconv.Atoi関数で変換した値のキャスト
	// オーバーフローを起こすサイズにキャストしてもpanicにならない
	// 変換後のint型からint16型などにキャストしない
	// 最初からstrconv.ParseInt型を用いる
	// gosecなどで検出する
	// int16より大きな値:"32768"
	s := strconv.FormatInt(math.MaxInt16+1, 10)
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	if int16(n) < 0 { // オーバーフロー
		fmt.Println(n) // 32768 が表示される
	}

	// ** stringsパッケージ */ ・・文字列関連の処理を行うパッケージ
	// スペースで分割してスライスにする: [a b c]
	fmt.Println(strings.Split("a b c", " "))

	// スライスを","で結合する: a,b,c
	fmt.Println(strings.Join([]string{"a", "b", "c"}, ","))

	// 繰り返す: hogehoge
	fmt.Println(strings.Repeat("hoge", 2))

	// プリフィックスを持つかどうか: true
	fmt.Println(strings.HasPrefix("hoge_fuga", "hoge"))

	// ** 文字列の置換・・strings.Replace関数を使う
	// 第1引数の文字列中の第2引数の文字列を第3引数の文字列で置換
	// 第4引数は置換回数で-1を指定するとすべて置換
	s1 := strings.Replace("郷に入っては郷に従え", "郷", "Go", 1)
	// Goに入っては郷に従え
	fmt.Println(s1)

	s2 := strings.Replace("郷に入っては郷に従え", "郷", "Go", -1)
	// Goに入ってはGoに従え
	fmt.Println(s2)

	s3 := strings.ReplaceAll("郷に入っては郷に従え", "郷", "Go")
	// Goに入ってはGoに従え
	fmt.Println(s3)

	// ** 複数文字列の置換 */ ・・strings.Replacer型を使う
	// strings.Replacer型を使う
	// strings.NewReplacer関数で置き換えたい文字列のペアを指定する
	// 実際に置換を行うのはReplaceメソッド
	// Writeメソッドを用いるとio.Writerに書き出せる
	// 郷 → Go、入れば → 入っては
	r := strings.NewReplacer("郷", "Go", "入れば", "入っては")
	s4 := r.Replace("郷に入れば郷に従え")
	fmt.Println(s4) // Goに入ってはGoに従え

	_, err = r.WriteString(os.Stdout, "郷に入れば郷に従え") // Goに入ってはGoに従え
	if err != nil {                                /* エラー処理 */
	}

	// ** コードポイント（rune）ごとの置換・・strings.Map関数を使う
	// 第1引数にrune型単位で置換する関数
	// 第2引数に置換したい文字列
	// 小文字を大文字に変換する関数
	toUpper := func(r rune) rune {
		if 'a' <= r && r <= 'z' {
			return r - ('a' - 'A')
		}
		return r
	}
	// HELLO, WORLD
	s5 := strings.Map(toUpper, "Hello, World")
	fmt.Println(s5)

	// ** 大文字・小文字の変換
	// 	unicode.ToUpper関数 / unicode.ToLower関数
	// rune単位で大文字・小文字に変換する関数
	// strings.ToUpper / strings.ToLower関数
	// 文字列単位で大文字・小文字に変換する関数
	s6 := strings.Map(unicode.ToUpper, "Hello, World")
	// HELLO, WORLD
	fmt.Println(s6)
	// hello, world
	fmt.Println(strings.Map(unicode.ToLower, s6))
	// HELLO, WORLD
	fmt.Println(strings.ToUpper("Hello, World"))
	// hello, world
	fmt.Println(strings.ToLower("Hello, World"))

	// ** bytesパッケージ */ ・・[]byte型向けの処理を提供する
	// stringsパッケージにある関数や型に似たものが多い
	// []byte型からstring型へのキャスト省く
	// 0x0B → 0xFF
	src := []byte{0x0A, 0x0B, 0x0C}
	b := bytes.ReplaceAll(src, []byte{0x0B}, []byte{0xFF})
	// 0A FF 0C
	fmt.Printf("% X\n", b)

	// ** 12.2. ioパッケージ */・・io.Pipe関数
	// 	パイプのように接続されたReaderとWriterを作る
	// io.Pipe関数で*io.PipeReader型と*io.PipeWriter型の値を生成
	r2, w2 := io.Pipe()
	go func() {
		// fmt.Fprint(w2, "Hello, 世界\n")//Hello, 世界
		w2.Close()
	}()
	io.Copy(os.Stdout, r2)

	// ** 読み込みバイト数を制限する・・io.LimitedReader型を用いる
	// Rフィールドには元のio.Readerを設定する
	// Nフィールドには読み込むバイト数を設定する
	r3 := &io.LimitedReader{
		R: strings.NewReader("Hello, 世界"),
		N: 5,
	}
	// Hello
	io.Copy(os.Stdout, r3)

	// **複数のio.Writerに書き込む
	// io.MultiWriter関数を用いる
	// 同じ内容が複数のio.Writerに書き込まれる
	var buf1, buf2 bytes.Buffer
	w := io.MultiWriter(&buf1, &buf2)
	fmt.Fprint(w, "Hello, 世界")
	// buf1: Hello, 世界
	fmt.Println("buf1:", buf1.String())
	// buf2: Hello, 世界
	fmt.Println("buf2:", buf2.String())

	// ** 複数のio.Readerから読み込む */ ・・io.MultiReader関数を用いる
	// 複数のio.Readerを直列につなげたようなio.Readerを生成
	// 分割された複数のファイルから読み込む場合などに一度にメモリに載せなくて済む
	// すでに読み込んだ部分を先頭に詰めるなどに応用できる
	r4 := strings.NewReader("Hello, ")
	r5 := strings.NewReader("世界\n")
	r6 := io.MultiReader(r4, r5)
	// Hello, 世界
	io.Copy(os.Stdout, r6)

	// ** io.TeeReader関数・・読み込みと同時に書き込むio.Readerを作る
	// 引数のio.Readerをベースに読み込まれると同時に引数のio.Writerに書き込む
	var buf3 bytes.Buffer
	r7 := strings.NewReader("Hellooooo, 世界\n")
	tee := io.TeeReader(r7, &buf3)
	// Hellooooo, 世界
	io.Copy(os.Stdout, tee) // bufにも書き込まれる
	// Hellooooo, 世界
	fmt.Print(buf3.String())

	// **12.3. 正規表現*/
	// ** 正規表現のコンパイル
	// - regexp.Compile関数を用いる
	// パッケージ変数で1度しか行わない場合はMustCompile関数を使う
	// *regexp.Regexp型が返される
	// 使えるシンタックス：https://golang.org/s/re2syntax
	var validID = regexp.MustCompile(`^[a-z]+\[[0-9]+\]$`)
	fmt.Println(validID.MatchString("adam[23]")) //true
	// 関数内で行う場合はエラー処理をする
	validID2, err := regexp.Compile(`^[a-z]+\[[0-9]+\]$`)
	if err != nil { /* エラー処理 */
	}
	fmt.Println(validID2.MatchString("adam[23]")) //true

	// ** 正規表現のマッチ */ ・・指定した文字列などが正規表現にマッチするか
	// MatchメソッドやMatchStringメソッドを使う
	re, err := regexp.Compile(`(\d+)年(\d+)月(\d+)日`)
	if err != nil { /* エラー処理 */
	}
	// バイト列（[]byte型）がマッチするか
	fmt.Println(re.Match([]byte("1986年01月12日"))) // true
	// 文字列がマッチするか
	fmt.Println(re.MatchString("1986年01月12日")) // true
	// io.RuneReaderがマッチするか
	var r8 io.RuneReader = strings.NewReader("1986年01月12日") // true
	fmt.Println(re.MatchReader(r8))

	// **マッチした部分を返す*/ ・・正規表現にマッチする文字列などを探す
	// FindメソッドやFindStringメソッドを用いる
	// FindAllメソッドやFindStringAllメソッドは個数を指定できる
	// -1はマッチするすべてを取得
	re2, err := regexp.Compile(`\d+`)
	if err != nil { /* エラー処理 */
	}
	// 最初にマッチするバイト列を取得
	fmt.Printf("%q\n", re2.Find([]byte("1986年01月12日"))) //"1986"
	// すべてのマッチするバイト列を取得
	fmt.Printf("%q\n", re2.FindAll([]byte("1986年01月12日"), -1)) //["1986" "01" "12"]
	// 最初にマッチする文字列を取得
	fmt.Printf("%q\n", re2.FindString("1986年01月12日")) //"1986"
	// すべてのマッチする文字列を取得
	fmt.Printf("%q\n", re2.FindAllString("1986年01月12日", -1)) //["1986" "01" "12"]

	// ** マッチした部分のインデックスを返す*/・・正規表現にマッチする部分のインデックスを返す
	// Find*Indexメソッドを用いる
	// マッチする部分の始端と終端のインデックスが返ってくる
	re3, err := regexp.Compile(`\d+`)
	if err != nil { /* エラー処理 */
	}
	// [0 4]
	fmt.Println(re3.FindIndex([]byte("1986年01月12日")))
	// [[0 4] [7 9] [12 14]]
	fmt.Println(re3.FindAllIndex([]byte("1986年01月12日"), -1))
	// [0 4]
	fmt.Println(re3.FindStringIndex("1986年01月12日"))
	// [[0 4] [7 9] [12 14]]
	fmt.Println(re3.FindAllStringIndex("1986年01月12日", -1))

	// ** キャプチャされた部分を取得*/・・Find*Submatch*メソッドを使う
	re4, err := regexp.Compile(`(\d+)[^\d]`) // 数字がキャプチャされる
	if err != nil {                          /* エラー処理 */
	}
	// ["1986年" "1986"]
	fmt.Printf("%q\n", re4.FindSubmatch([]byte("1986年01月12日")))
	fmt.Printf("%q\n", re4.FindStringSubmatch("1986年01月12日"))
	// [["1986年" "1986"] ["01月" "01"] ["12日" "12"]]
	fmt.Printf("%q\n", re4.FindAllSubmatch([]byte("1986年01月12日"), -1))
	fmt.Printf("%q\n", re4.FindAllStringSubmatch("1986年01月12日", -1))
	// [0 7 0 4]
	fmt.Println(re4.FindSubmatchIndex([]byte("1986年01月12日")))
	fmt.Println(re4.FindStringSubmatchIndex("1986年01月12日"))
	// [[0 7 0 4] [7 12 7 9] [12 17 12 14]]
	fmt.Println(re4.FindAllSubmatchIndex([]byte("1986年01月12日"), -1))
	fmt.Println(re4.FindAllStringSubmatchIndex("1986年01月12日", -1))

	// ** キャプチャした部分の展開*/・・キャプチャした部分をテンプレートに展開する
	// ExpandメソッドやExpandStringメソッドを使う
	// FindAllStringSubmatchIndexメソッドなどでインデックスを取得する
	// (?P<var_name>regexp)で名前をつけてキャプチャする
	re5, err := regexp.Compile(`(?P<Y>\d+)年(?P<M>\d+)月(?P<D>\d+)日`)
	if err != nil { /* エラー処理 */
	}
	content := "1986年01月12日\n2020年03月24日"
	template := "$Y/$M/$D\n" // "${1}/${2}/${3}"でも可
	var result []byte
	for _, submatches := range re5.FindAllStringSubmatchIndex(content, -1) {
		result = re5.ExpandString(result, template, content, submatches)
	}
	// "1986/01/12\n2020/03/24\n"
	fmt.Printf("%q", result)

}
