package store

// Хранилище - База Данных

import (
  "context"
  "database/sql"

  _ "github.com/jackc/pgx/v5/stdlib"
)

// тип: БД-хранилище записей URL
type StoreDB struct {
  db *sql.DB  // БД
}

//// тип: запись в файле хранилища
//// {"uuid":"2","short_url":"edVPg3ks","original_url":"http://ya.ru"}
//type StoreDBRecord struct {
//  ShortURL    string `json:"short_url"`
//  OriginalURL string `json:"original_url"`
//}

// таблица URL в БД
const sql_CREATE_TABLE = `
CREATE TABLE IF NOT EXISTS url (
  id           serial PRIMARY KEY,
  short_url    varchar(40)  NOT NULL,
  original_url varchar(500) NOT NULL
)`

//------------------------------------------------------------------------------
// Создание нового объекта типа StoreDB
// - dsn - Data Source Name
func NewStoreDB( dsn string ) (*StoreDB, error) {

  // открываем БД
  db, err := sql.Open( "pgx", dsn )
  if err != nil { return nil, err }

  _, err = db.Exec( sql_CREATE_TABLE )  // создаём таблицу в БД если её ещё нет
  if err != nil {
    db.Close()
    return nil, err
  }

  return &StoreDB{ db: db }, nil // return StoreDB, err
} // func

//------------------------------------------------------------------------------
// Проверяет связь с хранилищем
func (sd *StoreDB) Ping(ctx context.Context) error {

  //ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
  //defer cancel()

  return sd.db.PingContext(ctx)
} // func

//------------------------------------------------------------------------------
// Считывает все данные из хранилища в mapURL
func (sd *StoreDB) Load(ctx context.Context) error {

  rows, err := sd.db.QueryContext(ctx,"SELECT short_url, original_url from url")
  if err != nil { return err }
  defer rows.Close()  // закрываем rows перед выходом из функции

  for rows.Next() {  // перебираем записи
    if err = ctx.Err(); err != nil { return err }  // проверяем прерывание контекста

    var shortURL, originalURL string
    err = rows.Scan( &shortURL, &originalURL )  // получаем значения полей записи
    if err != nil { return err }

    mapURL[shortURL] = originalURL  // запоминаем в Памяти пару shortURL - originalURL
  }
  err = rows.Err()  // ошибка была?

  return err
} // func

//------------------------------------------------------------------------------
// Сохраняет originalURL в хранилище, возвращает его shortURL
func (sd *StoreDB) SaveURL(ctx context.Context, originalURL string) (shortURL string, err error) {

  shortURL = generateNewShortURL()  // генерируем новый уникальный shortURL

  // добавляем запись в таблицу
  _, err = sd.db.ExecContext( ctx, "INSERT INTO url (short_url, original_url) VALUES ($1, $2)", shortURL, originalURL )
  if err != nil { return "", err }

  mapURL[shortURL] = originalURL  // запоминаем в Памяти пару shortURL - originalURL

  return shortURL, nil  // OK
} // func

//------------------------------------------------------------------------------
// Получает originalURL из хранилища по его shortURL
func (sd *StoreDB) GetURL(ctx context.Context, shortURL string) (originalURL string, err error) {

  return mapURL[shortURL], nil  // ищем нужный URL в памяти, если не нашли, то --> ""
} // func

//------------------------------------------------------------------------------
// Закрывает хранилище
func (sd *StoreDB) Close() error {
  return sd.db.Close()
} // func

//------------------------------------------------------------------------------

