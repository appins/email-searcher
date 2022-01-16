package main

import (
	"fmt"
	"net/smtp"
	"time"
)

func testEmail(client *smtp.Client, email string) (bool, time.Duration, error) {
	errDuration := time.Duration(0)

	err := client.Noop()
	if err != nil {
		return false, errDuration, err
	}

	err = client.Mail("")
	if err != nil {
		return false, errDuration, err
	}

	startTime := time.Now()
	err = client.Rcpt(email)
	client.Reset()
	if err != nil {
		return false, errDuration, nil
	}

	return true, time.Since(startTime), nil
}

func testEmails(emails []string, host string) error {
	client, err := smtp.Dial(host + ":25")
	if err != nil {
		return err
	}

	client.Hello("client.example.com")

	var validEmails []string
	var longSuspects []string
	for _, email := range emails {
		fmt.Printf("Testing email %s... ", email)
		valid, elapsed, err := testEmail(client, email)
		if err != nil {
			return err
		}

		if valid {
			validEmails = append(validEmails, email)
			fmt.Println("seemed to work, ( took", elapsed, ")")
			if elapsed.Seconds() >= 1 {
				longSuspects = append(longSuspects, email)
			}
		} else {
			fmt.Println("invalid email.")
		}
	}

	fmt.Println()
	fmt.Println("The following emails we found to be valid:")
	for _, email := range validEmails {
		fmt.Println(email)
	}

	if len(validEmails) == len(emails) {
		fmt.Println()
		fmt.Println("HOWEVER!!!")
		fmt.Println("It seems that no emails failed to verify. It's very")
		fmt.Println("likely the email server is configured to not allow")
		fmt.Println("probing for email addresses like we have here...")
		fmt.Println("Sorry :(")
		fmt.Println()
		fmt.Println("However, I've found that sometimes the one real email can")
		fmt.Println("have a longer delay compared to all of the others.")
		fmt.Println("These emails took at least 1 second to process:")
		for _, email := range longSuspects {
			fmt.Println(email)
		}
	}

	return nil
}
