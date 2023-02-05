package notifier

import (
	"bytes"
	"cakebot/internal/data"
	"cakebot/internal/data/validator"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
	"os"
	"strings"
	"time"
)

type EmailMessage struct {
	Message struct {
		Subject string `json:"subject,omitempty"`
		Body    struct {
			ContentType string `json:"contentType,omitempty"`
			Content     string `json:"content,omitempty"`
		} `json:"body,omitempty"`
		ToRecipients []struct {
			EmailAddress struct {
				Address string `json:"address,omitempty"`
			} `json:"emailAddress,omitempty"`
		} `json:"toRecipients,omitempty"`
	} `json:"message,omitempty"`
}

type ToRecipients []struct {
	EmailAddress struct {
		Address string `json:"address,omitempty"`
	} `json:"emailAddress,omitempty"`
}

type Attendees []struct {
	EmailAddress struct {
		Address string `json:"address"`
		Name    string `json:"name"`
	} `json:"emailAddress"`
	Type string `json:"type"`
}

type LdapData struct {
	Cn          []string `json:"cn"`
	Mail        []string `json:"mail"`
	ObjectClass []string `json:"objectClass"`
	Sn          []string `json:"sn"`
	Mobile      []string `json:"mobile"`
	UID         []string `json:"uid"`
}

type calendarEvent struct {
	Subject string `json:"subject"`
	Body    struct {
		ContentType string `json:"contentType"`
		Content     string `json:"content"`
	} `json:"body"`
	Start struct {
		DateTime string `json:"dateTime"`
		TimeZone string `json:"timeZone"`
	} `json:"start"`
	End struct {
		DateTime string `json:"dateTime"`
		TimeZone string `json:"timeZone"`
	} `json:"end"`
	Location struct {
		DisplayName string `json:"displayName"`
	} `json:"location"`
	Attendees []struct {
		EmailAddress struct {
			Address string `json:"address"`
			Name    string `json:"name"`
		} `json:"emailAddress"`
		Type string `json:"type"`
	} `json:"attendees"`
}

type Notifier struct {
	client *http.Client
}

func New(cfg clientcredentials.Config) Notifier {
	ctx := context.Background()
	oauthClient := &clientcredentials.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		TokenURL:     cfg.TokenURL,
		Scopes:       []string{".default"},
	}
	authorizedClient := oauthClient.Client(ctx)
	newNotifier := Notifier{client: authorizedClient}
	return newNotifier
}

func (notify Notifier) EmailNotifications() bool {
	emailNotifications := os.Getenv("EMAIL_NOTIFICATIONS")

	if emailNotifications == strings.ToUpper("TRUE") {
		return true
	}

	return false

}

func (notify Notifier) SMSNotifications() bool {
	emailNotifications := os.Getenv("SMS_NOTIFICATIONS")

	if emailNotifications == strings.ToUpper("TRUE") {
		return true
	}
	return false

}

func (notify Notifier) TeamsNotifications() bool {
	emailNotifications := os.Getenv("TEAMS_NOTIFICATIONS")

	if emailNotifications == strings.ToUpper("TRUE") {
		return true
	}
	return false

}

func (notify Notifier) CalendarNotificationsEnabled() bool {
	emailNotifications := os.Getenv("CALENDAR_NOTIFICATION")

	if emailNotifications == strings.ToUpper("TRUE") {
		return true
	}
	return false

}

func (notify Notifier) enabled() bool {
	if notify.EmailNotifications() || notify.SMSNotifications() || notify.TeamsNotifications() || notify.CalendarNotificationsEnabled() {
		return true
	}
	return false
}

func (notify Notifier) Notify(cake data.CakeModel, model data.CakeModel) error {

	if notify.enabled() {

		for {
			for {
				currentTime := time.Now()
				cake, err := cake.GetCurrent()
				fmt.Println(cake.When, cake.User_ID)
				if err != nil {
					fmt.Println(err)
				}

				if cake.Notified == 0 && data.CorrectNotitifactionTime(currentTime, cake) {
					err = notify.SendEmail(cake)
					err = notify.SendSms(cake)
					err = notify.SendTeamsMessage(cake)
					err = notify.SendCalendarInvite(cake)
					if err != nil {
						// we have notified now, if all is well we set notification flag.
						fmt.Println("notifications failed", err)
					}
					model.SetNotified(cake)
					break

				}
				time.Sleep(time.Minute * 5)
			}
		}
	} else {
		fmt.Println("notifications are not enabled")
	}
	return nil
}

