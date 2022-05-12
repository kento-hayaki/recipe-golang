package main

import (
	"fmt"
	"reflect"
)

func main() {
	// ** インタフェースと抽象化 */ ・・インタフェースはメソッドの集まり
	// メソッドのリストがインタフェースで規定しているものと
	// 一致する型はインタフェースを実装していることになる
	type Stringer interface {
		String() string
	}
	// インタフェースを実装していることになる
	var s Stringer = Hex(100)//Hex型にしている
	fmt.Println(s.String())//64
	fmt.Println(reflect.TypeOf(s))//main.Hex
	fmt.Printf("%T\n", s) // main.Hex

	// **empty interface
	// メソッドセットが空なインタフェース
	// つまりどの型の値も実装していることになる
	// JavaのObject型のような使い方ができる
	var v interface{}
	v = 100
	fmt.Println(v)//100
	v = "hoge"
	fmt.Println(v)//hoge
	v = 100

	//関数に型をつける
	var s1 fmt.Stringer = Func(func() string { return "test" })
	fmt.Println(s1)//test

	// **スライスとインタフェース
	// 実装していてもスライスは互換がない
	// コピーするには愚直にforで回すしかない
	ns := []int{1, 2, 3, 4}
	fmt.Println(ns)
	// できない
	//var vs []interface{} = ns //cannot use ns (variable of type []int) as []interface{} value in variable declarationcompiler

	//インタフェースの実装チェック・・コンパイル時に実装しているかチェックする
	// インタフェース型の変数に代入してみる
	var _ fmt.Stringer = Func(nil)//コンパイルできるから実装できているとみなせる
	
	// ** 型アサーション */ インタフェース.(型)
	// インタフェース型の値を任意の型にキャストする
	// 第2戻り値にキャストできるかどうかが返る
	// 第2戻り値を省略するとパニックが発生する
	var v1 interface{}
	v1 = 100
	n,ok := v1.(int)
	fmt.Println(n, ok)//100 true
	// s,ok := v1.(string)//: cannot assign string to s (type Stringer) in multiple assignment:
	// 				   //string does not implement Stringer (missing String method)
	// fmt.Println(s, ok)
	
	// ** 型スイッチ * /・・型によって処理をスイッチする
	// 代入文は省略可能
	switch vt := v.(type) {
	case int:
		fmt.Println(vt*2)//200が表示される
	case string:
		fmt.Println(vt+"hoge")
	default:
		fmt.Println("default")
	}

	//** io.Readerとio.Writer */ 入出力の抽象化
	// 	入出力を抽象化したioパッケージで提供される型
	// それぞれ1つのメソッドしか持たないので実装が楽
	// 入出力をうまく抽象化し、さまざまな型を透過的に扱える
	// ファイル、ネットワーク、メモリ etc…
	// パイプのように簡単に入出力を繋げられる
	type Reader interface {
		Read(p []byte) (n int, err error)
	}
	type Writer interface {
		Write(p []byte) (n int, err error)
	}
	F(I(100))//100 I
	F(B(false))//false B
	F(S("test"))//test S

	//** 埋め込みとフィールド
	f := Fuga{Hoge{N:100}} // Fuga{Hoge:Hoge{N:100}} でもOK
	// Hoge型のフィールドにアクセスできる
	fmt.Println(f.N) //100
	// 型名を指定してアクセスできる
	fmt.Println(f.Hoge.N) //100

	// ** 埋め込みの特徴 */ 型リテラルでなければ埋め込められる
	// typeで定義したものや組み込み型
	// インタフェースも埋め込められる
	// インタフェースの実装
	// 埋め込んだ値のメソッドもカウント
	var s2 Stringer
	h := Hex(100)
	s2 = h
	fmt.Println(s2.String())//64

	h2 := Hex2{h}
	s2 = h2
	fmt.Println(s2.String())//64

	h3 := Hex2{100}
	s2 = h3
	fmt.Println(s2.String())//64


}
//Stringerの実装
type Hex int
func (h Hex) String() string {
	return fmt.Sprintf("%x", int(h))
}

// Hex2もStringerを実装
type Hex2 struct{ Hex }

//関数に型をつける
//関数にインターフェスを実装
//関数にかたをつけた関数を定義している
type Stringer interface {
	String() string
}
type Func func() string
type Func2 func(x int) string
func (f Func) String() string { return f() }

//Stringer interfaceを実装する３つの型たち
type I int
func (i I) String() string {
	return "I"
}
type B bool
func (b B) String() string {
	return "B"
}
type S string
func (s S) String() string {
	return "S"
}
//受け取った値を上記の3つの具象型によって分岐する関数
func F(s Stringer) {
	switch x := s.(type) {
	case I:
		fmt.Println(int(x), "I")
	case B:
		fmt.Println(bool(x), "B")
	case S:
		fmt.Println(string(x), "S")
	}
}

//** 構造体の埋め込み */ ・・構造体に匿名フィールドを埋め込む機能
type Hoge struct {
	N int
}
// Fuga型にHoge型を埋め込む
type Fuga struct {
	Hoge // 名前のないフィールドになる
}

//** インタフェースと埋め込み */・・既存のインタフェースの振る舞いを変える
type Hoge2 interface{
	M()
	N()
}
type Fuga2 struct {Hoge}//ンタフェースを埋め込む
func (f Fuga2) M() { //Mの振る舞いを変える
	fmt.Println("Hi")
	f.Hoge.M() // 元のメソッドを呼ぶ
}
func HiHoge(h Hoge2) Hoge2 {
	return fuga{h} // 構造体作る
}

// ** インタフェースの埋め込み */・・インタフェースをインタフェースに埋め込む
// 複数のインタフェースを合成する
// 複雑なインタフェースが必要な場合
type Reader interface { Read(p []byte) (n int, err error) }
type Writer interface { Write(p []byte) (n int, err error) }
// ReaderとWriterを埋め込む
type ReadWriter interface {//ReadメソッドとWriteメソッドを実装する必要がある
    Reader
    Writer
}
