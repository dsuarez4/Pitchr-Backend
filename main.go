package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"encoding/json"
	"net/http"
	"log"

	"io/ioutil"
)
// error response contains everything we need to use http.Error
type handlerError struct {
	Error   error
	Message string
	Code    int
}
//type customer struct {
//	customerId int `json:"customerId"`
//	Tab *tab
//	Name string
//
//}
//type drink struct {
//	drinkId int `json:"drinkId"`
//	Name string
//	Price float32
//}

//type tab struct {
//
//	tabId int `json:"tabId"`
//	tabName string `json:"tabName"`
//	//hashmap of all the customers
//	Customers map[*customer]int `json:"customers"`
//
//}
// custom type to handle errors and formatting
type handler func(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError)

func (fn handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// here we could do some prep work before calling the handler if we wanted to

	// call the actual handler
	response, err := fn(w, r)

	// check for errors
	if err != nil {
		log.Printf("ERROR: %v\n", err.Error)
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Message), err.Code)
		return
	}
	if response == nil {
		log.Printf("ERROR: response from method is nil\n")
		http.Error(w, "Internal server error. Check the logs.", http.StatusInternalServerError)
		return
	}

	// turn the response into JSON
	bytes, e := json.Marshal(response)
	if e != nil {
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
		return
	}

	// send the response and log
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
	log.Printf("%s %s %s %d", r.RemoteAddr, r.Method, r.URL, 200)
}
// array of type tabs
var tabs = make([]Tab, 0)




func test(w http.ResponseWriter, r *http.Request) {

	log.Println("Responding to /hello request")
	log.Println(r.UserAgent())

	vars := mux.Vars(r)

	// pulled from rest endpoint
	name := vars["name"]

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Hello:", name)

}

func parseTabRequest(r *http.Request) (Tab, *handlerError) {
	// the tab payload is in the request body
	data, e := ioutil.ReadAll(r.Body)

	if e != nil {
		return Tab{}, &handlerError{e, "Could not read request", http.StatusBadRequest}
	}

	var payload Tab
	e = json.Unmarshal(data, &payload)

	if e != nil {
		return Tab{}, &handlerError{e, "Could not parse JSON", http.StatusBadRequest}
	}

	return payload, nil
}
// 104
func addTab(res http.ResponseWriter, req *http.Request) {
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
	//return payload
}

func listTabs(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	return tabs, nil
}

var id = 0
func getNextId() int {
	id += 1
	return id
}

func main() {


	fmt.Println("Hello SysAdmin")
	// new Gorilla router
	router := mux.NewRouter()
	// Add handler functions for routes
	//
	//
	//
	router.HandleFunc("/hello/{name}", test).Methods("GET")
	//
	router.HandleFunc("/tabs/{tabName}", addTab_handler).Methods("POST")
	//router.HandleFunc("/tabs", handler(listTabs)).Methods("GET")
	//
	http.Handle("/", router)


	daTab := newTab("Name of Test Tab #1")
	tabs = append(tabs, *daTab)


	log.Fatal(http.ListenAndServe(":8080", router))

}