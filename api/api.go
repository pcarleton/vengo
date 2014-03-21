package api

import (
  "fmt"
  "net/http"
  "net/url"
  "encoding/json"
  "io/ioutil"
  "strings"
  "reflect"
)

type Service struct {
  client *http.Client
  baseUrl string
}

func New(client *http.Client) *Service {
  return &Service{
    client: client,
    baseUrl: "https://api.venmo.com/v1/",
  }
}

type Payment struct {
  AccessToken string `json:"access_token,omitempty"`
  Phone string `json:"phone,omitempty"`
  Email string `json:"email,omitempty"`
  UserID string `json:"user_id,omitempty"`
  Note string `json:"note,omitempty"`
  Amount string `json:"amount,omitempty"` // Negative for a charge
  Audience string `json:"audience,omitempty"`
}

type MakePaymentCall struct {
  s *Service
  payment *Payment
}


func (s *Service) MakePayment(payment *Payment) *MakePaymentCall {
  return &MakePaymentCall{s: s, payment: payment}
}

type PaymentResponse struct {
  Data ResponseData `json:"data,omitempty"`
}

type ResponseData struct {
  Balance string `json:"balance,omitempty"`
  Pmt PmtData `json:"payment,omitempty"`
}

type PmtData struct {
  Status string `json:"status,omitempty"`
}

func (c *MakePaymentCall) Do() (*PaymentResponse, error) {
  params := StructToUrlValues(c.payment)
  urls := c.s.baseUrl + "payments?" + params.Encode()
  req, _ := http.NewRequest("POST", urls, nil)

  ctype := "application/x-www-form-urlencoded"
  req.Header.Set("Content-Type", ctype)
  res, err := c.s.client.Do(req)
  if err != nil {
    return nil, err
  }
  defer res.Body.Close()
  if err := CheckResponse(res); err != nil {
    return nil, err
  }
  ret := new(PaymentResponse)
  if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
    return nil, err
  }
  return ret, nil
}

func StructToUrlValues(p *Payment) url.Values {
  params := make(url.Values)
  v := reflect.ValueOf(*p)
  t := reflect.TypeOf(*p)
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
        return fmt.Sprintf("googleapi: Error %d: %s", e.Code, e.Message)
}
