package api

import (
  "encoding/json"
  "fmt"
  "net/http"
)

type Service struct {
  client *http.Client
  baseUrl string
  tokenString string
}

func New(tokenString string) *Service {
  return &Service{
    client: ClientFromToken(tokenString),
    baseUrl: "https://api.venmo.com/v1/",
    tokenString: tokenString,
  }
}

func NewTest(tokenString string) *Service {
  svc := New(tokenString)
  svc.baseUrl = "https://sandbox-api.venmo.com/v1/"
  return svc
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
  Balance float32 `json:"balance,omitempty"`
  Pmt Payment `json:"payment,omitempty"`
}

type Payment struct {
  Status string `json:"status,omitempty"`
  Target PaymentTarget `json:"target,omitempty"`
  DateCompleted string `json:"date_completed,omitempty"`
  Actor PaymentUser `json:"actor,omitempty"`
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
  FriendsCount int `json:"friends_count,omitempty"`
  IsFriend bool `json:"is_friend,omitempty"`
}

type MeRequest struct {
  AccessToken string `json:"access_token,omitempty"`
}

type MeResponse struct {
  Data MeResponseData `json:"data,omitempty"`
}

type MeResponseData struct {
  Balance string `json:"balance,omitempty"` // This one is returned as a string..
  User PaymentUser `json:"user,omitempty"`
}

type ListFriendsRequest struct {
  AccessToken string `json:"access_token,omitempty"`
  Before string `json:"before,omitempty"`
  After string `json:"after,omitempty"`
  Limit string `json:"limit,omitempty"`
}

type ListFriendsResponse struct {
  // Pagination
  Data []PaymentUser
}

func (s *Service) ListPayments(req *ListPaymentsRequest) (*ListPaymentsResponse, error) {
  ret := new(ListPaymentsResponse)
  req.AccessToken = s.tokenString
  if err := s.makeRequest("payments", "GET", *req, ret); err != nil {
    return nil, err
  }
  return ret, nil
}

func (s *Service) Me() (*MeResponse, error) {
  req := &MeRequest{AccessToken: s.tokenString}
  ret := new(MeResponse)
  if err := s.makeRequest("me", "GET", *req, ret); err != nil {
    return nil, err
  }
  return ret, nil
}

func (s *Service) MakePayment(payment *MakePaymentRequest) (*MakePaymentResponse, error) {
  ret := new(MakePaymentResponse)
  payment.AccessToken = s.tokenString
  if err := s.makeRequest("payments", "POST", *payment, ret); err != nil {
    return nil, err
  }
  return ret, nil
}

func (s *Service) ListFriends(userID string, req *ListFriendsRequest) (*ListFriendsResponse, error) {
  req.AccessToken = s.tokenString
  ret := new(ListFriendsResponse)
  target := fmt.Sprintf("users/%s/friends", userID)
  if err := s.makeRequest(target, "GET", *req, ret); err != nil {
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
