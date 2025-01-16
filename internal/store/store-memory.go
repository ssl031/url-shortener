package store

// Хранилище - Память

import (
  "context"
)

// тип: память-хранилище записей URL
type StoreMemory struct {
  // здесь ничего нет потому что записи URL сохраняются в памяти (в mapURL) при всех типах хранилищ
}

//------------------------------------------------------------------------------
// Создание нового объекта типа StoreMemory
func NewStoreMemory() (*StoreMemory, error) {

  return &StoreMemory{}, nil
} // func

//------------------------------------------------------------------------------
// Проверяет связь с хранилищем
func (sm *StoreMemory) Ping(ctx context.Context) error {

  return nil  // ping OK
} // func

//------------------------------------------------------------------------------
// Считывает все данные из хранилища в mapURL
func (sm *StoreMemory) Load(ctx context.Context) error {

  // при перезагрузке в памяти ничего не сохраняется, поэтому считывать нечего

  return nil  // ok
} // func

//------------------------------------------------------------------------------
// Сохраняет originalURL в хранилище, возвращает его shortURL
func (sm *StoreMemory) SaveURL(ctx context.Context, originalURL string) (shortURL string, err error) {

  shortURL = generateNewShortURL()  // генерируем новый уникальный shortURL

  mapURL[shortURL] = originalURL  // запоминаем пару shortURL - originalURL

  return shortURL, nil  // OK
} // func

//------------------------------------------------------------------------------
// Получает originalURL из хранилища по его shortURL
func (sm *StoreMemory) GetURL(ctx context.Context, shortURL string) (originalURL string, err error) {

  return mapURL[shortURL], nil  // ищем нужный URL в памяти, если не нашли, то --> ""
} // func

//------------------------------------------------------------------------------
// Закрывает хранилище
func (sm *StoreMemory) Close() error {
  return nil  // OK
} // func

//------------------------------------------------------------------------------

