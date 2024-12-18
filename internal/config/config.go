package config

import (
  "flag"
  "os"
)

var ServerAddress string  // адрес + порт на котором запускается сервис      # localhost:8080
var ServerBaseURL string  // базовый адрес результирующего сокращённого URL  # http://localhost:8080

//------------------------------------------------------------------------------
// Получает параметры конфигурации
func Get() {

  // получаем параметры из командной строки
  // параметры командной строки:        default
  flag.StringVar( &ServerAddress, "a", "localhost:8080",        "адрес + порт на котором запускается сервис" )
  flag.StringVar( &ServerBaseURL, "b", "http://localhost:8080", "базовый адрес результирующего сокращённого URL")
  flag.Parse()  // разбираем параметры командной строки

  // получаем параметры из ENV
  if sa := os.Getenv("SERVER_ADDRESS"); sa != "" { ServerAddress = sa }
  if bu := os.Getenv("BASE_URL");       bu != "" { ServerBaseURL = bu }

}
