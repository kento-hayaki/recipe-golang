package mypkg_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/jinzhu/configor"
)

type Hex int

func (h Hex) String() string {
	return fmt.Sprintf("%x", int(h))
}

func IsOdd(i int) bool {
	return i%2 != 0
}

func TestHex_String(t *testing.T) {
	want := "a"
	got := Hex(10).String()
	if got != want {
		// fmt.Println(want, got)// 出力されない
		t.Errorf("want %q, got %q", want, got)
	}
}

// ** testingパッケージでできること */
// 失敗理由を出力してテストを失敗させる
// テスト関数を継続：t.Error, t.Errorf
// テスト関数を終了：t.Fatal, t.Fatalf
// testing.T: TestXxxの引数
// t.Fatal: テストを終了し、失敗内容を出力する
// t.Error: テストを継続しつつ、失敗内容を出力する
// t.Run: サブテストを実行できる
// t.Parallel: テストを並列に実行する
// テストの並列実行
// t.Parallel（テスト関数の先頭で呼び出す）
// go testの-parallelオプションで並列数を指定
// ベンチマーク
// *testing.B型を使う
//b.ResetTimer: 経過時間とメモリカウンターをリセットする
// ブラックボックステスト
// testing/quickパッケージ
// あまり積極的には使わない
// go test -bench .　でテスト実行
func BenchmarkFizzBuzz(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// fizzbuzz.FizzBuzz(i)
	}
}

// ** テストの後処理
// t.Cleanup
// テスト終了時に行う関数を登録
// defer文だと並列実行してる際におかしくなる
// t.TempDir
// テスト終了時に消える一時ディレクトリを作成

// ** Exampleテスト
// テストされたサンプル
// Exampleで始まる関数を書く
// Go Docにサンプルとして出る
// Output: を書くとテストになる
func ExampleHex_String() {
	fmt.Println(Hex(10))
	//（1）は失敗
	// Output: a

}

// ** テストの並列実行 */ ・・t.Parallelメソッドを用いる
// テスト関数を並列に実行する許可を与える
// 初期段階で導入を検討しておく
// テストが遅くなってから導入すると大変
func Test(t *testing.T) {
	t.Parallel()
	/* (略) */
}

// ** サブテスト */・・子テストを実行するしくみ
// サブテストを指定して実行できる
// func TestAB(t *testing.T) {
// 	t.Run("A", func(t *testing.T) { t.Error("error") })
// 	t.Run("B", func(t *testing.T) { t.Error("error") })
// }

// 実行方法 go test -v mypkg_test.go -run TestAB/A

// ** テーブル駆動テスト */・・テスト対象のデータを羅列してテストする
var flagtests = []struct {
	in  string
	out string
}{
	{"%a", "[%a]"}, {"%-a", "[%-a]"}, {"%+a", "[%+a]"},
	{"%#a", "[%#a]"}, {"% a", "[% a]"},
}

// func TestFlagParser(t *testing.T) {
// 	var flagprinter flagPrinter
// 	for _, tt := range flagtests {
// 		s := Sprintf(tt.in, &flagprinter)
// 		if s != tt.out {
// 			t.Errorf("Sprintf(%q, &flagprinter) => %q, want %q", tt.in, s, tt.out)
// 		}
// 	}
// }

// ** サブテストとテーブル駆動テスト
// 利点
// 落ちた箇所が分かりやすい
// テストケースの名前が表示される
// 特定のサブテストだけ実行できる
// テストケースが大量な場合分かりやすい

func TestIsOdd(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		in   int
		want bool
	}{
		"+odd":  {5, true},
		"+even": {6, false},
		"-odd":  {-5, true},
		"-even": {-6, false},
		"zero":  {0, false},
		//"error": {1, false},//error
	}
	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got := IsOdd(tt.in); tt.want != got {
				t.Errorf("want IsOdd(%d) = %v, got %v", tt.in, tt.want, got)
			}
		})
	}
}

// ** テストヘルパー */・・テスト用のヘルパー関数
// ヘルパー関数はエラーを返さない
// *testing.Tを受け取ってテストを落とす
// (*testing.T).Helperを使って情報を補足する
// https://golang.org/pkg/testing/#T.Helper
// https://play.golang.org/p/uj-vKdHol-k
func testTempFile(t *testing.T) string {
	t.Helper()
	tf, err := ioutil.TempFile("", "test")
	if err != nil {
		//t.Fatal("err %s", err) //error
		t.Fatalf("err %s", err)
	}
	tf.Close()
	return tf.Name()
}

// ** coverprofile */ ・・テストのカバレッジを分析
//テストがどれだけ網羅的に行われたか調べる

// ** テスタビリティ * /

// ** テスタブルなコードと抽象化 */
// - テストしやすいコード
// -- 個々の機能が疎結合で単体でテストしやすい
// -- 外部との接続部分が抽象化されている
// --- データベース接続、ネットワークやファイルへのアクセス
// -- 抽象化されている部分をモックに差し替えれる
// --- moq
// --- インタフェース

// ** インタフェースを使う */・・外部に繋がる部分はモックに差し替え可能にする
type DB interface {
	Get(key string) string
	Set(key, value string) error
}

// DBはインタフェースなので実装を入れ替えれる
type Server struct{ DB DB }

// ** テストする部分だけ実装する
// 埋め込みを使って一部分だけモックを用意する
// 呼び出さないメソッドは実装しなくてもコンパイルエラーにならない
type getTestDB struct {
	// DBを埋め込むことで実装したことになる
	DB
}

// Getだけテスト用に実装する
func (db getTestDB) Get(key string) string { return "test" }

// **環境変数を使う
// 環境変数を使って切り替える
// os.Getenvで取得できる
// CIでテストを走らせるときに便利
// DBの接続先など環境に依存する値を保存する
// github.com/jinzhu/configorを使う
var Config = struct {
	DB string `env:"DB"`
}{}

func main() {
	configor.Load(&Config)
	fmt.Printf("config: %#v", Config)
}

// ** テストデータを用意する */ ・・どの環境でも使用できるテストデータを用意する
// - testdataというディレクトリに入れる（参考）https://pkg.go.dev/cmd/go#hdr-Test_packages
// testdataはパッケージとみなされない
// - テストの再現性を担保する
// ネットワークアクセスを発生させない
// テストデータ以外のファイルにアクセスしない

// ** 非公開な機能を使ったテスト
// - export_test.goという名前でファイルを作る
// _testがついているのでテストのときにしかビルドされない
// - テスト対象のパッケージと同じにする
// テストのときだけ参照できる関数などを作る

// mypkg.go
// package mypkg

// const maxValue = 100

// export_test.go
// package mypkg // テスト対象と同じパッケージ

// const ExportMaxValue = maxValue
