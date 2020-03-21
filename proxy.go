package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type children struct {
	ID       string    `json:"id" db:"id"`
	ItemID   string    `json:"item_id" db:"item_id"`
	StopTime time.Time `json:"stop_time" db:"stop_time"`
}

type item struct {
	ID         string     `json:"id" db:"id"`
	Title      string     `json:"title" db:"title"`
	CategoryID string     `json:"category_id" db:"category_id"`
	Price      float32    `json:"price" db:"price"`
	CurrencyID string     `json:"currency_id" db:"currency_id"`
	StartTime  time.Time  `json:"start_time" db:"start_time"`
	StopTime   time.Time  `json:"stop_time" db:"stop_time"`
	Children   []children `json:"children"`
}

func getItemData(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	db, dbErr := sql.Open("mysql", "root:sniPer$3@/ml_proxy")

	if dbErr != nil {
		panic(dbErr.Error())
	}
	defer db.Close()

	stmtOut, stmtErr := db.Prepare("SELECT * FROM item WHERE id = ?")

	if stmtErr != nil {
		panic(stmtErr.Error())
	}
	defer stmtOut.Close()

	//var id sql.NullString
	var mlResponse item
	var mlChildrenResponse []children

	queryErr := stmtOut.QueryRow(params["id"]).Scan()

	if queryErr == sql.ErrNoRows {
		response1, err1 := http.Get(fmt.Sprintf("https://api.mercadolibre.com/items/%s", params["id"]))
		response2, err2 := http.Get(fmt.Sprintf("https://api.mercadolibre.com/items/%s/children", params["id"]))

		if err1 != nil {
			fmt.Printf("Http request fails with error %s", err1)
		} else if err2 != nil {
			fmt.Printf("Http request fails with error %s", err2)
		} else {
			data1, _ := ioutil.ReadAll(response1.Body)
			data2, _ := ioutil.ReadAll(response2.Body)

			json.Unmarshal(data2, &mlChildrenResponse)
			json.Unmarshal(data1, &mlResponse)

			stmtIn1, _ := db.Prepare("INSERT INTO item (id, category_id, title, price, currency_id, start_time, stop_time) VALUES (?,?,?,?,?,?,?)")
			_, insErr := stmtIn1.Exec(mlResponse.ID, mlResponse.CategoryID, mlResponse.Title, mlResponse.Price, mlResponse.CurrencyID, mlResponse.StartTime, mlResponse.StopTime)
			
			if insErr != nil {
				panic(insErr.Error())
			}
			fmt.Printf("INSERT ITEM TO CACHE DB\n")

			for _, c := range mlChildrenResponse {
				c.ItemID = mlResponse.ID
				stmtIn2, _ := db.Prepare("INSERT INTO child (id, item_id, stop_time) VALUES (?,?,?)")
				stmtIn2.Exec(c.ID, c.ItemID, c.StopTime)
				fmt.Printf("INSERT CHILD TO CACHE DB\n")
			}

			mlResponse.Children = mlChildrenResponse
			resp, _ := json.Marshal(mlResponse)

			w.Header().Set("Content-Type", "application/json")
			w.Write(resp)
		}
	}
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/items/{id}", getItemData).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	handleRequests()
}
