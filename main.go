package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}
func main() {
	var adr = flag.String("addr", ":8080", "Adres aplikacji internetowej.")
	flag.Parse()
	r := NewRoom()

	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)

	go r.run()
	log.Println("Uruchamianie serwera WWW pod adresem:", adr)
	if err := http.ListenAndServe(*adr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
