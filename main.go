package main

import (
	"flag"
	"fmt"
	"net"
	"net/smtp"
	"os"
	"time"
)

func main() {
	var name, domain string

	flag.StringVar(&name, "name", "", "The name of the person we'd like to find the email of")
	flag.StringVar(&domain, "domain", "", "The domain we're this person's email is registered on")
	flag.Parse()

	if name == "" {
		fmt.Println("Name must be defined!")
		os.Exit(1)
	}
	if domain == "" {
		fmt.Println("Domain must be defined!")
		os.Exit(1)
	}

	fmt.Printf("name='%s' domain='%s'\n", name, domain)

	fmt.Println("Testing ISP...")

	if !testIsp() {
		fmt.Println()
		fmt.Println("It seems that we couldn't connect to the gmail mail server\n" +
			"This is likely due to either due to not being connected\n" +
			"to the internet or your ISP blocking outgoing requests on\n" +
			"port 25. You may want to try on another network or hotspot.")
		return
	}

	fmt.Println("ISP test succeeded!")

	fmt.Println("Getting mail server...")
	mxs, err := net.LookupMX(domain)
	if err != nil {
		fmt.Println("Error finding mail server for", domain, ":", err)
		return
	}
	if len(mxs) == 0 {
		fmt.Println("No mail server records found for", domain)
		fmt.Println("Perhaps there isn't a mail server configured?")
	}

	host := mxs[0].Host
	fmt.Println("Using host", host)

	emails := generateEmails(name, domain)
	testEmails(emails, host)

}

func testIsp() bool {
	// Any error occuring in this function means there is probably
	// a connection issue (e.g. ISP blocking outgoing requests)
	timeout, _ := time.ParseDuration("5s")
	dialer := net.Dialer{Timeout: timeout}
	conn, err := dialer.Dial("tcp", "gmail-smtp-in.l.google.com:25")
	if err != nil {
		return false
	}
	client, err := smtp.NewClient(conn, "gmail-smtp-in.l.google.com:25")
	if err != nil {
		return false
	}

	client.Hello("client.example.com")
	client.Mail("")
	err = client.Rcpt("alexandersonone@gmail.com")
	client.Quit()

	if err != nil {
		return false
	}
	return true
}
