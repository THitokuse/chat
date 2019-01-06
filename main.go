package main

import (
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

  r := newRoom()
  // r.tracer = trace.New(os.Stdout)
  //ルート
  http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
  http.Handle("/room", r)
  //チャットルームを開始します。
  go r.run()
  //webサーバを起動します。
  log.Println("webサーバーを開始します。ポート：", *addr)
  if err := http.ListenAndServe(*addr, nil); err != nil {
    log.Fatal("ListenAndServe:", err)
  }
}
