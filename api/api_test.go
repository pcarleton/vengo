package api

import (
  "testing"
)

func TestToUrl(t *testing.T) {

  pmt := &Payment{
    Amount: "1.1",
    AccessToken: "zz",
  }

  
  expected := "access_token=zz&amount=1.1" 
  if msg := StructToUrlValues(pmt).Encode(); msg != expected {
    t.Errorf("Got %v, expected %v", msg, expected)
  }
}

