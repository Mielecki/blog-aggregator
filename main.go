package main

import (
	"fmt"

	"github.com/Mielecki/blog-aggregator/internal/config"
)

func main() {
	config_file, err := config.Read()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(config_file)
	config_file.SetUser("lane")
	config_file, err = config.Read()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(config_file)
}
