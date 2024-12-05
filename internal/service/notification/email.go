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
	smtpUser string
	smtpPass string
)

func init() {
	logger.Log.Debug("Executing init() function of 'notification' package: Loading SMTP configurations from '.env' file")

	var err error
	smtpHost, err = utils.GetEnv("SMTP_HOST")
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

	smtpUser, err = utils.GetEnv("SMTP_USER")
	if err != nil || smtpUser == "" {
		logger.Log.Fatalf("SMTP_USER is not set or invalid: %v", err)
	}
	logger.Log.Debug("SMTP_USER loaded successfully.")

	smtpPass, err = utils.GetEnv("SMTP_PASS")
	if err != nil || smtpPass == "" {
		logger.Log.Fatalf("SMTP_PASS is not set or invalid: %v", err)
	}
	logger.Log.Debug("SMTP_PASS loaded successfully.")

	logger.Log.Infof("SMTP configurations loaded successfully from '.env' file.")
}

func SendReservationConfirmationEmail(to string, reservationDetails string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", smtpUser)
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", "Reservation Confirmation")
	mailer.SetBody("text/html", reservationDetails)

	dialer := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	if err := dialer.DialAndSend(mailer); err != nil {
		logger.Log.Errorf("Failed to send email to %s: %v", to, err)
		return err
	}

	logger.Log.Infof("Email sent successfully to %s", to)
	return nil
}
