package main

import (
	"fmt"

	"github.com/p-hti/heimdallr-server/internal/config"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)
}
