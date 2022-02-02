package main

import (
	"fmt"
	"github.com/a-shine/butter"
	"github.com/a-shine/butter/node"
	"github.com/a-shine/butter/retrieve"
	"github.com/a-shine/butter/store"
	"net/http"
	"os"
)

type WikiUser struct {
	node *node.Node
}

func (user *WikiUser) store(w http.ResponseWriter, r *http.Request) {
	data := r.FormValue("article")

	var keywords []string
	for _, keyword := range r.Form["keyword"] {
		keywords = append(keywords, keyword)
	}

	uuid := store.NaiveStore(user.node, keywords, data)
	fmt.Fprintf(w, string(uuid))
}

func (user *WikiUser) retrieve(w http.ResponseWriter, r *http.Request) {
	uuid := r.FormValue("uuid")

	data := retrieve.NaiveRetrieve(user.node, uuid)
	fmt.Fprintf(w, string(data))
}

func addEntry(w http.ResponseWriter, r *http.Request) {
	dat, err := os.ReadFile("pages/add.html")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(w, string(dat))
}

func findEntry(w http.ResponseWriter, r *http.Request) {
	dat, err := os.ReadFile("pages/find.html")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(w, string(dat))
}

func hello(w http.ResponseWriter, req *http.Request) {
	dat, err := os.ReadFile("pages/welcome.html")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(w, string(dat))
}

func dummyClientBehaviour(node *node.Node) {
	// do nothing
}

func main() {
	butterNode, _ := node.NewNode(0, 512, dummyClientBehaviour, false)
	fmt.Println("Node created with address:", butterNode.Address())

	go butter.Spawn(&butterNode, false)

	user := WikiUser{&butterNode}

	http.HandleFunc("/", hello)
	http.HandleFunc("/add", addEntry)
	http.HandleFunc("/find", findEntry)
	http.HandleFunc("/store", user.store)
	http.HandleFunc("/retrieve", user.retrieve)

	http.ListenAndServe(":8001", nil)
}
