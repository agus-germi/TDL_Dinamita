package notification

import (
	"strconv"

	"github.com/agus-germi/TDL_Dinamita/logger"
	"github.com/agus-germi/TDL_Dinamita/utils"
	"gopkg.in/gomail.v2"
)

var (
	smtpHost string
	smtpPort int
)

func init() {
	logger.Log.Debug("Executing init() function of 'notification' package: Loading SMTP_HOST and SMTP_PORT from '.env' file")

	smtpHost, err := utils.GetEnv("SMTP_HOST")
	if err != nil || smtpHost == "" {
		logger.Log.Fatalf("SMTP_HOST is not set or invalid: %v", err)
	}
	logger.Log.Debugf("Value read from SMTP_HOST: %s", smtpHost)

	smtpPortStr, err := utils.GetEnv("SMTP_PORT")
	if err != nil || smtpPortStr == "" {
		logger.Log.Fatalf("SMTP_PORT is not set or invalid: %v", err)
	}
	logger.Log.Debugf("Value read from SMTP_PORT: %s", smtpPortStr)

	smtpPort, err = strconv.Atoi(smtpPortStr)
	if err != nil {
		logger.Log.Fatalf("Error trying to convert SMTP_PORT environment variable to int: %v", err)
	}

	logger.Log.Infof("SMTP Host and Port loaded successfully from '.env' file.")
}

func SendReservationConfirmationEmail(to string, reservationDetails string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "reservation@app.com")
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", "Reservation Confirmation")
	mailer.SetBody("text/html", reservationDetails)

	dialer := gomail.NewDialer(smtpHost, smtpPort, "agmgerminario@gmail.com", "iwos arnt ngur yfte")

	if err := dialer.DialAndSend(mailer); err != nil {
		logger.Log.Errorf("Failed to send email to %s: %v", to, err)
		return err
	}

	logger.Log.Infof("Email sent successfully to %s", to)
	return nil
}
