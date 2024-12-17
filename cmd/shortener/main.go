package main

import (
  "fmt"
  "io"
  "net/http"
  "math/rand"
)

var mapURL = make( map[string]string )  // карта mapURL[id] -> url
// желательно сделать защиту этой map

//------------------------------------------------------------------------------
// генерирует случайный id (строка 8 символов)
func generateRandomID() string {
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

  // генерируем новый id
  var id string
  for {
    id = generateRandomID()  // генерируем id
    if _, exists := mapURL[id]; !exists { break }  // проверяем чтобы не было повтора id
  } // for

  mapURL[id] = string(url)  // запоминаем пару id - url

  w.Header().Set("content-type","text/plain")
  w.WriteHeader(http.StatusCreated)

  fmt.Fprintf(w,"http://localhost:8080/%s",id)

} // func mainPage

//------------------------------------------------------------------------------
func idPage( w http.ResponseWriter, r *http.Request ) {

  // получаем параметр id из запроса  # GET /mYFl7FlK  --> id="mYFl7FlK"
  id := r.URL.Path
  if len(id) > 0 && id[0] == '/' { id = id[1:] }  // убираем первый символ /
  //fmt.Printf("id=[%s]\n",id)

  url := mapURL[id]
  if url == "" { BadRequest(w,r); return }  // если нет такого id

  http.Redirect( w, r, url, http.StatusTemporaryRedirect )

} // func mainPage

//------------------------------------------------------------------------------
func BadRequest( w http.ResponseWriter, r *http.Request ) {
  http.Error( w, "Bad Request", http.StatusBadRequest )
}

//------------------------------------------------------------------------------
func main() {

  mux := http.NewServeMux()
  mux.HandleFunc( "POST /{$}", rootPage )
  mux.HandleFunc( "GET /",     idPage )
  mux.HandleFunc( "/",         BadRequest )

  err := http.ListenAndServe( "localhost:8080", mux )
  if err != nil { panic(err) }

} // func main

//------------------------------------------------------------------------------
