package notification

import (
	"log"

	"gopkg.in/gomail.v2"
)

func SendReservationConfirmationEmail(to string, reservationDetails string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "reservations@app.com")
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", "Reservation Confirmation")
	mailer.SetBody("text/html", reservationDetails)

	dialer := gomail.NewDialer("smtp.gmail.com", 587, "agmgerminario@gmail.com", "iwos arnt ngur yfte")

	if err := dialer.DialAndSend(mailer); err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Printf("Email sent successfully to %s", to)
	return nil
}
