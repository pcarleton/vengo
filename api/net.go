package api

import (
  "fmt"
  "net/url"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "strings"
  "reflect"
  "code.google.com/p/goauth2/oauth"
)

func ClientFromToken(tokenString string) *http.Client {
  token := &oauth.Token{
    AccessToken: tokenString,
  }
  t := &oauth.Transport{
                Token:     token,
                Transport: http.DefaultTransport,
  }
  return t.Client()
}

func StructToUrlValues(input interface{}) url.Values {
  params := make(url.Values)
  v := reflect.ValueOf(input)
  fmt.Println(v.Kind())
  t := reflect.TypeOf(input)
  for i := 0; i < v.NumField(); i++ {
    fieldValue := v.Field(i)
    if (fieldValue.Len() != 0) {
      // Get json name
      urlName := strings.Split(t.Field(i).Tag.Get("json"), ",")[0]
      params.Set(urlName, fieldValue.String())
    }
  }
  return params
}

func CheckResponse(res *http.Response) error {
        if res.StatusCode >= 200 && res.StatusCode <= 299 {
                return nil
        }
        slurp, err := ioutil.ReadAll(res.Body)
        if err == nil {
                jerr := new(errorReply)
                err = json.Unmarshal(slurp, jerr)
                if err == nil && jerr.Error != nil {
                        return jerr.Error
                }
        }
        return fmt.Errorf("vengo: got HTTP response code %d and error reading body: %v", res.StatusCode, err)
}

type errorReply struct {
        Error *Error `json:"error"`
}

type Error struct {
        Code    int    `json:"code"`
        Message string `json:"message"`
}

func (e *Error) Error() string {
        return fmt.Sprintf("vengo: Error %d: %s", e.Code, e.Message)
}
