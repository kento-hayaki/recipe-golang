package utils

import (
	"bufio"
	"fmt"
	"os"
)

// ** golangで扱える共通関数を定義する。

// ** input周り
// 文字列を1行入力
func StrStdin() (stringInput string) {
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	stringInput = scanner.Text()

	return
}

// ** ファイル操作周り
// ファイルを開き一行ずつ配列に格納する
func FromFileToArray(filePath string) []string {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "File %s could not read: %v\n", filePath, err)
		os.Exit(1)
	}

	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if serr := scanner.Err(); serr != nil {
		fmt.Fprintf(os.Stderr, "File %s scan error: %v\n", filePath, err)
	}

	return lines
}

// ** 出力周り
func (a *App) println(args ...any) {
	fmt.Fprintln(a.output, args...)
}

func (a *App) print(args ...any) {
	fmt.Fprint(a.output, args...)
}

func (a *App) waitEnter(message string) error {
	buffer := make([]byte, 256)
	a.print(message)
	_, err := a.input.Read(buffer)
	if err != nil {
		return err
	}
	return nil
}