func (notify Notifier) SendSms(cake *data.Cake) error {
	SmsEnabled := os.Getenv("SMS_NOTIFICATIONS")
	if strings.ToUpper(SmsEnabled) == "TRUE" {

		ldapData, err := notify.ExtractLdapData(cake.Ldapdata)
		if err != nil {
			fmt.Println("could not extract ldap data", err)
			return data.LdapExtractionError
		}

		if ldapData.Mobile[0] != "" {
			PhoneNr := fmt.Sprintf("%s@smsgw.sms", ldapData.Mobile[0])
			Recipient, RecipientErr := SetEmailRecipient(PhoneNr)
			if RecipientErr != nil {
				return err
			}
			bodyContent := fmt.Sprintf("Cakebot har bestemt at førstkommende fredag (%s) er din tur til å lage kake", cake.When)
			template := notify.NewEmailTemplate(bodyContent, Recipient)
			notify.SendMail(template)
		} else {
			fmt.Println(data.MailNotFound)
		}
		return nil
	}
	return data.SMSNotificationsNotEnabled
}

func (notify Notifier) ExtractLdapData(ldap []uint8) (*LdapData, error) {
	userLdapData := LdapData{}
	err := json.Unmarshal(ldap, &userLdapData)
	if err != nil {
		return nil, err
	}
	return &userLdapData, nil

}

func (notify Notifier) SendEmail(cake *data.Cake) error {
	emailEnabled := os.Getenv("EMAIL_NOTIFICATIONS")
	if strings.ToUpper(emailEnabled) == "TRUE" {

		ldapData, err := notify.ExtractLdapData(cake.Ldapdata)
		if err != nil {
			fmt.Println("could not extract ldap data", err)
			return data.LdapExtractionError
		}

		if ldapData.Mail[0] != "" {
			Recipient, RecipientErr := SetEmailRecipient(ldapData.Mail[0])
			if RecipientErr != nil {
				return err
			}
			bodyContent := fmt.Sprintf("Cakebot har bestemt at førstkommende fredag (%s) er din tur til å lage kake", cake.When)
			template := notify.NewEmailTemplate(bodyContent, Recipient)
			notify.SendMail(template)
		} else {
			fmt.Println(data.MailNotFound)
		}
		return nil
	}
	return data.MailNotificationsNotEnabled
}
func (notify Notifier) SendTeamsMessage(*data.Cake) error {

	return nil
}

func (notify Notifier) EventResponses() {

}

func (notify Notifier) Events() {

}

func (notify Notifier) SendCalendarInvite(*data.Cake) error {

	return nil
}

func (notify Notifier) NewEmailTemplate(bodyContent string, theRecipients ToRecipients) *EmailMessage {
	aNewMessage := EmailMessage{}
	aNewMessage.Message.Body.ContentType = "Text"
	aNewMessage.Message.Body.Content = bodyContent
	aNewMessage.Message.ToRecipients = theRecipients
	return &aNewMessage
}

func (notify Notifier) SendMail(theEmailMessage *EmailMessage) {
	mailUserID := os.Getenv("AUTH_USERCALENDERID")

	if mailUserID != "" {
		mailEndPoint := fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%s/sendMail", mailUserID)
		contentType := "application/json"

		messageJson, err := json.Marshal(theEmailMessage)

		response, err := notify.client.Post(mailEndPoint, contentType, bytes.NewBuffer(messageJson))
		if err != nil {
			fmt.Println("http POST", err)
		}

		if response.StatusCode == 202 {
			fmt.Printf("Cakebot successfully sent email to %s\n", theEmailMessage.Message.ToRecipients[0].EmailAddress)

		} else {
			fmt.Printf("Cakebot did not succeed in sending email.. Full response below:\n %v \n", response)
		}
	} else {
		fmt.Println("AUTH_USERCALENDERID is not set. Cannot send SMS.")
		os.Exit(1)
	}

}

func SetEmailRecipient(Recipient string) (ToRecipients, error) {

	TheRequiredAttendees := make(ToRecipients, 1)
	v := validator.New()
	v.Check(Recipient != "", "email", "must be provided")
	v.Check(validator.Matches(Recipient, validator.EmailRX), "email", "must be a valid email address")
	if len(v.Errors) == 0 {
		TheRequiredAttendees[0].EmailAddress.Address = Recipient
		return TheRequiredAttendees, nil
	}
	return nil, data.InvalidEmail

}

func (notify Notifier) WebHook() {

}
