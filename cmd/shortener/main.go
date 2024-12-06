package main

import (
  "fmt"
  "io"
  "net/http"
  "math/rand"
)

var mapIdUrl map[string]string  // карта id - url

//------------------------------------------------------------------------------
// генерирует случайный id (строка 8 символов)
func generate_random_id() string {
  var chars = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

  b := make([]byte, 8)
  for i := range b {
    b[i] = chars[rand.Intn(len(chars))]
  }
  return string(b)
}

//------------------------------------------------------------------------------
func rootPage( w http.ResponseWriter, r *http.Request ) {

  url, err := io.ReadAll(r.Body)  // получаем url из тела запроса
  if err != nil { http.Error( w, err.Error(), http.StatusInternalServerError ); return }

  id := generate_random_id()  // генерируем новый случайный id
  mapIdUrl[id] = string(url)  // запоминаем пару id - url

  w.Header().Set("content-type","text/plain")
  w.WriteHeader(http.StatusCreated)

  fmt.Fprintf(w,"http://localhost:8080/%s",id)

} // func mainPage

//------------------------------------------------------------------------------
func idPage( w http.ResponseWriter, r *http.Request ) {

  id := r.PathValue("id")  // получаем параметр id из запроса

  url := mapIdUrl[id]
  if url == "" { BadRequest(w,r); return }  // если нет такого id

  http.Redirect( w, r, url, http.StatusTemporaryRedirect )

} // func mainPage

//------------------------------------------------------------------------------
func BadRequest( w http.ResponseWriter, r *http.Request ) {
  http.Error( w, "Bad Request", http.StatusBadRequest )
}

//------------------------------------------------------------------------------
func main() {

  mapIdUrl = make(map[string]string)

  mux := http.NewServeMux()
  mux.HandleFunc( "POST /{$}", rootPage )
  mux.HandleFunc( "GET /{id}", idPage )
  mux.HandleFunc( "/",         BadRequest )

  err := http.ListenAndServe( "localhost:8080", mux )
  if err != nil { panic(err) }

} // func main

//------------------------------------------------------------------------------
