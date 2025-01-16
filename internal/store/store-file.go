package store

// Хранилище - Текстовый файл с записями в формате json

import (
  "context"
  "io"
  "encoding/json"
  "os"
)

// тип: файл-хранилище записей URL
type StoreFile struct {
  file    *os.File       // файл хранилища
  encoder *json.Encoder  // json-encoder который записывает результат в файл file
}

// тип: запись в файле хранилища
// {"uuid":"2","short_url":"edVPg3ks","original_url":"http://ya.ru"}
type StoreFileRecord struct {
  ShortURL    string `json:"short_url"`
  OriginalURL string `json:"original_url"`
}

//------------------------------------------------------------------------------
// Создание нового объекта типа StoreFile
func NewStoreFile( filename string ) (*StoreFile, error) {

  // открываем файл хранилища
  file, err := os.OpenFile( filename, os.O_RDWR|os.O_CREATE, 0644)  // чтение-запись, создать файл если его нет, rw-r--r--
  if err != nil { return nil, err }

  _, err = file.Seek( 0, io.SeekEnd )  // переходим в конец файла - для последующего добавления новых записей
  if err != nil { return nil, err }

  return &StoreFile{
    file:    file,
    encoder: json.NewEncoder(file),
  }, nil // return StoreFile, err
} // func

//------------------------------------------------------------------------------
// Проверяет связь с хранилищем
func (sf *StoreFile) Ping(ctx context.Context) error {

  //if sf.file == nil { return errors.New("file is nil") }  // ping ERROR

  return nil  // ping OK
} // func

//------------------------------------------------------------------------------
// Считывает все данные из хранилища в mapURL
func (sf *StoreFile) Load(ctx context.Context) error {

  // переходим в начало файла - для чтения всех записей
  _, err := sf.file.Seek( 0, io.SeekStart )
  if err != nil { return err }

  // считываем записи, декодируем их из json и сохраняем в mapURL
  decoder := json.NewDecoder(sf.file)
  for {
    if err = ctx.Err(); err != nil { return err }  // проверяем прерывание контекста

    r := StoreFileRecord{}        // подготовим пустой объект
    err = decoder.Decode( &r )  // декодируем запись из json
    if err == io.EOF { break }
    if err != nil { return err }

    mapURL[r.ShortURL] = r.OriginalURL  // запоминаем пару shortURL - originalURL
  } // for

  // переходим в конец файла - для последующего добавления новых записей
  _, err = sf.file.Seek( 0, io.SeekEnd )
  return err
} // func

//------------------------------------------------------------------------------
// Сохраняет originalURL в хранилище, возвращает его shortURL
func (sf *StoreFile) SaveURL(ctx context.Context, originalURL string) (shortURL string, err error) {

  shortURL = generateNewShortURL()  // генерируем новый уникальный shortURL

  r := StoreFileRecord{  // запись для сохранения
    ShortURL:    shortURL,
    OriginalURL: originalURL,
  }

  err = sf.encoder.Encode( &r )  // кодируем запись в json и выводим в файл
  if err != nil { return "", err }

  mapURL[shortURL] = originalURL  // запоминаем пару shortURL - originalURL

  return shortURL, nil  // OK
} // func

//------------------------------------------------------------------------------
// Получает originalURL из хранилища по его shortURL
func (sf *StoreFile) GetURL(ctx context.Context, shortURL string) (originalURL string, err error) {

  return mapURL[shortURL], nil  // ищем нужный URL в памяти, если не нашли, то --> ""
} // func

//------------------------------------------------------------------------------
// Закрывает хранилище
func (sf *StoreFile) Close() error {
  return sf.file.Close()
} // func

//------------------------------------------------------------------------------

