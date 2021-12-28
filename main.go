package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Employee struct {
	Id int
	Name string
	City string
}

func dbConn() (db *sql.DB) {
	db, err := sql.Open("mysql", "root:@/golang")

	if err != nil {
		panic(err)
	}

	return db
}

func userCreate(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	name := r.FormValue("name")
	city := r.FormValue("city")

	insert, err := db.Prepare("INSERT INTO employee(name, city) VALUES(?,?)")

	if err != nil {
		panic(err)
	}

	insert.Exec(name, city)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	data := make(map[string]string)

	data["name"] = name
	data["city"] = city

	jsonResp, _ := json.Marshal(data)

	w.Write(jsonResp)

	defer db.Close()
}

func userLists(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	result, err := db.Query("SELECT * FROM employee")

	if err != nil {
		panic(err)
	}

	emp := Employee{}

	res := []Employee{}

	for result.Next() {
		var id int
		var name, city string

		err = result.Scan(&id, &name, &city)

		if err != nil {
			panic(err.Error())
		}

		emp.Id = id
		emp.Name = name
		emp.City = city

		res = append(res, emp)
	}

	jsonData, _ := json.Marshal(res)

	w.Write(jsonData)

	defer db.Close()
}


func main() {
	log.Println("Server starated on: http://localhost:8080")

	http.HandleFunc("/v1/users", userLists)
	http.HandleFunc("/v1/user", userCreate)
	http.ListenAndServe(":8080", nil)
}