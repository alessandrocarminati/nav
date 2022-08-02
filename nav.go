package main

//import (
//	"fmt"
//	)

func main() {
	t:=Connect_token{"dbs.hqhome163.com",5432,"alessandro","<password>","kernel"}
	db := Connect_db(&t)

	Navigate(db, 153735, make(map[int]bool), make(map[int]Entry))
}
