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

//------------------------------------------------------------------------------
// Создание нового объекта типа StoreDB
// - dsn - Data Source Name
func NewStoreDB( dsn string ) (*StoreDB, error) {

  // открываем БД
  db, err := sql.Open( "pgx", dsn )
  if err != nil { return nil, err }

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

/*
  // переходим в начало файла - для чтения всех записей
  _, err := sf.file.Seek( 0, io.SeekStart )
  if err != nil { return err }

  // считываем записи, декодируем их из json и сохраняем в mapURL
  decoder := json.NewDecoder(sf.file)
  for {
    if err = ctx.Err(); err != nil { return err }  // проверяем прерывание контекста

    r := StoreDBRecord{}        // подготовим пустой объект
    err = decoder.Decode( &r )  // декодируем запись из json
    if err == io.EOF { break }
    if err != nil { return err }

    mapURL[r.ShortURL] = r.OriginalURL  // запоминаем пару shortURL - originalURL
  } // for

  // переходим в конец файла - для последующего добавления новых записей
  _, err = sf.file.Seek( 0, io.SeekEnd )
  return err
*/
  return nil
} // func

//------------------------------------------------------------------------------
// Сохраняет originalURL в хранилище, возвращает его shortURL
func (sd *StoreDB) SaveURL(ctx context.Context, originalURL string) (shortURL string, err error) {

  shortURL = generateNewShortURL()  // генерируем новый уникальный shortURL
/*
  r := StoreDBRecord{  // запись для сохранения
    ShortURL:    shortURL,
    OriginalURL: originalURL,
  }

  err = sf.encoder.Encode( &r )  // кодируем запись в json и выводим в файл
  if err != nil { return "", err }
*/
  mapURL[shortURL] = originalURL  // запоминаем пару shortURL - originalURL

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

