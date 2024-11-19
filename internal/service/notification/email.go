package notification

import (
	"gopkg.in/gomail.v2"
	"log"
)

func SendReservationConfirmationEmail(to string, reservationDetails string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "mi_email@gmail.com")
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", "Reservation Confirmation")
	mailer.SetBody("text/plain", reservationDetails)

	dialer := gomail.NewDialer("smtp.gmail.com", 587, "mi_email@gmail.com", "mi_email_password")

	if err := dialer.DialAndSend(mailer); err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Printf("Email sent successfully to %s", to)
	return nil
}
