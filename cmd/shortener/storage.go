package main

import (
  "io"
  "encoding/json"
  "os"
)

// тип: хранилище записей ShortURL - OriginalURL
type Storage struct {
  file    *os.File       // файл хранилища
  encoder *json.Encoder  // json-encoder который записывает результат в файл file
}

// тип: запись в файле хранилища
// {"uuid":"2","short_url":"edVPg3ks","original_url":"http://ya.ru"}
type StorageRecord struct {
  ShortURL    string `json:"short_url"`
  OriginalURL string `json:"original_url"`
}

//------------------------------------------------------------------------------
// Создание нового объекта типа Storage
func NewStorage( filename string ) (*Storage, error) {

  // открываем файл хранилища
  file, err := os.OpenFile( filename, os.O_RDWR|os.O_CREATE, 0644)  // чтение-запись, создать файл если его нет, rw-r--r--
  if err != nil { return nil, err }

  _, err = file.Seek( 0, io.SeekEnd )  // переходим в конец файла - для последующего добавления новых записей
  if err != nil { return nil, err }

  return &Storage{
    file:    file,
    encoder: json.NewEncoder(file),
  }, nil // return Storage, err

} // func

//------------------------------------------------------------------------------
// Считывает все данные из хранилища в mapURL
func (st *Storage) Load( mapURL map[string]string ) error {

  // переходим в начало файла - для чтения всех записей
  _, err := st.file.Seek( 0, io.SeekStart )
  if err != nil { return err }

  // считываем записи, декодируем их из json и сохраняем в mapURL
  decoder := json.NewDecoder(st.file)
  for {
    r := StorageRecord{}        // подготовим пустой объект
    err = decoder.Decode( &r )  // декодируем запись из json
    if err == io.EOF { break }
    if err != nil { return err }

    mapURL[r.ShortURL] = r.OriginalURL  // сохраняем пару ShortURL - OriginalURL
  } // for

  // переходим в конец файла - для последующего добавления новых записей
  _, err = st.file.Seek( 0, io.SeekEnd )
  return err
} // func

//------------------------------------------------------------------------------
// Добавление записи в хранилище
func (st *Storage) WriteRecord( shortURL, originalURL string ) error {

  r := StorageRecord{
    ShortURL:    shortURL,
    OriginalURL: originalURL,
  }

  return st.encoder.Encode( &r )
} // func

//------------------------------------------------------------------------------
// Закрывает хранилище
func (st *Storage) Close() error {
  return st.file.Close()
} // func

//------------------------------------------------------------------------------

