package sendmail

import (
	"fmt"
	"net/smtp"
	"strings"
)

// SendMail is to send email by smtp.
func SendMail(host, mailFrom, mailTo, mailType, subject, body string) error {
	msg := []byte(fmt.Sprintf("To: %s\r\nFrom: %s>\r\nSubject: %s\r\nContent-Type: text/%s; charset=UTF-8\r\n\r\n%s", mailTo, mailFrom, subject, mailType, body))
	sendTo := strings.Split(mailTo, ";")
	err := smtp.SendMail(host, nil, mailFrom, sendTo, msg)
	return err
}
