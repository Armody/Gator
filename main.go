package main

import (
	"fmt"
	"log"

	"github.com/Armody/Gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg)

	err = cfg.SetUser("Armody")
	if err != nil {
		log.Fatal(err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg)
}
