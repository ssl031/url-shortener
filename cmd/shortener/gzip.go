package main

import (
  "io"
  "compress/gzip"
  "net/http"
  "strings"
)

// compressResponseWriter реализует интерфейс http.ResponseWriter и позволяет прозрачно для сервера
// сжимать передаваемые данные и выставлять правильные HTTP-заголовки
type compressResponseWriter struct {
  w  http.ResponseWriter  // оригинальный ResponseWriter
  zw *gzip.Writer         // Writer для сжатия данных
}

//------------------------------------------------------------------------------
func NewCompressResponseWriter(w http.ResponseWriter) *compressResponseWriter {
  return &compressResponseWriter{
    w  : w,                  // оригинальный ResponseWriter
    zw : gzip.NewWriter(w),  // Writer для сжатия данных и записи их в оригинальный ResponseWriter
  }
} // func

//------------------------------------------------------------------------------
func (c *compressResponseWriter) Header() http.Header {
  return c.w.Header()
} // func

//------------------------------------------------------------------------------
func (c *compressResponseWriter) Write(p []byte) (int, error) {
  return c.zw.Write(p)
} // func

//------------------------------------------------------------------------------
func (c *compressResponseWriter) WriteHeader(statusCode int) {
  if statusCode < 300 {
    c.w.Header().Set("Content-Encoding","gzip")
  }
  c.w.WriteHeader(statusCode)
} // func

//------------------------------------------------------------------------------
// Close закрывает gzip.Writer и досылает все данные из буфера
func (c *compressResponseWriter) Close() error {
  return c.zw.Close()
} // func

//------------------------------------------------------------------------------
// compressReadCloser реализует интерфейс io.ReadCloser и позволяет прозрачно для сервера
// декомпрессировать получаемые от клиента данные
type compressReadCloser struct {
  r  io.ReadCloser  // оригинальный ReadCloser
  zr *gzip.Reader   // Reader распакованных данных
}

//------------------------------------------------------------------------------
func NewCompressReadCloser(r io.ReadCloser) (*compressReadCloser, error) {

  zr, err := gzip.NewReader(r)  // Reader распакованных данных от оригинального ReadCloser
  if err != nil { return nil, err }

  return &compressReadCloser{
    r  : r,   // оригинальный ReadCloser
    zr : zr,  // Reader распакованных данных от оригинального ReadCloser
  }, nil  // return
} // func

//------------------------------------------------------------------------------
func (c compressReadCloser) Read(p []byte) (n int, err error) {
  return c.zr.Read(p)
} // func

//------------------------------------------------------------------------------
func (c *compressReadCloser) Close() error {
  if err := c.r.Close(); err != nil { return err }
  return c.zr.Close()
} // func


//------------------------------------------------------------------------------
func gzipMiddleware(h http.Handler) http.Handler {
  return http.HandlerFunc( func(w http.ResponseWriter, r *http.Request) {
    // по умолчанию устанавливаем оригинальный http.ResponseWriter как тот, который будем передавать следующей функции
    ow := w

    // если клиент умеет получать от сервера сжатые данные в формате gzip
    if strings.Contains( r.Header.Get("Accept-Encoding"), "gzip" ) {
      logger.Debug("compress(gzip) response")

      cw := NewCompressResponseWriter(w)  // оборачиваем оригинальный http.ResponseWriter новым с поддержкой сжатия
      ow = cw  // и меняем оригинальный http.ResponseWriter на новый

      defer cw.Close()  // не забываем отправить клиенту все сжатые данные после завершения middleware
    } // if

    // если клиент отправил серверу сжатые данные в формате gzip
    if strings.Contains( r.Header.Get("Content-Encoding"), "gzip" ) {
      // оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
      cr, err := NewCompressReadCloser(r.Body)
      if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
      } // if
      // меняем тело запроса на новое
      r.Body = cr
      defer cr.Close()
    } // if

    h.ServeHTTP( ow, r )  // вызываем оригинальный хендлер h
  } ) // return func
} // func

//------------------------------------------------------------------------------

