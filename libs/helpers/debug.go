package helpers

import (
	"encoding/json"
	"fmt"
)

func Dump(d interface{}) {
	fmt.Println("-----")
	b, _ := json.MarshalIndent(d, "", "    ")
	println(string(b))
	fmt.Println("-----")
}

func Splash(d interface{}) {
	fmt.Println("##### ##### #####")
	fmt.Println("##### ##### #####")
	b, _ := json.MarshalIndent(d, "", "    ")
	println(string(b))
	fmt.Println("##### ##### #####")
	fmt.Println("##### ##### #####")
}
