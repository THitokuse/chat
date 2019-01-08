package main

import (
  "github.com/stretchr/gomniauth"
  "github.com/stretchr/gomniauth/providers/facebook"
  "github.com/stretchr/gomniauth/providers/github"
  "github.com/stretchr/gomniauth/providers/google"
  "log"
  "net/http"
  "text/template"
  "path/filepath"
  "sync"
  "flag"
  // "os"
  // "trace"
)

//templは1つのテンプレートを表します
type templateHandler struct {
  once      sync.Once
  filename  string
  templ     *template.Template
}
//ServeHTTPはHTTPリクエストを処理します
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  t.once.Do(func() {
    t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
  })
  t.templ.Execute(w, r)
}

func main() {
  var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
  //フラグを解釈します
  flag.Parse()

  //Gomniauthのセットアップ
  gomniauth.SetSecurityKey("セキュリティーキー")
  gomniauth.WithProviders(
    facebook.New("クライアントID", "秘密の値", "http://localhost:8080/auth/callback/facebook"),
    github.New("クライアントID", "秘密の値", "http://localhost:8080/auth/callback/github"),
    google.New("クライアントID", "秘密の値", "http://localhost:8080/auth/callback/google"),
  )

  r := newRoom()
  // r.tracer = trace.New(os.Stdout)
  //ルート
  http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
  http.Handle("/login", &templateHandler{filename: "login.html"})
  http.HandleFunc("/auth/", loginHandler)
  http.Handle("/room", r)
  //チャットルームを開始します。
  go r.run()
  //webサーバを起動します。
  log.Println("webサーバーを開始します。ポート：", *addr)
  if err := http.ListenAndServe(*addr, nil); err != nil {
    log.Fatal("ListenAndServe:", err)
  }
}
