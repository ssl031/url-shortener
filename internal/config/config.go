package config

import (
  "flag"
  "os"
)

var LogLevel      string  // уровень выводимых лог-записей                      # debug, info, warn, error, dpanic, panic, fatal
var ServerAddress string  // адрес + порт на котором запускается сервис         # localhost:8080
var ServerBaseURL string  // базовый адрес результирующего сокращённого URL     # http://localhost:8080
var FileStorage   string  // имя файла для сохранения записей ShortURL-OrigURL  # ""
var DatabaseDSN   string  // строка для подключения к БД Postgres               # ""

//------------------------------------------------------------------------------
// Получает параметры конфигурации
func Get() {

  // получаем параметры из командной строки
  // параметры командной строки:        default
  flag.StringVar( &LogLevel,      "ll", "info",                 "уровень выводимых лог-записей (debug,info,warn,error,dpanic,panic,fatal)")
  flag.StringVar( &ServerAddress, "a", "localhost:8080",        "адрес + порт на котором запускается сервис" )
  flag.StringVar( &ServerBaseURL, "b", "http://localhost:8080", "базовый адрес результирующего сокращённого URL")
  flag.StringVar( &FileStorage,   "f",  "",                     "имя файла для сохранения URL")
  flag.StringVar( &DatabaseDSN,   "d",  "",                     "строка Database-DSN для подключения к БД Postgres")
  flag.Parse()  // разбираем параметры командной строки

  // получаем параметры из ENV
  if ll := os.Getenv("LOG_LEVEL");         ll != "" { LogLevel      = ll }
  if sa := os.Getenv("SERVER_ADDRESS");    sa != "" { ServerAddress = sa }
  if bu := os.Getenv("BASE_URL");          bu != "" { ServerBaseURL = bu }
  if fs := os.Getenv("FILE_STORAGE_PATH"); fs != "" { FileStorage   = fs }
  if ds := os.Getenv("DATABASE_DSN");      ds != "" { DatabaseDSN   = ds }
}
