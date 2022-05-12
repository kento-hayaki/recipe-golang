package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	//fmt.Println(mylib.Hello())
	//fmt.Println(greeting.Do())
	var f float64 = 10
	var n int = int(f) // 型キャスト・・ある方から別の型に変更すること
	println(n)

	// ** 型リテラル */ リテラル = 識別子（名前）が付与されてないもの
	// int型のスライスの型リテラルを使った変数定義
	var ns []int // []はスライスを表す。要素がint型
	println("ns=", ns) // ns= [0/0]0x0
	// mapの型リテラルを使った変数定義
	var m map[string]int
	println("m=", m) //m= 0x0


	// ** 構造体リテラル */ フィールドを指定して初期化（構造体リテラル）
	// 構造体リテラルの例
	// 変数定義の文法  =>  var 変数名 型
	p := struct {
		//型リテラル
		name string
		age  int
	}{
		// フィールドの値
		name: "Gopher",
		age:  10,
	}
	type Vertex struct {
		X, Y float64
	}

	// println(p) // error. illegal types for operand: print struct { name string; age int }
	// フィールドにアクセスする例
	p.age++ // p.age = p.age + 1と同じ
	println(p.name, p.age) // Gopher 11

	// ** 配列 */ 同じ型のデータを集めて並べたデータ構造
	// 要素の型はすべて同じ
	// 要素数が違えば別の型
	// 要素数は変更できない
	// 型は型リテラルで記述することが多い
	// 型と要素数がセット
	// いろいろな初期化方法
	// ゼロ値で初期化
	var ns1 [5]int
	// 配列リテラルで初期化
	var ns2 = [5]int{10, 20, 30, 40, 50}
	// 要素数を値から推論
	ns3 := [...]int{10, 20, 30, 40, 50}
	// 5番目が50、10番目が100で他が0の要素数11の配列
	ns4 := [...]int{5: 50, 10: 100}

	fmt.Println(ns1, ns2, ns3, ns4) //[0 0 0 0 0] [10 20 30 40 50] [10 20 30 40 50] [0 0 0 0 0 50 0 0 0 0 100]

	nns := [...]int{10, 20, 30, 40, 50}
	// 要素にアクセス
	println(nns[3]) // 添字は変数でもよい 40
	// 長さ
	println(len(nns))// 5
	// スライス演算
	fmt.Println(nns[1:2]) //[20]


	// ** スライス */ 配列の一部を切り出したデータ構造
	// スライスと配列の関係 ・・ スライスはベースとなる配列が存在している
	s_ns := nns[1:3] //[20 30]
	fmt.Println(s_ns)

	// ゼロ値はnil
	var s_ns1 []int
	// 長さと容量を指定して初期化
	// 各要素はゼロ値で初期化される
	s_ns1 = make([]int, 3, 10)// make(型リテラル,要素の長さ, 容量) => var array [10]int,  ns := array[0:3] // または array[:3] と同じような処理
	// スライスリテラルで初期化
	// 要素数は指定しなくてよい
	// 自動で配列は作られる
	var s_ns2 = []int{10, 20, 30, 40, 50} // => var array2 = [...]int{10, 20, 30, 40, 50} , ms := array2[0:5] // または array[:] と同じような処理

	// 5番目が50、10番目が100で他が0の要素数11のスライス
	s_ns3 := []int{5: 50, 10: 100}
	fmt.Println(s_ns1, s_ns2, s_ns3)//[0 0 0] [10 20 30 40 50] [0 0 0 0 0 50 0 0 0 0 100]

	//スライス操作
	ns5 := []int{10, 20, 30, 40, 50}
	// 要素にアクセス
	println(ns5[3]) // 40
	// 長さ
	println(len(ns5))// 5
	// 要素の追加
	// 容量が足りない場合は背後の配列が再確保される
	ns5 = append(ns5, 60, 70)// 挙動。容量が足りる場合 => 1.新しい要素をコピーするlenを更新する. 2.lenを更新する

	println(len(ns5))// 7
	// 容量
	println(cap(ns5)) //10

	//配列・スライスへのスライス演算
	ns6 := []int{10, 20, 30, 40, 50}
	n6, m6 := 2, 4

	// n番目以降のスライスを取得する
	fmt.Println(ns6[n6:]) // [30 40 50]

	// 先頭からm-1番目までのスライスを取得する
	fmt.Println(ns6[:m6]) // [10 20 30 40]

	// capを指定する
	ms6 := ns6[:m6:m6]
	fmt.Println(ms6, cap(ms6)) // 4

	//Slice でappendを使ったテクニック
	// カット
	ns6 = append(ns6[:1], ns6[3:]...)
	fmt.Println(ns6) //[10 40 50]
	//削除
	ns6 = append(ns6[:2], ns6[3:]...)
	fmt.Println(ns6) //[10 40 50]
	// or
	//ns6 = ns6[:2+copy(ns6[2:], ns6[3:])] //コンパイルエラーになっていた

	// ** マップ */
	// キーと値をマッピングさせるデータ構造
	// 	キーと値の型を指定する
	// キーには「==」で比較できる型しかNG

	//マップの初期化
	// ゼロ値はnil
	var m2 map[string]int
	// makeで初期化
	m3 := make(map[string]int)
	// 容量を指定できる
	m4 := make(map[string]int, 10)
	// リテラルで初期化
	m5 := map[string]int{"x": 10, "y": 20}
	// 空の場合
	m7 := map[string]int{}
	fmt.Println(m2,m3,m4,m5,m7) //map[] map[] map[] map[x:10 y:20] map[]

	// マップ操作
	m8 := map[string]int{"x": 10, "y": 20}
	// キーを指定してアクセス
	println(m8["x"])
	// キーを指定して入力
	m8["z"] = 30
	// 存在を確認する
	n8, ok := m8["z"]
	println(n8, ok)
	// キーを指定して削除する
	delete(m8, "z")
	// 削除されていることを確認
	n8, ok = m8["z"] // ゼロ値とfalseを返す
	println(n, ok)

	//スライスの要素がスライスの場合（2次元スライス
	var n9 [][]int
	fmt.Println(n9) // []
	//マップの値がスライスの場合
	var n10 map[string][]int
	fmt.Println(n10) // map[]


	// ** ユーザー定義型 */ typeで名前を付けて新しい型を定義する
	//type 型名 基底型
	// 組み込み型を基にする
	type MyInt int
	// 他のパッケージの型を基にする
	type MyWriter io.Writer
	// 型リテラルを基にする
	type Person struct {
		Name string
	}

	// ** 型エイリアス（Go 1.9以上）*/
	// 型のエイリアスを定義できる
	// 完全に同じ型
	// キャスト不要
	type Applicant = http.Client
	//型名を出力する%Tが同じ元の型名を出す
	fmt.Printf("%T", Applicant{})//http.Client


	// ** 関数 */
	//関数の使い方
	fmt.Println(add(1, 2)) // 3
	//多値の受け取り方

	fmt.Println(swap2(1,2))// 2,1
	nswap1, mswap1 := 10, 20
	swap3(&nswap1, &mswap1)
	println(nswap1, mswap1)//20 10


	//関数はファーストクラスオブジェクト
	fs := make([]func() string,2) //string型を返す関数
	fs[0] = func() string { return "hoge" }
	fs[1] = func() string { return "fuga" }
	for _, f := range fs {
		fmt.Println(f())
	}
	//hoge
	//fuga

	// ** 値のコピー */
	p2 := struct{age int;name string}{age:10,name: "Gopher"}
	p3 := p2 // コピー
	p3.age = 20
	println(p3.age, p3.name)//20 Gopher
	println(p2.age, p2.name)//10 Gopher


	// **  ポインタ */ 変数の格納先を表す値
	//値で渡される型の値に対して破壊的な操作を加える際に利用する
	// 破壊的な操作 = 関数を出てもその影響が残る
	//内部でポインタが用いられているデータ型
	// コンポジット型の一部 スライス マップ チャネル
	var xp int
	fxp(&xp)
	println(xp)

	// ** メソッド */ ..レシーバと紐付けられた関数
	// 	データとそれに対する操作を紐付けるために用いる
	// ドットでメソッドにアクセスする

	// 100をHex型として代入
	var hex Hex = 100
	// Stringメソッドを呼び出す
	fmt.Println("hex",hex.String()) //64 ???? ポインタじゃないから？
	frece := hex.String
	fmt.Println(frece())//64

	// ** レシーバ */ ・・メソッドに関連付けられた変数
	// 	メソッド呼び出し時には通常の引数と同じような扱いになる > コピーが発生する
	// ポインタを用いることでレシーバへの変更を呼び出し元に伝えることができる > レシーバがポインタの場合もドットでアクセスする
	var receve T
	(&receve).f()//hi
	receve.f()//ここ2つは同じ意味

}

// ** 関数 */
//関数の定義
func add(x int, y int) int {
	return x + y
}
// 型をまとめて記述できる
//返り値を複数設定できる
func swap(x, y int) (int, int) {
	return y, x //カンマ区切りで返り値を返す
}

//名前付き戻り値
//明示しない場合は戻り値用の変数の値が返される
func swap2(x, y int) (x2, y2 int) {
	y2, x2 = x, y
	return
}
func swap3(n *int, m *int) {
	*n, *m = *m, *n
	return
}

//ポインタを用いた関数
//intのポインタ型
func fxp(xp *int) {
	*xp = 100//*でポインタのサス先に値を入れる
}


// メソッド
type Hex int
func (h Hex) String() string {
	return fmt.Sprintf("%x", int(h))
}

//レシーバ
type T int
func (t *T) f() { println("hi") }
