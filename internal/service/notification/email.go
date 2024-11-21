package notification

import (
	"log"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func SendReservationConfirmationEmail(to string, reservationDetails string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "reservations@app.com")
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", "Reservation Confirmation")
	mailer.SetBody("text/html", reservationDetails)

	/*emailPassword := os.Getenv("EMAIL_PASSWORD") //TODO: agregar contraseñas al .env y luego tiene que ir en el lugar de la password en dialer
	if emailPassword == "" {
		log.Printf("Email password not set in environment variables")
		return fmt.Errorf("email password not set")
	}*/

	dialer := gomail.NewDialer("smtp.gmail.com", 587, "agmgerminario@gmail.com", "iwos arnt ngur yfte")

	if err := dialer.DialAndSend(mailer); err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Printf("Email sent successfully to %s", to)
	return nil
}
