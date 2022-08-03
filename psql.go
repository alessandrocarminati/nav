package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Connect_token struct {
	Host   string
	Port   int
	User   string
	Pass   string
	Dbname string
}

type Entry struct {
	Sym_id   int
	Symbol   string
	Exported bool
	Type     string
	Subsys   []string
	Fn       string
}

type Edge struct {
	Caller int
	Callee int
}

func Connect_db(t *Connect_token) *sql.DB {
	fmt.Println("connect")
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", (*t).Host, (*t).Port, (*t).User, (*t).Pass, (*t).Dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}
	fmt.Println("connected")
	return db
}

func get_entries_by_name(db *sql.DB, name string) ([]Entry, error) {
	var entries []Entry
	var tmp_entry, entry Entry
	var last_entry_id int = -1
	var s sql.NullString

	query := "select sym_id, symbol, exported, type, subsys, fn from (select * from symbols, kernel_file where symbols.fn_id=kernel_file.id) as dummy left outer join tags on dummy.fn_id=tags.fn_id where symbol=$1 order by sym_id"
	rows, err := db.Query(query, name)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&tmp_entry.Sym_id, &tmp_entry.Symbol, &tmp_entry.Exported, &tmp_entry.Type, &s, &tmp_entry.Fn); err != nil {
			fmt.Println(err)
			return entries, err
		}
		if last_entry_id != tmp_entry.Sym_id {
			if last_entry_id != -1 {
				entries = append(entries, entry)
			}
			entry.Sym_id = tmp_entry.Sym_id
			entry.Symbol = tmp_entry.Symbol
			entry.Exported = tmp_entry.Exported
			entry.Type = tmp_entry.Type
			entry.Fn = tmp_entry.Fn
			entry.Subsys = []string{}
		} else {
			if s.Valid {
				entry.Subsys = append(entry.Subsys, s.String)
			}
		}
		last_entry_id = tmp_entry.Sym_id
	}
	if last_entry_id != -1 {
		entries = append(entries, entry)
	}
	return entries, nil
}

func get_entry_by_id(db *sql.DB, symbol_id int, cache map[int]Entry) (Entry, error) {
	var e Entry
	var s sql.NullString

	//	fmt.Println("query")
	e, ok := cache[symbol_id]
	if !ok {
		query := "select sym_id, symbol, exported, type, subsys, fn from (select * from symbols, kernel_file where symbols.fn_id=kernel_file.id) as dummy left outer join tags on dummy.fn_id=tags.fn_id where sym_id=$1"
		rows, err := db.Query(query, symbol_id)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(&e.Sym_id, &e.Symbol, &e.Exported, &e.Type, &s, &e.Fn); err != nil {
				fmt.Println("this error hit3")
				fmt.Println(err)
				return e, err
			}
			//		fmt.Println(e)
			if s.Valid {
				e.Subsys = append(e.Subsys, s.String)
			}
		}
		if err = rows.Err(); err != nil {
			fmt.Println("this error hit")
			return e, err
		} else {
			cache[symbol_id] = e
		}
	}
	return e, nil
}

func get_successors_by_id(db *sql.DB, symbol_id int, cache map[int]Entry) ([]Entry, error) {
	var e Edge
	var res []Entry

	query := "select caller, callee from calls where caller =$1"
	rows, err := db.Query(query, symbol_id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&e.Caller, &e.Callee); err != nil {
			fmt.Println("this error hit1 ")
			return nil, err
		}
		successors, _ := get_entry_by_id(db, e.Callee, cache)
		res = append(res, successors)
	}
	if err = rows.Err(); err != nil {
		fmt.Println("this error hit2 ")
		return nil, err
	}
	return res, nil
}

func Not_in(list []int, v int) bool {

	for _, a := range list {
		if a == v {
			return false
		}
	}
	return true
}

func Navigate(db *sql.DB, symbol_id int, visited map[int]bool, cache map[int]Entry) {
	var rname, lname, lpname string

	visited[symbol_id] = true
	lname = ""
	entry, err := get_entry_by_id(db, symbol_id, cache)
	if err != nil {
		rname = "Unknown"
	}
	rname = entry.Symbol
	successors, err := get_successors_by_id(db, symbol_id, cache)
	if err == nil {
		for _, curr := range successors {
			entry, err := get_entry_by_id(db, curr.Sym_id, cache)
			if err != nil {
				lname = "Unknown"
			}
			lname = entry.Symbol
			if lpname != lname {
				fmt.Printf("%s -> %s\n", rname, lname)
			}
			lpname = lname
			if !visited[curr.Sym_id] {
				Navigate(db, curr.Sym_id, visited, cache)
			}
		}
	}
}
