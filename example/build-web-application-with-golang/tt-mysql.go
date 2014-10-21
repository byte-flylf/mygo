// go-sql-driver/mysql驱动
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func main() {
	db, err := sql.Open("mysql", "root:mypassword@/test")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	time.Sleep(120)

	stmIns, err := db.Prepare("INSERT INTO squareNum VALUES(?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmIns.Close()

	for i := 0; i < 50; i++ {
		// 支持重连
		for {
			err = db.Ping()
			if err != nil {
				fmt.Println("ping fail", err)
				time.Sleep(5 * time.Second)
			} else {
				break
			}
		}
		_, err = stmIns.Exec(i, (i * i))
		if err != nil {
			fmt.Println("exec fail", err)
		}
		// 手动关闭mysqld，模拟网络中断
		//time.Sleep(2 * time.Second)
	}

	stmOut, err := db.Prepare("SELECT squareNumber FROM squareNum WHERE number = ?")
	if err != nil {
		panic(err.Error())
	}
	var squareNum int
	err = stmOut.QueryRow(13).Scan(&squareNum)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("The square number of 13 is: %d\n", squareNum)

	err = stmOut.QueryRow(1).Scan(&squareNum)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("The square number of 1 is: %d\n", squareNum)
}
