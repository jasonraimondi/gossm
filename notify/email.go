package notify

import (
	"fmt"
	"net/smtp"
	"strings"
)

type EmailSettings struct {
	SMTP     string
	Port     int
	Username string
	Password string
	From     string
	To       []string
}

type EmailNotifier struct {
	Settings *EmailSettings
	auth     smtp.Auth
}

func (es *EmailSettings) Validate() (bool, error) {
	errEmailProperty := func(property string) error {
		return fmt.Errorf("missing email property %s", property)
	}
	switch {
	case es.Username == "":
		return false, errEmailProperty("username")
	case es.Password == "":
		return false, errEmailProperty("password")
	case es.SMTP == "":
		return false, errEmailProperty("smtp")
	case es.Port == 0:
		return false, errEmailProperty("port")
	case es.From == "":
		return false, errEmailProperty("from")
	case len(es.To) == 0:
		return false, errEmailProperty("to")
	}
	return true, nil
}

func (e *EmailNotifier) Initialize() {
	e.auth = smtp.PlainAuth(
		"",
		e.Settings.Username,
		e.Settings.Password,
		e.Settings.SMTP,
	)
}

func (e *EmailNotifier) String() string {
	return fmt.Sprintf("email %s at %s:%d", e.Settings.Username, e.Settings.SMTP, e.Settings.Port)
}

func (e *EmailNotifier) Notify(text string) (bool, error) {
	formattedReceipets := strings.Join(e.Settings.To, ", ")
	msg := "From: " + e.Settings.From + "\n" +
		"To: " + formattedReceipets + "\n" +
		"Subject: GOSSM Notification\n\n" +
		text + " not reached."

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", e.Settings.SMTP, e.Settings.Port),
		e.auth,
		e.Settings.From,
		e.Settings.To,
		[]byte(msg),
	)
	if err != nil {
		return false, fmt.Errorf("error sending email: %s", err)
	}
	return true, nil
}
