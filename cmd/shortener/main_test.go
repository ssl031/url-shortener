package main

import (
  "io"
  "net/http"
  "net/http/httptest"
  "strings"
  "testing"
  "github.com/stretchr/testify/assert"
)

//var mapShortURL = make( map[string]string )  // карта mapShortURL[shortURL] -> targetURL

type urlPairT struct{ targetURL, shortURL string }  // пара targetURL - shortURL
var  urlPairs []urlPairT  // пары полученные при тестировании

//------------------------------------------------------------------------------
func TestRootPage(t *testing.T) {

  tests := []struct {
    name      string
    targetURL string
    wantCode  int
  }{
    { name: "create shortURL OK #1", targetURL: "https://practicum.yandex.ru/", wantCode: http.StatusCreated },
    { name: "create shortURL OK #2", targetURL: "http://2ip.ru/",               wantCode: http.StatusCreated },
  } // tests

  for _, test := range tests {
    t.Run(test.name, func(t *testing.T) {

      req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.targetURL))  // создаём запрос
      w := httptest.NewRecorder()  // создаём ResponseRecorder (implementation of http.ResponseWriter)

      rootPage(w, req)   // вызываем обработчик
      res := w.Result()  // получаем ответ (Response)

      assert.Equal(t, test.wantCode, res.StatusCode, "Код ответа")  // проверяем код ответа

      assert.Equal(t, "text/plain", res.Header.Get("Content-Type"), "Content-Type")  // проверяем Content-Type

      resBody, _ := io.ReadAll(res.Body)  // получаем тело ответа
      res.Body.Close()
      shortURL := string(resBody)         // в теле должен быть shortURL
      assert.Contains(t, shortURL, "http://localhost:8080/", "Полученная короткая ссылка")

      urlPairs = append( urlPairs, urlPairT{ test.targetURL, shortURL } )  // запоминаем целевой URL и его короткую ссылку

    }) // func, Run
  } // for tests

} // func

//------------------------------------------------------------------------------
func TestIdPage(t *testing.T) {
  // проверка полученных коротких ссылок

  // добавляем "плохую" короткую ссылку - для проверки BadRequest
  urlPairs = append( urlPairs, urlPairT{ "", "http://localhost:8080/BAD-SHORT-URL" } )

  for _, pair := range urlPairs {
    t.Run("get by short-url "+pair.shortURL, func(t *testing.T) {

      req := httptest.NewRequest( http.MethodGet, pair.shortURL, nil )  // создаём запрос
      req.SetPathValue( "id", pair.shortURL[22:] )  // установим параметр id  (делаем работу за ServeMux)  # 22 - длина строки http://localhost:8080/

      w := httptest.NewRecorder()  // создаём ResponseRecorder (implementation of http.ResponseWriter)

      idPage(w, req)     // вызываем обработчик
      res := w.Result()  // получаем ответ (Response)
      defer res.Body.Close()

      wantCode := http.StatusTemporaryRedirect  // ожидаем код ответа 307 Temporary Redirect
      if pair.targetURL == "" { wantCode = http.StatusBadRequest }  // если тестируем "плохую" короткую ссылку (targetURL==""), то ожидаем код ответа 400 Bad Request

      assert.Equal(t, wantCode, res.StatusCode, "Код ответа")  // проверяем код ответа

      assert.Equal(t, pair.targetURL, res.Header.Get("Location"), "Ссылка переадресации (Header.Location)")  // проверяем ссылку переадресации - должна быть равна targetURL

    }) // func, Run
  } // for tests

} // func TestUserViewHandler

//------------------------------------------------------------------------------

//------------------------------------------------------------------------------
