package main

import (
  "net/http"
  "time"

  "go.uber.org/zap"
)

// реализация http.ResponseWriter для логгера
type loggingResponseWriter struct {
  http.ResponseWriter  // оригинальный http.ResponseWriter
  status int  // статус ответа
  size   int  // размер ответа
}

var logger *zap.Logger  // логгер

//------------------------------------------------------------------------------
// Инициализация логгера logger
func loggerInit() {
  var err error

  logger, err = zap.NewDevelopment()  // создаём логгер
  if err != nil { panic(err) }

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
func withLogging(h http.HandlerFunc) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
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

  } // return func
} // func

//------------------------------------------------------------------------------

