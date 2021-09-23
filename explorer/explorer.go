package explorer

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/mingi3442/mingicoin/blockchain"
)

const (
	port        string = ":4000"
	templateDir string = "explorer/templates/"
)

var templates *template.Template

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func home(rw http.ResponseWriter, r *http.Request) { //rw(ResponeWriter)에는 유저에게 보내고 싶은 데이터를 적는다
	data := homeData{"Home!", blockchain.GetBlockchain().AllBlock()}
	templates.ExecuteTemplate(rw, "home", data)
}
func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", nil)
	case "POST":
		r.ParseForm()
		data := r.Form.Get("blockData")
		blockchain.GetBlockchain().AddBlock(data)
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}
}
func Start(port int) {
	handler := http.NewServeMux()
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))     //handle하기 전 load하는 코드 / (standard package이용)
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml")) // Must func은 error발생시 error를 반환 / (templates variable 사용)
	handler.HandleFunc("/", home)
	handler.HandleFunc("/add", add)
	fmt.Printf("Listening on http://localhost:%d \n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler)) //첫 번째 인자는 주소값이고 두번째 인자는 handler이다
}
