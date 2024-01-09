package main

import (
	"fmt"

	"github.com/palmcivet7/go-blockchain/utils"
)

func main() {
	fmt.Println(utils.FindNeighbours("127.0.0.1", 5000, 0, 3, 5000, 5003))
}