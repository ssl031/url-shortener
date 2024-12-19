package main

import (
  "fmt"
  "io"
  "net/http"
  "math/rand"

  "github.com/go-chi/chi/v5"
  "go.uber.org/zap"

  "github.com/ssl031/url-shortener/internal/config"
)

var mapURL = make( map[string]string )  // карта mapURL[id] -> url  # желательно сделать защиту этой map

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

  //fmt.Printf("rootPage\n")

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

  fmt.Fprintf( w, "%s/%s", config.ServerBaseURL, id)  // http://localhost:8080/uD2wgoIb

} // func mainPage

//------------------------------------------------------------------------------
func idPage( w http.ResponseWriter, r *http.Request ) {

  // получаем параметр id из запроса  # GET /mYFl7FlK  --> id="mYFl7FlK"
  id := r.URL.Path
  if len(id) > 0 && id[0] == '/' { id = id[1:] }  // убираем первый символ /
  //id := chi.URLParam(r,"id")

  //fmt.Printf("idPage id=[%s]\n",id)

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
  var err error

  config.Get()  // получаем конфигурацию
  //fmt.Printf("ServerAddress = [%s]\n",config.ServerAddress)
  //fmt.Printf("ServerBaseURL = [%s]\n",config.ServerBaseURL)

  loggerInit()  // инициализируем logger
  defer logger.Sync()  // при завершении выведем оставшиеся сообщения из буфера

  rt := chi.NewRouter()

  rt.Post("/",     withLogging(rootPage) )
  rt.Get ("/{id}", withLogging(idPage) )

  //mux := http.NewServeMux()
  //mux.HandleFunc( "POST /{$}", rootPage )
  //mux.HandleFunc( "GET /",     idPage )
  //mux.HandleFunc( "/",         BadRequest )

  logger.Info("Starting server",zap.String("address",config.ServerAddress))

  err = http.ListenAndServe( config.ServerAddress, rt )
  if err != nil { panic(err) }

} // func main

//------------------------------------------------------------------------------

