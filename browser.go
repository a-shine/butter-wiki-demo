package main

import (
	"fmt"
	"github.com/a-shine/butter"
	"github.com/a-shine/butter/node"
	"github.com/a-shine/cs347-cw/pcg"
	"html/template"
	"net/http"
	"strings"
)

type WikiUser struct {
	overlayInterface *pcg.Peer
}

func (user *WikiUser) store(w http.ResponseWriter, r *http.Request) {
	data := r.FormValue("article")
	uuid := pcg.Store(user.overlayInterface, data)
	fmt.Fprintf(w, uuid)
}

func (user *WikiUser) retrieve(w http.ResponseWriter, r *http.Request) {
	uuid := r.FormValue("uuid")
	data, err := pcg.NaiveRetrieve(user.overlayInterface, strings.TrimSpace(uuid))
	if err != nil {
		fmt.Fprintf(w, "Unable to find the information on the network")
	} else {
		fmt.Fprintf(w, string(data))
	}
}

func addEntry(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"pages/add.html",
		"pages/base.html",
	}
	temp, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Println(err)
	}
	temp.Execute(w, nil)
}

func findEntry(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"pages/find.html",
		"pages/base.html",
	}
	temp, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Println(err)
	}
	temp.Execute(w, nil)
}

func hello(w http.ResponseWriter, req *http.Request) {
	// Initialize a slice containing the paths to the two files. Note that the
	// home.page.tmpl file must be the *first* file in the slice.
	files := []string{
		"pages/welcome.html",
		"pages/base.html",
	}
	temp, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Println(err)
	}
	temp.Execute(w, nil)
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
