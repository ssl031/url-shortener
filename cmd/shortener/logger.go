package main

import (
  "net/http"
  "time"

  "go.uber.org/zap"
)

// реализация http.ResponseWriter для логера
type loggingResponseWriter struct {
  http.ResponseWriter  // оригинальный http.ResponseWriter
  status int  // статус ответа
  size   int  // размер ответа
}

var logger *zap.Logger = zap.NewNop()  // логер, по умолчанию = no-op-логер, который не выводит никаких сообщений

//------------------------------------------------------------------------------
// Инициализация логера logger
func loggerInit(level string) error {

  // преобразуем уровень логирования из текста в zap.AtomicLevel
  atomicLevel, err := zap.ParseAtomicLevel(level)
  if err != nil { return err }

  // создаём конфигурацию логера
  cfg := zap.NewDevelopmentConfig()
  cfg.Level = atomicLevel  // устанавливаем уровень логирования

  // создаём логер на основе конфигурации
  zl, err := cfg.Build()
  if err != nil { return err }

  logger = zl  // устанавливаем логер
  return nil
}

//------------------------------------------------------------------------------
func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
  // записываем ответ, используя оригинальный http.ResponseWriter
  size, err := lrw.ResponseWriter.Write(b)  // вызываем оригинальный http.ResponseWriter.Write()
  lrw.size += size  // сохраняем размер ответа (суммируем)
  return size, err
}

//------------------------------------------------------------------------------
func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
  lrw.ResponseWriter.WriteHeader(statusCode)  // вызываем оригинальный http.ResponseWriter.WriteHeader()
  lrw.status = statusCode // сохраняем статус статуса
}

//------------------------------------------------------------------------------
// добавляет логирование запросов и ответов
func withLogging(h http.Handler) http.Handler {
  return http.HandlerFunc( func(w http.ResponseWriter, r *http.Request) {
    start := time.Now()  // время начала обработки запроса

    lrw := loggingResponseWriter {
      ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
    } // 

    h.ServeHTTP( &lrw, r ) // вызываем оригинальный обработчик запроса

    logger.Info( "request",
      zap.String("method",r.Method),              // метод запроса
      zap.String("uri",r.RequestURI),             // URI запроса
      zap.Duration("duration",time.Since(start)), // продолжительность выполнения запроса
      zap.Int("status",lrw.status),  // статус ответа
      zap.Int("size",  lrw.size),    // размер ответа
    ) // logger

  } ) // return http.HandlerFunc func
} // func

//------------------------------------------------------------------------------

