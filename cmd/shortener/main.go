package main

import (
  "io"
  "net/http"
)

//------------------------------------------------------------------------------
func rootPage( w http.ResponseWriter, r *http.Request ) {

  url, err := io.ReadAll(r.Body)  // получаем url из тела запроса
  if err != nil { http.Error( w, err.Error(), http.StatusInternalServerError ); return }
  _ = url

  w.Header().Set("content-type","text/plain")
  w.WriteHeader(http.StatusCreated)

  w.Write([]byte("http://localhost:8080/EwHXdJfB"))

} // func mainPage

//------------------------------------------------------------------------------
func idPage( w http.ResponseWriter, r *http.Request ) {

  id := r.PathValue("id")  // получаем параметр id из запроса

  if id != "EwHXdJfB" { BadRequest(w,r); return }

  http.Redirect( w, r, "https://practicum.yandex.ru/", http.StatusTemporaryRedirect )

} // func mainPage

//------------------------------------------------------------------------------
func BadRequest( w http.ResponseWriter, r *http.Request ) {
  http.Error( w, "Bad Request", http.StatusBadRequest )
}

//------------------------------------------------------------------------------
func main() {

  mux := http.NewServeMux()
  mux.HandleFunc( "POST /{$}", rootPage )
  mux.HandleFunc( "GET /{id}", idPage )
  mux.HandleFunc( "/",         BadRequest )

  err := http.ListenAndServe( "localhost:8080", mux )
  if err != nil { panic(err) }

} // func main

//------------------------------------------------------------------------------
