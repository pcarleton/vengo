package main

import (
  "fmt"
  "flag"
  "github.com/pcarleton/vengo/api"
  "os"
)

func init() {
  registerDemo("list", listDemo)
  registerDemo("pay", makePaymentDemo)
}

func usage() {
  fmt.Fprintf(os.Stderr, "Usage: go run main <access-token> <demo-name> [demo args]\n\nPossible demos:\n\n")
        for n, _ := range demoFunc {
                fmt.Fprintf(os.Stderr, "  * %s\n", n)
        }
        os.Exit(2)
}

func main() {
  flag.Parse()
  if flag.NArg() == 0 {
    fmt.Println("Need access token")
    return
  }

  demo, ok := demoFunc[flag.Arg(1)]
  if !ok {
    usage()
  }
  svc := api.NewTest(flag.Arg(0))
  demo(svc, flag.Args()[1:])
}

var (
  demoFunc = make(map[string]func(*api.Service, []string))
)

func makePaymentDemo(svc *api.Service, argv []string) {
  paymentReq := &api.MakePaymentRequest{
    Phone: "15555555555",
    Note: "testing!",
    Amount: "0.10",
    Audience: "private",
  }

  res, err := svc.MakePayment(paymentReq)
  if err != nil {
    fmt.Printf("Error making payment: %v\n", err)
    return
  }
  fmt.Printf("Sucess! Got response: %+v\n", res)
}

func listDemo(svc *api.Service, argv []string) {
  listReq := &api.ListPaymentsRequest{
    Limit: "1",
  }
  res, err := svc.ListPayments(listReq)
  if err != nil {
    fmt.Printf("Error listing payments: %v\n", err)
    return
  }
  fmt.Printf("Sucess! Got response: %+v\n", res)
}

func registerDemo(name string, main func(svc *api.Service, argv []string)) {
  if demoFunc[name] != nil {
    panic(name + " already exists!")
  }
  demoFunc[name] = main
}
