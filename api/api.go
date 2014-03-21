package api

import (
  "fmt"
  "net/http"
  "net/url"
  "bytes"
  "encoding/json"
  "io/ioutil"
  "unicode"
)

type Service struct {
  client *http.Client
}

func New(client *http.Client) *Service {
  return &Service{client}
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

func (c *MakePaymentCall) Do() (interface{}, error) {
  buf := new(bytes.Buffer)
  //err := json.NewEncoder(buf).Encode(c.payment)
  //if err != nil {
  //  return nil, err
  //}
  otherBody := `{"access_token":"EQcxUksCSBjPPxtjrHkgYBPuArZjF5JR","note":"testing!","amount":-1.00,"audience":"private"}`
  buf.WriteString(otherBody)
  fmt.Println(buf)
  urls := "https://api.venmo.com/v1/payments"
  params := make(url.Values)
  params.Set("access_token", c.payment.AccessToken)
  params.Set("phone", c.payment.Phone)
  params.Set("amount", c.payment.Amount)
  params.Set("note", c.payment.Note)
  ctype := "application/x-www-form-urlencoded"
  urls += "?" + params.Encode()
  req, _ := http.NewRequest("POST", urls, nil)
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

func UrlEncode(p *Payment) string {
  return ""
}

func Spacerize(s string) string {
  buf := new(bytes.Buffer)

  inARow := false
  for i, char := range s {
    lower := unicode.ToLower(char)
    if lower != char && i != 0 {
      if !inARow {
        buf.WriteRune('_')
      }
      inARow = true
    }  else {
      inARow = false
    }

    buf.WriteRune(lower)
  }
  return buf.String()
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
