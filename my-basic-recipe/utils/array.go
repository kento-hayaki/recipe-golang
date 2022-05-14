// ** 配列操作周り

package utils

import (
	"fmt"
	"sort"
)

// 連想配列のintの部分で昇順ソート。10件表示
func GetResultOrderByMiss() {
	sort.Slice(db, func(i, j int) bool {
		return db[i].miss < db[j].miss
	})
	for i := 0; i < 10 && i < len(db); i++ {
		fmt.Println(db[i])
	}
}

// 連想配列のtime.Timeの部分で降順ソート。10件表示
func GetResultOrderByPlayDate() {
	sort.Slice(db, func(i, j int) bool {
		return db[j].playDate.Before(db[i].playDate)
	})
	for i := 0; i < 10 && i < len(db); i++ {
		fmt.Println(db[i])
	}
}
