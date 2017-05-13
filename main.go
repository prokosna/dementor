package main

import (
	"fmt"
	"log"
	"os"

	"github.com/prokosna/dementor/lib"
)

func main() {
	err := dementor.InitConf()
	fmt.Println(err)
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
}
