package main

import (
  "context"
  "fmt"
  "io"
  "encoding/json"
  "net/http"
  "time"

  "github.com/go-chi/chi/v5"
  "go.uber.org/zap"

  "github.com/ssl031/url-shortener/internal/config"
  "github.com/ssl031/url-shortener/internal/store"
)

// app содержит в себе все зависимости и логику приложения
type app struct {
  store store.Store
}

//------------------------------------------------------------------------------
func newApp(s store.Store) *app {
  return &app{ store: s }
}

//------------------------------------------------------------------------------
func (a *app) ping( w http.ResponseWriter, r *http.Request ) {
  // GET /ping  --> 200 OK / 500 Internal Server Error

  ctx, cancel := context.WithTimeout( r.Context(), 2*time.Second )
  defer cancel()

  err :=  a.store.Ping(ctx)  // проверяем связь с хранилищем
  logger.Debug("store.ping", zap.Error(err))

  if err == nil {
    w.WriteHeader(http.StatusOK)                   // 200 OK
  } else {
    w.WriteHeader(http.StatusInternalServerError)  // 500 Internal Server Error
  }
} // func

//------------------------------------------------------------------------------
func (a *app) shorten( w http.ResponseWriter, r *http.Request ) {
  // POST / http://mail.ru  --> http://localhost:8080/uD2wgoIb

  ctx := r.Context()

  //fmt.Printf("shorten handler\n")

  url, err := io.ReadAll(r.Body)  // получаем url из тела запроса
  if err != nil { http.Error( w, err.Error(), http.StatusInternalServerError ); return }

  //// генерируем новый id
  //var id string
  //for {
  //  id = generateRandomID()  // генерируем id
  //  if _, exists := mapURL[id]; !exists { break }  // проверяем чтобы не было повтора id
  //} // for
  //
  //err = storage.WriteRecord( id, string(url) )  // сохраняем пару id - url в хранилище
  //if err != nil { http.Error( w, err.Error(), http.StatusInternalServerError ); return }
  //
  //mapURL[id] = string(url)  // запоминаем пару id - url

  // сохраняем URL в хранилище
  shortURL, err := a.store.SaveURL( ctx, string(url) )
  if err != nil {
    logger.Debug("cannot save url in store", zap.Error(err))
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  w.Header().Set("content-type","text/plain")
  w.WriteHeader(http.StatusCreated)

  fmt.Fprintf( w, "%s/%s", config.ServerBaseURL, shortURL)  // http://localhost:8080/uD2wgoIb

} // func

//------------------------------------------------------------------------------
type apiRequest struct {
  URL string `json:"url"`
}

type apiResponse struct {
  Result string `json:"result"`
}

func (a *app) apiShorten( w http.ResponseWriter, r *http.Request ) {
  // POST /api/shorten {"url":"http://mail.ru"}  --> {"result":"http://localhost:8080/uD2wgoIb"}

  ctx := r.Context()

  logger.Debug("got request", zap.String("method",r.Method), zap.String("path",r.URL.Path))

  var req apiRequest  // запрос (из json)
  // преобразуем json-запрос в объект req
  dec := json.NewDecoder(r.Body)
  if err := dec.Decode(&req); err != nil {
    logger.Debug("cannot decode request JSON body", zap.Error(err))
    //w.WriteHeader(http.StatusInternalServerError)
    http.Error( w, "cannot decode request JSON body", http.StatusBadRequest )
    return
  } // if

  //// генерируем новый id
  //var id string
  //for {
  //  id = generateRandomID()  // генерируем id
  //  if _, exists := mapURL[id]; !exists { break }  // проверяем чтобы не было повтора id
  //} // for
  //
  //err := storage.WriteRecord( id, string(req.URL) )  // сохраняем пару id - url в хранилище
  //if err != nil { http.Error( w, err.Error(), http.StatusInternalServerError ); return }
  //
  //mapURL[id] = req.URL  // запоминаем пару id - url

  // сохраняем URL в хранилище
  shortURL, err := a.store.SaveURL( ctx, req.URL )
  if err != nil {
    logger.Debug("cannot save url in store", zap.Error(err))
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  // готовим объект ответа
  resp := apiResponse { Result : config.ServerBaseURL+"/"+shortURL }  // http://localhost:8080/uD2wgoIb

  w.Header().Set("content-type","application/json")
  w.WriteHeader(http.StatusCreated)

  // преобразуем объект ответа в json
  enc := json.NewEncoder(w)
  if err := enc.Encode(resp); err != nil {
    logger.Debug("error encoding response", zap.Error(err))
    return
  }

  //logger.Debug("response", zap.String("json",???))
} // func

//------------------------------------------------------------------------------
func (a *app) redirect( w http.ResponseWriter, r *http.Request ) {
  // GET /uD2wgoIb  --> Redirect Location http://mail.ru

  ctx := r.Context()

  // получаем параметр shortURL из запроса  # GET /mYFl7FlK  --> shortURL="mYFl7FlK"
  shortURL := r.URL.Path
  if len(shortURL) > 0 && shortURL[0] == '/' { shortURL = shortURL[1:] }  // убираем первый символ /
  //shortURL := chi.URLParam(r,"shortURL")

  //fmt.Printf("redirect handler: shortURL=[%s]\n",shortURL)

  //url := mapURL[shortURL]

  originalURL, err :=  a.store.GetURL( ctx, shortURL )
  if err != nil {
    logger.Debug("cannot get url from store", zap.Error(err))
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  if originalURL == "" { BadRequest(w,r); return }  // если нет такого shortURL

  http.Redirect( w, r, originalURL, http.StatusTemporaryRedirect )  // перенаправляем на originalURL

} // func

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

  err = loggerInit( config.LogLevel )  // инициализируем logger
  if err != nil { panic(err) }
  defer logger.Sync()  // при завершении выведем оставшиеся сообщения из буфера

  // открываем хранилище
  var storage store.Store
  if        config.DatabaseDSN != "" {  // если указаны параметры подключения к БД, то для хранилища используем БД
    storage, err = store.NewStoreDB( config.DatabaseDSN )
    logger.Debug("store = DB", zap.String("dsm",config.DatabaseDSN))

  } else if config.FileStorage != "" {  // если указано имя файла-хранилища, то для хранилища используем Файл
    storage, err = store.NewStoreFile( config.FileStorage )
    logger.Debug("store = File", zap.String("filename",config.FileStorage))

  } else {                       // для хранилища используем Память
    storage, err = store.NewStoreMemory()
    logger.Debug("store = Memory")
  }
  if err != nil { logger.Fatal("Open Store",zap.Error(err)) }  // была ошибка открытия хранилища?
  defer storage.Close()  // при завершении - закроем хранилище

  app := newApp( storage )  // создаём экземпляр приложения

  // считываем все данные из хранилища в mapURL
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  err = app.store.Load(ctx)
  if err != nil { logger.Fatal("Store Load",zap.Error(err)) }
  //fmt.Println(mapURL)
  cancel()  // отменяем контекст

  rt := chi.NewRouter()

  rt.Use( withLogging, gzipMiddleware )  // middleware logger, gzip
  rt.Get ("/ping",        app.ping )       // GET  /ping                                 --> 200 OK / 500 Internal Server Error
  rt.Post("/",            app.shorten )    // POST /            http://mail.ru           --> http://localhost:8080/uD2wgoIb
  rt.Post("/api/shorten", app.apiShorten ) // POST /api/shorten {"url":"http://mail.ru"} --> {"result":"http://localhost:8080/uD2wgoIb"}
  rt.Get ("/{shortURL}",  app.redirect )   // GET  /uD2wgoIb                             --> Redirect Location http://mail.ru

  logger.Info("Starting server",zap.String("address",config.ServerAddress))
  defer logger.Info("STOP server")

  err = http.ListenAndServe( config.ServerAddress, rt )
  if err != nil { panic(err) }

} // func main

//------------------------------------------------------------------------------

