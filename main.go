package main

import (
	trace "Chat/Trace"
	"Chat/auth"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
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
	var adr = flag.String("addr", ":8080", "App address.")
	flag.Parse()
	gomniauth.SetSecurityKey("zsh92@21#lk32reajh3h8dhsb73boud8qy")
	gomniauth.WithProviders(
		github.New("id", "security key",
			"http://localhost:8080/auth/callback/github"), //TODO add github key and id
		google.New("id", "security_key",
			"http://localhost:8080/auth/callback/google"), //TODO add google key and id
		facebook.New("id", "security_key",
			"httpL//localhost:8080/auth/callback/facebook"), //TODO add facebook key and id
	)
	r := NewRoom()
	r.tracer = trace.New(os.Stdout)

	http.Handle("/chat", auth.MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", auth.LoginHandler)
	http.Handle("/room", r)

	go r.run()
	log.Println("Running WWW server:", adr)
	if err := http.ListenAndServe(*adr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
