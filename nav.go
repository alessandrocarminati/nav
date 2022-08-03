package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

func main() {
	var funcID int

	hostPtr := flag.String("host", "localost", "db host")
	portPtr := flag.Int("port", 5432, "db port")
	userPtr := flag.String("user", "test", "db user")
	pwdPtr := flag.String("password", "test", "db password")
	dbnamePtr := flag.String("db", "kernel", "db name")
	funcNamePtr := flag.String("function", "", "function ID/function name")

	flag.Parse()

	t := Connect_token{*hostPtr, *portPtr, *userPtr, *pwdPtr, *dbnamePtr}
	fmt.Println("DB Connection:", t)
	db := Connect_db(&t)

	//integer => we are looking for function_id
	funcID, err := strconv.Atoi(*funcNamePtr)
	if err != nil {
		var index int = 1
		entries, err := get_entries_by_name(db, *funcNamePtr)
		if err == nil {
			for _, entry := range entries {
				fmt.Printf("%d: %+v\n", index, entry)
				index++
			}
		} else {
			fmt.Println("Symbol not found!")
			os.Exit(1)
		}
		fmt.Print("Select symbol: ")
		fmt.Scanf("%d", &index)
		if index < 1 || index > len(entries) {
			fmt.Println("Wrong selection!")
			os.Exit(1)
		}
		funcID = entries[index-1].Sym_id
		fmt.Println(funcID)
	}
	Navigate(db, funcID, make(map[int]bool), make(map[int]Entry))
}
