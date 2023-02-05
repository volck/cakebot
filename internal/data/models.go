package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound           = errors.New("record not found")
	ErrEditConflict             = errors.New("edit conflict")
	CouldNotCreateTemplate      = errors.New("could not create template")
	LdapExtractionError         = errors.New("could not extract ldap data")
	MailNotFound                = errors.New("mail not found in ldap data")
	InvalidEmail                = errors.New("email is invalid")
	MailNotificationsNotEnabled = errors.New("mail notification not enabled")
	SMSNotificationsNotEnabled  = errors.New("SMS notification not enabled")
	ExpectedSingleRowAffected   = errors.New("expected single row affected")
	CannotNotify                = errors.New("cannot start notification service")
)

type Models struct {
	Cake  CakeModel
	Users UserModel
}

func NewModels(db *sql.DB, client KeycloakClient) Models {
	return Models{
		Cake:  CakeModel{DB: db},
		Users: UserModel{KeycloakClient: client, DB: db},
	}
}
