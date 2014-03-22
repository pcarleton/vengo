package main

import (
  "fmt"
  "flag"
  "net/http"
  "github.com/pcarleton/vengo/api"
  "code.google.com/p/goauth2/oauth"
)

func main() {
  flag.Parse()
  if flag.NArg() == 0 {
    fmt.Println("Need access token")
    return
  }

  token := &oauth.Token{
    AccessToken: flag.Arg(0),
  }
  t := &oauth.Transport{
                Token:     token,
                Transport: http.DefaultTransport,
  }
  client := t.Client()
  svc := api.New(client)

  //payment := &api.Payment{
  //  AccessToken: flag.Arg(0),
  //  Phone: flag.Arg(1),
  //  Note: "testing!",
  //  Amount: "-1.0",
  //  Audience: "private",
  //}

  //res, err := svc.MakePaymento(payment)

  listReq := &api.ListPaymentsRequest{
    AccessToken: flag.Arg(0),
    Limit: "1",
  }
  res, err := svc.ListPayments(listReq)
  if err != nil {
    fmt.Printf("Error making payment: %v\n", err)
    return
  }
  fmt.Printf("Sucess! Got response: %+v\n", res)
}
