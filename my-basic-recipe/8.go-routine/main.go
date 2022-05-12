// ** Concurrency is not Parallelism
// - 並行と並列は別ものである by RobPike
// 並行：Concurrency
// 並列：Parallelism
// - Concurrency
// 同時にいくつかの質の異なることを扱う
// - Parallelism
// 同時にいくつかの質の同じことを扱う

// ** 並列と並行の違い
// ■ 並列 Concurrency：同時にいくつかの質の異なることを扱う
// ■ 並行 Parallelism：同時にいくつかの質の同じことを扱う

// ** ゴールーチンとConcurrency
// - ゴールーチンでConcurrencyを実現
// 複数のゴールーチンで同時に複数のタスクをこなす
// 各ゴールーチンに役割を与えて分業する
// - 軽量なスレッドのようなもの
// LinuxやUnixのスレッドよりコストが低い
// 1つのスレッドの上で複数のゴールーチンが動く -> 複数のコアで動くとは限らない ＝ 真に並列に動くとは限らない
// - ゴールーチンの作り方
// goキーワードをつけて関数を呼び出す
// go f()

// ** 無名関数とゴールーチン
package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"runtime"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	go func() {
		fmt.Println("別のゴールーチン")
	}()
	fmt.Println("mainゴールーチン")
	// 出力
	// mainゴールーチン
	// 別のゴールーチン
	time.Sleep(50 * time.Millisecond) // sleepしないとすぐに終了する。main関数を抜けきる前に設定したルーチンの実行される

	// go routine 動作確認
	defer fmt.Println("main done")
	go func() {
		defer fmt.Println("goroutine1 done")
		time.Sleep(3 * time.Second)
	}()

	go func() {
		defer fmt.Println("goroutine2 done")
		time.Sleep(1 * time.Second)
	}()
	time.Sleep(5 * time.Second)
	// 出力結果
	// goroutine2 done
	// goroutine1 done
	// main done

	// ** チャネルでデータ競合を避ける
	// 競合の例
	// ゴールーチン間のデータ競合 −1−
	var v int
	go func() { // ゴールーチン-1
		time.Sleep(1 * time.Second)
		v = 100 //共有の変数を使う

	}()
	go func() { // ゴールーチン-2
		time.Sleep(1 * time.Second)
		fmt.Println(v) //処理順序が保証されない。そのため競合が起こり0が出力される
	}()
	time.Sleep(2 * time.Second)

	// ゴールーチン間のデータ競合 −2−
	n := 1
	go func() {
		for i := 2; i <= 5; i++ {
			fmt.Println(n, "*", i)
			n *= i //競合
			time.Sleep(100)
		}
	}()
	for i := 1; i <= 10; i++ {
		fmt.Println(n, "+", i)
		n += 1 //競合
		time.Sleep(100)
	}
	// 出力 下記の通りごちゃごちゃ 実行ごとに違う結果
	// 1 + 1
	// 1 * 2
	// 4 + 2
	// 4 * 3
	// 15 * 4
	// 60 + 3
	// 61 + 4
	// 62 * 5
	// 310 + 5
	// 311 + 6
	// 312 + 7
	// 313 + 8
	// 314 + 9
	// 315 + 10

	// ** データ競合の解決
	// - 問題点
	// どのゴールーチンが先にアクセスするか分からない
	// 値の変更や参照が競合する
	// - 解決方法
	// 1つの変数には1つのゴールーチンからアクセスする
	// チャネルを使ってゴールーチン間で通信をする
	// またはロックをとる（syncパッケージ）
	// メモリの共有による通信を行わない。
	// メモリを共有するのではなく、通信することでメモリを共有する

	// ** チャネル
	// ** チャネルの特徴
	// - 送受信できる型
	// チャネルを定義する際に型を指定する
	// - バッファ
	// チャネルにバッファを持たせることができる
	// 初期化時に指定できる
	// 指定しないと容量0となる
	// - 送受信時の処理のブロック
	// 送信時にチャネルのバッファが一杯だとブロック
	// 受信時にチャネル内が空だとブロック

	// ** チャネルの基本 -1-
	// 初期化
	// ch1 := make(chan int)     // make(chan int, 0)と同じ
	// ch2 := make(chan int, 10) // 容量を指定

	// // 送信
	// ch1 <- 10      // 受け取られるまでブロック
	// ch2 <- 10 + 20 // 一杯であればブロック

	// // 受信
	// n1 := <-ch1       //送信されるまでブロック
	// n2 := <-ch2 + 100 // 空であればブロック

	// ** チャネルの基本 −2−
	ch := make(chan int) // 容量0
	go func() {          // ゴールーチン-1
		ch <- 100
	}()
	go func() { // ゴールーチン-2
		v := <-ch
		fmt.Println(v) // 100
	}()
	time.Sleep(2 * time.Second)

	// ** 複数のチャネルから同時に受信 */・・
	// select-case -1- : select-caseを用いる
	// select-case -2-
	ch1 := make(chan int)
	ch2 := make(chan string)
	go func() { ch1 <- 100 }()
	go func() { ch2 <- "hi" }()

	select {
	case v1 := <-ch1:
		fmt.Println("v1 = ", v1)
	case v2 := <-ch2:
		fmt.Println("v2 = ", v2)
	}
	// 出力 // hi or 100  //後に代入されたチャネルの方だけ実行される？

	// ** nilチャネル
	ch3 := make(chan int)
	var ch4 chan string //ゼロ値はnil
	go func() { ch3 <- 100 }()
	go func() { ch4 <- "hi" }()

	select {
	case v1 := <-ch3:
		fmt.Println("v1 = ", v1)
	case v2 := <-ch4: // nilの場合は無視
		fmt.Println("v2 = ", v2)
	}

	// ** ファーストクラスオブジェクト
	// - チャネルはファーストクラスオブジェクト
	// 変数へ代入可能
	// 引数に渡す
	// 戻り値で返す
	// チャネルのチャネル //chan chan int など
	// - timeパッケージ
	// http://golang.org/pkg/time/#After
	// 5分間待つ
	//<-time.After(0 * time.Minute) //5分たったら現在時刻が送られてくるチャネルを返す

	// ** チャネルを引数や戻り値にする
	ch5 := makeCh()
	go func() { ch5 <- 100 }()
	fmt.Println(recvCh(ch5)) //100

	// ** 双方向チャネル */
	ch6 := makeCh2()
	go func() { ch6 <- 100 }()
	fmt.Println(recvCh2(ch6)) //200

	// ** 単方向チャネル */
	ch7 := makeCh3()
	go func(ch7 chan<- int) { ch7 <- 100 }(ch7) //送信専用のチャネル
	fmt.Println(recvCh3(ch7))                   // 100

	// ** Concurrencyの実現
	// - 複数のゴールーチンで分業する
	// タスクの種類によってゴールーチンを作る
	// Concurrencyを実現
	// - チャネルでやりとりする
	// ゴールーチン間はチャネルで値を共有する
	// 複雑すぎる場合はロックを使うことも
	// - for-selectパターン
	// ゴールーチンごとに無限ループを作る
	// メインのゴールーチンはselectで結果を受信
	// ** 標準入力から受け取った文字列を出力するコード
	// ch8 := input(os.Stdin)
	// for {
	// 	fmt.Print(">")
	// 	fmt.Println(<-ch8)
	// }

	// ** チャネル以外でデータ競合を避ける */
	// ** syncパッケージ
	// - チャネル以外を使う理由
	// チャネルだけを使っているとコードが難解になる場合がある
	// 複数のチャネルが登場したり
	// 競合を防ぎたいデータが複数ある場合
	// - syncパッケージ
	// データの競合を防ぐロックなどを提供するパッケージ
	// sync/atomicではアトミックな演算をするための型などを提供

	// ** ロック
	// - sync.Mutex
	//Lockメソッドを呼ぶとUnlockメソッドが呼ばれるまでLockメソッドの呼び出しでブロックする
	var m sync.Mutex // ゼロ値で使える
	m.Lock()
	go func() {
		time.Sleep(3 * time.Second)
		m.Unlock()
		fmt.Println("unlock 1")
	}()
	m.Lock() // ここでブロック
	m.Unlock()
	fmt.Println("unlock 2")
	// 出力
	// unlock 1
	// unlock 2

	// ** 書き込み・読み込みロック
	// - sync.RWMutex
	// Mutexに読み込み用のRLockとRUnlockが入ったもの
	//   R | W
	// R ○ | ☓
	// W ☓ | ☓
	var m1 sync.RWMutex
	m1.RLock()
	go func() {
		time.Sleep(3 * time.Second)
		m1.RUnlock()
		fmt.Println("unlock 1")
	}()
	m1.RLock() //読み込みロックだけではブロックしない
	m1.RUnlock()
	fmt.Println("unlock 2")
	// 出力
	// unlock 1

	//** 複数のゴールーチンの待機
	// - sync.WaitGroupを使う
	// Addメソッドに渡した数の合計の回数だけDoneメソッドを呼ぶ
	// Waitメソッドで処理をブロックして待機する
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(2 * time.Second)
		fmt.Println("done1")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(100 * time.Millisecond)
		fmt.Println("done2")
	}()

	wg.Wait()
	fmt.Println("done all")
	// 出力
	// done2
	// done1
	// done all

	// ** エラーを返すゴールーチンの待機・・golang.org/x/sync/errgroupを使う
	// 失敗した場合に最初のエラーが取得できる
	// WithContextを使うと1つでもエラーを起こすとキャンセルされる
	var eg errgroup.Group
	eg.Go(func() error { return errors.New("Error") })
	eg.Go(func() error { return errors.New("Error") })
	if err := eg.Wait(); err != nil {
		fmt.Println("エラーを返すゴールーチンの待機") // エラーを返すゴールーチンの待機
	}

	// ** 1度しか実行しない関数・・sync.Onceを使う
	// 1回以上Doメソッドを呼んでも意味がない
	// 複数のゴールーチンから1回しか呼ばないようにするために利用する
	var once sync.Once
	once.Do(f)
	once.Do(f) //2回目は実行されない
	fmt.Println("done")

	// ** ゴルーチンとチャネルを深く理解する
	// ** ゴルーチンのスケジューラの挙動・・ゴルーチンが切り替わるタイミング
	// - チャネルへの読み書き
	// ブロックされる場合のみ
	// - システムコール
	// 全てではなく待ちが発生するもの
	// - time.Sleepの呼び出し
	// - メモリの割り当て
	// - runtime.Goschedの呼び出し

	// ** 並列度
	// - 並列に動かせるゴルーチンの数
	// runtime.GOMAXPROCSで設定が可能
	// 環境変数のGOMAXPROCSでも設定ができる
	// runtime.NumCPUで論理CPUの数が返ってくる
	fmt.Println(runtime.NumCPU) //0x1006600
	// デフォルトはruntime.NumCPUの数
	// - 並列度が1の場合
	// 並列に動かないだけでうまく使えば有効
	// Google App Engineの第1世代の場合は並列度が1
	// 処理がブロックされるタイミングでうまく並行処理してやる
	// DBへのアクセスなど

	// ** チャネルのclose
	// closeの挙動
	// closeは送信側が行う
	// 同じチャネルは2度閉じれない
	// panicが起こる
	// 閉じられたチャネルには送信できない
	// panicが起こる
	// 受信するとゼロ値とfalseが返ってくる
	// closeを使ったブロードキャスト
	// 複数の受信箇所に一気にブロードキャストしたい
	// closeした瞬間に受信場所にゼロ値が送られる
	// 処理の終了を伝えるのに使われる
	// https://qiita.com/tenntenn/items/dd6041d630af7feeec52

	// ** コンテキスト */ ・・contextインタフェース
	// ゴールーチンをまたいだキャンセル処理
	// ゴールーチンをまたいで値を共有する
	//処理の締め切り・キャンセル信号・API境界やプロセス間を横断する必要のあるリクエストスコープな値を伝達させる
	// contextが役に立つのは、一つの処理が複数のゴールーチンをまたいで行われる場合

	// ** コンテキストとキャンセル処理 */・・ゴルーチンをまたいだ処理のキャンセルに使う
	// cancel関数が呼ばれるか親のDoneチャネルがクローズされるとDoneチャネルがクローズされる
	bc := context.Background()
	ctx, cancel := context.WithCancel(bc)
	defer cancel()
	for n := range gen(ctx) {
		fmt.Println(n)
		//1
		// 2
		// 3
		// 4
		// 5
		if n == 5 {
			break
		}
	}

	// ** タイムアウト ・・context.WithTimeoutを用いる
	bc2 := context.Background()
	t := 50 * time.Millisecond
	ctx2, cancel := context.WithTimeout(bc2, t)
	defer cancel()
	select {
	case <-time.After(1 * time.Second):
		fmt.Println("overslept")
	case <-ctx2.Done():
		fmt.Println(ctx2.Err()) //context deadline exceeded
	}

	// ** コンテキストに値を持たせる
	// WithValueで値を持たせる
	// 例：キャッシュを充てない

}

