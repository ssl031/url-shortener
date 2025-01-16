package store

import (
  "context"
  "math/rand"
)

// интерфейс: абстрактное хранилище URL
type Store interface {

  // проверяет связь с хранилищем
  Ping(ctx context.Context) error

  // считывает все данные из хранилища в mapURL
  Load(ctx context.Context) error

  // сохраняет originalURL в хранилище, возвращает его shortURL
  SaveURL(ctx context.Context, originalURL string) (shortURL string, err error)

  // получает originalURL из хранилища по его shortURL
  GetURL(ctx context.Context, shortURL string) (originalURL string, err error)

  // закрывает хранилище
  Close() error
}


// запомненные URL: mapURL[shortURL] -> originalURL    # mapURL["J34oMyvD"] -> "http://mail.ru/"
var mapURL = make( map[string]string )


//------------------------------------------------------------------------------
// Генерирует случайный id (строка 8 символов)
func generateRandomID() string {
  var chars = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

  b := make([]byte, 8)
  for i := range b {
    b[i] = chars[rand.Intn(len(chars))]
  }
  return string(b)
}

//------------------------------------------------------------------------------
// Генерирует новый уникальный shortURL
func generateNewShortURL() string {

  var shortURL string
  for {
    shortURL = generateRandomID()  // генерируем случайный shortURL
    if _, exists := mapURL[shortURL]; !exists { break }  // если такого shortURL ещё не было, то возьмём его
  } // for

  return shortURL
}

//------------------------------------------------------------------------------

