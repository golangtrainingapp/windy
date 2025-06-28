package main

import (
	"fmt"
	"github.com/golangtrainingapp/windy"
)

func main() {
	resp, err := windy.GetWeather(53.1900, -112.2500, "mxJW8fEadecqILVj7RWBdhUfJ38Ou0Bv")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp)
}
