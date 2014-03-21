package api

import (
  "testing"
)

func TestToUrl(t *testing.T) {

  pmt := &Payment{
    Amount: "1.1",
    AccessToken: "zz",
  }

  
  expected := "amount=1.1&access_token=zz" 
  if msg := UrlEncode(pmt); msg != expected {
    t.Errorf("Got %v, expected %v", msg, expected)
  }
}

func TestSpacerize(t *testing.T) {
  tests := []struct {
    input string
    expected string
  }{
    { "Normal", "normal"},
    { "AccessToken", "access_token"},
    { "UserID", "user_id"},
  }

  for _, test := range tests {
    if got := Spacerize(test.input); got != test.expected {
      t.Errorf("Expected %v, got: %v", test.expected, got)
    }
  }

}
