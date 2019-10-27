package main

/*
TODO
	* テストコードを足す
	* main関数を分ける
	* 可読性
	* テスタビリティ
	* エラー処理（panicを使わない）
*/

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/stepupgo/stepupgo2-1/server"
)

func main() {
	server.Run()
}
