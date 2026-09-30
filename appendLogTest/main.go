package main

import (
	"appendLogTest/database"
	"fmt"
)

func main() {
	fp, err := database.LogCreate("test.log")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	database.LogAppend("hello", fp)
	database.LogAppend("this is my second entry", fp)
	database.LogAppend("third entry", fp)

	fmt.Println("Done")
}
