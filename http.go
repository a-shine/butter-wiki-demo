package main

import (
	"fmt"
	"github.com/a-shine/butter"
	"github.com/a-shine/butter/node"
	"github.com/a-shine/cs347-cw/pcg"
	"net/http"
	"os"
	"strings"
)

type WikiUser struct {
	overlayInterface *pcg.Peer
}

func (user *WikiUser) store(w http.ResponseWriter, r *http.Request) {
	data := r.FormValue("article")
	uuid := pcg.PCGStore(user.overlayInterface, data)
	fmt.Fprintf(w, uuid)
}

func (user *WikiUser) retrieve(w http.ResponseWriter, r *http.Request) {
	uuid := r.FormValue("uuid")
	data := pcg.NaiveRetrieve(user.overlayInterface, strings.TrimSpace(uuid))
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

func main() {
	butterNode, _ := node.NewNode(0, 512)
	fmt.Println("Node created with address:", butterNode.Address())

	overlay := pcg.NewPCG(butterNode, 512) // Creates a new overlay network
	pcg.AppendRetrieveBehaviour(overlay.Node())
	pcg.AppendGroupStoreBehaviour(overlay.Node())

	go butter.Spawn(&overlay, false)

	user := WikiUser{&overlay}

	http.HandleFunc("/", hello)
	http.HandleFunc("/add", addEntry)
	http.HandleFunc("/find", findEntry)
	http.HandleFunc("/store", user.store)
	http.HandleFunc("/retrieve", user.retrieve)

	http.ListenAndServe(":8000", nil)
}
