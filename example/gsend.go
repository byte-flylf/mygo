// Copyright 2013 - by Jim Lawless
// License: MIT / X11
// See: http://www.mailsend-online.com/license2013.php
//
// Bear with me ... I'm a Go noob.
// http://lawlessguy.wordpress.com/2013/07/05/sending-gmail-with-go-golang/
 
package main
 
import (
   "log"
   "flag"
   "net/smtp"
   "fmt"
)
 
func main() {
   to := flag.String("t","","destination Internet mail address")
   from := flag.String("f","","the sender's GMail address")
   pwd := flag.String("p","","the sender's password")
   subject := flag.String("s","","subject line of email")
   msg := flag.String("m","","a one-line email message")
   flag.Usage=func() {
      fmt.Printf("Syntax:\n\tGSend [flags]\nwhere flags are:\n")
      flag.PrintDefaults()
   }
 
   fmt.Printf("GSend v 1.01 by Jim Lawless\n")
 
   flag.Parse()
 
   if flag.NFlag() != 5 {
      flag.Usage()
      return
   }
 
   body := "To: " + *to + "\r\nSubject: " +
      *subject + "\r\n\r\n" + *msg
   auth := smtp.PlainAuth("",*from,*pwd,"smtp.gmail.com")
   err := smtp.SendMail("smtp.gmail.com:587",auth,*from,
      []string{*to},[]byte(body))
   if err != nil {
      log.Fatal(err)
   }
}
