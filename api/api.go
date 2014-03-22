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

type MakePaymentRequest struct {
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
  Data []Payment `json:"data,omitempty"`
}

type MakePaymentResponse struct {
  Data ResponseData `json:"data,omitempty"`
}

type ResponseData struct {
  Balance string `json:"balance,omitempty"`
  Pmt Payment `json:"payment,omitempty"`
}

type Payment struct {
  Status string `json:"status,omitempty"`
  Target PaymentTarget `json:"target,omitempty"`
  DateCompleted string `json:"date_completed,omitempty"`
  Actor User `json:"actor,omitempty"`
  Note string `json:"note,omitempty"`
  Amount float32 `json:"amount,omitempty"`
  Action string `json:"action,omitempty"`
  DateCreated string `json:"date_created,omitempty"`
  ID string `json:"id,omitempty"`
}

type PaymentTarget struct {
  Phone string `json:"phone,omitempty"`
  Type string `json:"type,omitempty"`
  Email string `json:"email,omitempty"`
  User PaymentUser `json:"user,omitempty"`
}

type PaymentUser struct {
  Username string `json:"username,omitempty"`
  FirstName string `json:"first_name,omitempty"`
  LastName string `json:"last_name,omitempty"`
  DisplayName string `json:"display_name,omitempty"`
  About string `json:"about,omitempty"`
  ProfilePictureURL string `json:"profile_picture_url,omitempty"`
  ID string `json:"id,omitempty"`
  DateJoined string `json:"date_joined,omitempty"`
}

func (s *Service) ListPayments(req *ListPaymentsRequest) (*ListPaymentsResponse, error) {
  ret := new(ListPaymentsResponse)
  if err := s.makeRequest("payments", "GET", *req, ret); err != nil {
    return nil, err
  }
  return ret, nil
}


func (s *Service) MakePayment(payment *MakePaymentRequest) (*MakePaymentResponse, error) {
  ret := new(MakePaymentResponse)
  if err := s.makeRequest("payments", "POST", *payment, ret); err != nil {
    return nil, err
  }
  return ret, nil
}


func (s *Service) makeRequest(targetUrl, method string, req, response interface{}) error {
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
