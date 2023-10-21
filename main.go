package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"orderAssembly/storage/sql"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(makeErr("order numbers not sent"))

		return
	}

	var orderNumbers []int

	ordersStr := strings.Split(os.Args[1], ",")

	for _, v := range ordersStr {
		num, err := strconv.Atoi(v)
		if err != nil {
			fmt.Println(err)
			return
		}
		orderNumbers = append(orderNumbers, num)
	}

	db, err := createDB()
	if err != nil {
		log.Fatal(err)
	}

	// Вызвать один раз. При повторных не сломается, но будет делать бессмысленную работу.
	err = db.CreateTables()
	if err != nil {
		fmt.Println(err)
		return
	}
	// Надо вызвать один раз, иначе будут дублирующиеся записи.
	err = db.CreateRecords()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = db.PrintAssemblyPage(orderNumbers)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func makeErr(strErr string) error {
	return fmt.Errorf("error: %s", strErr)
}

func createDB() (*sql.Database, error) {
	return sql.New("./storage/basic.db")
}
