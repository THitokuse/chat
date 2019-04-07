package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"
	"trace"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

//templは1つのテンプレートを表します
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

//ServeHTTPはHTTPリクエストを処理します
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
		fmt.Println(data["UserData"])
	}

	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	//フラグを解釈します
	flag.Parse()

	//Gomniauthのセットアップ
	gomniauth.SetSecurityKey("セキュリティーキー")
	gomniauth.WithProviders(
		facebook.New(os.Getenv("GO_CHAT_FACEBOOK_GOMNIAUTH_CLIENT_ID"), os.Getenv("GO_CHAT_FACEBOOK_GOMNIAUTH_SECRET_KEY"), "http://localhost:8080/auth/callback/facebook"),
		github.New(os.Getenv("GO_CHAT_GITHUB_GOMNIAUTH_CLIENT_ID"), os.Getenv("GO_CHAT_GITHUB_GOMNIAUTH_SECRET_KEY"), "http://localhost:8080/auth/callback/github"),
		google.New(os.Getenv("GO_CHAT_GOOGLE_GOMNIAUTH_CLIENT_ID"), os.Getenv("GO_CHAT_GOOGLE_GOMNIAUTH_SECRET_KEY"), "http://localhost:8080/auth/callback/google"),
	)

	r := newRoom(UseFileSystemAvatar)
	r.tracer = trace.New(os.Stdout)
	//ルート
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars"))))
	http.Handle("/upload", &templateHandler{filename: "upload.html"})

	//ユーザーログアウト
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	//チャットルームを開始します。
	go r.run()
	//webサーバを起動します。
	log.Println("webサーバーを開始します。ポート：", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

	http.HandleFunc("/uploader", uploaderHandler)
}
