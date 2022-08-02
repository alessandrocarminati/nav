package main

//import (
//	"fmt"
//	)

func main() {
	t := Connect_token{"localhost", 15432, "test", "test", "kernel"}
	db := Connect_db(&t)

	Navigate(db, 153735, make(map[int]bool), make(map[int]Entry))
}