// ** チャネルを引数や戻り値にする
func makeCh() chan int {
	return make(chan int)
}
func recvCh(recv chan int) int {
	return <-recv
}

// ** 双方向チャネル */
func makeCh2() chan int {
	return make(chan int)
}
func recvCh2(recv chan int) int {
	go func() { recv <- 200 }() //間違った使い方ができる
	return <-recv
}

// ** 単方向チャネル */
func makeCh3() chan int {
	return make(chan int)
}
func recvCh3(recv <-chan int) int { //受信専用のチャネル
	return <-recv
}

// ** 標準入力から受け取った文字列を出力するコード
func input(r io.Reader) <-chan string {
	// TODO: チャネルを作る
	//var ch chan string
	ch := make(chan string)
	go func() {
		s := bufio.NewScanner(r)
		for s.Scan() {
			if s.Err() != nil {
				// non-EOF error.
				log.Fatal(s.Err())
			}
			// TODO: チャネルに読み込んだ文字列を送る
			ch <- s.Text()
		}
		// TODO: チャネルを閉じる
		defer close(ch)
	}()
	// TODO: チャネルを返す
	return ch
}

func f() { fmt.Println("Do!!") }

// Contextインタフェース
// ゴールーチンをまたいだキャンセル処理
// ゴールーチンをまたいで値を共有する
type Context interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
}

// ** コンテキストとキャンセル処理 */・・ゴルーチンをまたいだ処理のキャンセルに使う
// cancel関数が呼ばれるか親のDoneチャネルがクローズされるとDoneチャネルがクローズされる
func gen(ctx context.Context) <-chan int {
	dst := make(chan int)
	n := 1
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case dst <- n:
				n++
			}
		}
	}()
	return dst
}

//** コンテキストに値を持たせる
// WithValueで値を持たせる
// 例：キャッシュを充てない
type withoutCacheKey struct{}

func WithoutCache(ctx context.Context) context.Context {
	if IsIgnoredCache(ctx) {
		return ctx
	}
	return context.WithValue(ctx, withoutCacheKey{}, struct{}{})
}
func IsIgnoredCache(ctx context.Context) bool {
	return ctx.Value(withoutCacheKey{}) != nil
}
