package explorer

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/ParkChangSun/nomadcoin/blockchain"
)

var templates *template.Template

const (
	templateDir string = "explorer/templates/"
)

var port string

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

// func home(rw http.ResponseWriter, r *http.Request) {
// 	data := homeData{"home", blockchain.GetBlockchain().AllBlocks()}
// 	templates.ExecuteTemplate(rw, "home", data)
// }

// func add(rw http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case "GET":
// 		templates.ExecuteTemplate(rw, "add", nil)
// 	case "POST":
// 		r.ParseForm()
// 		data := r.Form.Get("blockData")
// 		blockchain.GetBlockchain().AddBlock(data)
// 		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
// 	}
// }

func Start(aPort int) {
	port = fmt.Sprintf(":%d", aPort)
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))
	handler := http.NewServeMux()
	// handler.HandleFunc("/", home)
	// handler.HandleFunc("/add", add)
	log.Fatal(http.ListenAndServe(port, handler))
}
