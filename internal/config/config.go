package config

import "flag"

var ServerAddress string  // адрес + порт на котором запускается сервис      # localhost:8080
var ServerBaseURL string  // базовый адрес результирующего сокращённого URL  # http://localhost:8080

//------------------------------------------------------------------------------
// Получает параметры конфигурации
func Parse() {
  // параметры командной строки:        default
  flag.StringVar( &ServerAddress, "a", "localhost:8080",        "адрес + порт на котором запускается сервис" )
  flag.StringVar( &ServerBaseURL, "b", "http://localhost:8080", "базовый адрес результирующего сокращённого URL")

  // разбираем параметры командной строки
  flag.Parse()
}
