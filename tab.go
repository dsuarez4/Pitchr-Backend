package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"log"
	"fmt"
)

type Tab struct {

	TabId int `json:"tabId"`
	TabName string `json:"tabName"`
	//hashmap of all the key:customers value:customerId
	Customers map[*customer]int `json:"customers"`
}

func newTab(tabName string) *Tab {
	return &Tab {
		TabId: getNextId(),
		TabName: tabName,
		Customers: make(map[*customer] int),
	}
}

func addTab_handler(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)
	tabName := vars["tabName"]

	//payload, e := parseTabRequest(req)
	//if e != nil {
	//	//FIX LATER
	//	return payload
	//}

	newtab := new(Tab)
	newtab.TabName = tabName
	newtab.TabId = getNextId()

	decoder := json.NewDecoder(req.Body)
	r_error := decoder.Decode(&newtab)
	if r_error != nil {
		log.Println(r_error.Error())
		http.Error(res, r_error.Error(), http.StatusInternalServerError)
		return
	}

	tabs = append(tabs, *newtab)
	outgoingJson, err := json.Marshal(newtab)

	if err != nil {
		log.Println(err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	fmt.Fprintln(res, "Hello:", string(outgoingJson))
	//
	//payload.TabId = getNextId()
	//
	//tabs = append(tabs, payload)
	//
}