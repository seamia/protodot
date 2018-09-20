package main

import (
	"fmt"
	"os"
)

const key = "Katmanoooooooooooo"

func main() {

	fmt.Println("before [", os.Getenv(key), "]")
	if err := os.Setenv(key, "Yooooooooooooo"); err != nil {
		fmt.Println("ERROR: ", err)
	}
	fmt.Println("after [", os.Getenv(key), "]")
	
	fmt.Println(os.ExpandEnv("$USER lives in ${HOME}."))
	fmt.Println(os.ExpandEnv("$TMP ======= ${TEMP}."))
}