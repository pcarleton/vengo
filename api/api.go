package api

import (
  "net/http"
  "encoding/json"
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

type ListPaymentsRequest struct {
  AccessToken string `json:"access_token,omitempty"`
  Limit string `json:"limit,omitempty"`
  Before string `json:"before,omitempty"` // ISO 8601 format
  After string `json:"after,omitempty"`
}

type ListPaymentsResponse struct {

}
func (s *Service) ListPayments(req *ListPaymentsRequest) ListPaymentsResponse {
  return ListPaymentsResponse{}
}

type MakePaymentResponse struct {
  Data ResponseData `json:"data,omitempty"`
}

type ResponseData struct {
  Balance string `json:"balance,omitempty"`
  Pmt PmtData `json:"payment,omitempty"`
}

type PmtData struct {
  Status string `json:"status,omitempty"`
}

func (s *Service) MakePayment(payment *Payment) (*MakePaymentResponse, error) {
  ret := new(MakePaymentResponse)
  if err := s.MakeRequest("payments", "POST", *payment, ret); err != nil {
    return nil, err
  } 
  return ret, nil
}


func (s *Service) MakeRequest(targetUrl, method string, req, response interface{}) error {
  params := StructToUrlValues(req)
  urls := s.baseUrl + targetUrl + "?" + params.Encode()
  request, _ := http.NewRequest(method, urls, nil)

  ctype := "application/x-www-form-urlencoded"
  request.Header.Set("Content-Type", ctype)
  res, err := s.client.Do(request)
  if err != nil {
    return err
  }
  defer res.Body.Close()
  if err := CheckResponse(res); err != nil {
    return err
  }
  if err := json.NewDecoder(res.Body).Decode(response); err != nil {
    return err
  }
  return nil
}
