package main

func main() {
	//* 変数*/
	// 変数定義と代入が一緒
	var n int = 100
	println(n) // 100
	// 変数定義と代入が別
	var m int
	m = 100
	println(m) // 100
	// 型を省略(int型になる)
	var l = 100
	println(l) // 100

	// varを省略
	// 関数内のみでしかできない
	x := 100
	println(x) // 100
	// まとめる
	var (
		y = 100
		z = 200
	)
	println(y, z) // 100, 200

	text := "Hello World"
	println(text) // Hello World

	//** 定数*/
	//変数に代入される場合はデフォルトの型になる
	const n2 = 10000000000000000000 / 10000000000000000000
	var m2 = n2 // mはint型
	println(m2) // 1

	const text2 = "Hello World"
	println(text2) // Hello World

	// 	名前付き定数定義の右辺が省略できる
	// グループ化された名前付き定数で用いられる
	const (
		a = 1 + 2
		b
		c
	)
	println(a, b, c) //全部3

	//iota・・・連続した定数を作るための仕組み
	// グループ化された名前付き定数の定義で使われる
	// 0から始まり1ずつ加算される値として扱われる
	const (
		d = iota //省略できる（iota=1）
		e
	)
	const (
		f = 1 << iota //式の中でも使える 1 << 0
		g             //1 << 1
		h             //1<< 2
	)
	println(d, e, f, g, h) // 0 1 1 2 4

	//** 制御構文 */
	// 特徴
	// ()がいらない
	// {}の位置は下記の通り、{}は必須
	if x == 1 {
		println("xは1")
	} else if x == 2 {
		println("xは2")
	} else {
		println("xは1でも2でもない")
	}

	// 代入文を書く
	//aaaはifとelseのブロックで使える
	if aaa := f3(); aaa > 0 {
		println(aaa)
	} else {
		println(2 * aaa)
	}
	//println(aaa) //できない

	//条件分岐switch
	// breakがいらないå
	// 大量のif-elseをつなぐより見通しがよくなる
	switch a {
	case 1, 2: // 1または2のとき
		println("a is 1 or 2") //何もしないとbreakになる
	default:
		println("default")
	}

	//** caseに式が使える */
	// fallthroughを使う => 処理を次の節（caseやdefault）に進めます
	//"a is 1 \n a in not
	const a1 = 1
	switch {
	case a1 == 1:
		println("a is 1")
		fallthrough
	default:
		println("a in not 1")
	}

	//** for */
	//繰り返しはforしかない
	// 初期値;継続条件;更新文
	for i := 0; i <= 100; i = i + 1 {
	}

	// 継続条件のみ
	jj := 0
	for jj <= 100 {
		jj++
	}
	// 無限ループ
	for {
		break //breakを用いるとループから抜け出せる
	}
	// rangeを使った繰り返し
	for i, v := range []int{1, 2, 3} {
		println(i, v) //0 1
		// 1 2
		// 2 3
	}
	var i3 int
LOOP:
	for {
		println("i3=", i3) //0 \n 1
		switch {
		case i3%2 == 1:
			break LOOP
		}
		i3++
	}

}

func f3() int {
	return 3
}
