package notifier

import (
	"cakebot/internal/data"
	"net/http"
	"reflect"
	"testing"
)

func TestNotifier_ByCalendarInvite(t *testing.T) {
	type fields struct {
		client *http.Client
		cake   data.CakeModel
	}
	type args struct {
		in0 *data.Cake
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notify := Notifier{
				client: tt.fields.client,
			}
			notify.SendCalendarInvite(tt.args.in0)
		})
	}
}

func TestNotifier_ByEmail(t *testing.T) {
	type fields struct {
		client *http.Client
		cake   data.CakeModel
	}
	type args struct {
		cake *data.Cake
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notify := Notifier{
				client: tt.fields.client,
			}
			if err := notify.SendEmail(tt.args.cake); (err != nil) != tt.wantErr {
				t.Errorf("SendEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNotifier_BySms(t *testing.T) {
	type fields struct {
		client *http.Client
		cake   data.CakeModel
	}
	type args struct {
		cake *data.Cake
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notify := Notifier{
				client: tt.fields.client,
			}
			if err := notify.SendSms(tt.args.cake); (err != nil) != tt.wantErr {
				t.Errorf("SendSms() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNotifier_ExtractLdapData(t *testing.T) {
	type fields struct {
		client *http.Client
		cake   data.CakeModel
	}
	type args struct {
		ldap []uint8
	}

	TestFields := fields{
		client: nil,
		cake:   data.CakeModel{},
	}
	newStuff := []uint8{123, 34, 99, 110, 34, 58, 91, 34, 71, 97, 108, 105, 108, 101, 111, 32, 71, 97, 108, 105, 108, 101, 105, 34, 93, 44, 34, 109, 97, 105, 108, 34, 58, 91, 34, 103, 97, 108, 105, 101, 108, 101, 111, 64, 108, 100, 97, 112, 46, 102, 111, 114, 117, 109, 115, 121, 115, 46, 99, 111, 109, 34, 93, 44, 34, 111, 98, 106, 101, 99, 116, 67, 108, 97, 115, 115, 34, 58, 91, 34, 105, 110, 101, 116, 79, 114, 103, 80, 101, 114, 115, 111, 110, 34, 44, 34, 111, 114, 103, 97, 110, 105, 122, 97, 116, 105, 111, 110, 97, 108, 80, 101, 114, 115, 111, 110, 34, 44, 34, 112, 101, 114, 115, 111, 110, 34, 44, 34, 116, 111, 112, 34, 93, 44, 34, 115, 110, 34, 58, 91, 34, 71, 97, 108, 105, 108, 101, 105, 34, 93, 44, 34, 117, 105, 100, 34, 58, 91, 34, 103, 97, 108, 105, 101, 108, 101, 111, 34, 93, 44, 34, 109, 111, 98, 105, 108, 101, 34, 58, 91, 34, 43, 52, 55, 57, 49, 55, 54, 54, 55, 53, 48, 34, 93, 125}
	ExpectedObject := LdapData{
		Cn:          []string{"Galileo Galilei"},
		Mail:        []string{"galieleo@ldap.forumsys.com"},
		ObjectClass: []string{"inetOrgPerson", "organizationalPerson", "person", "top"},
		Sn:          []string{"Galilei"},
		Mobile:      []string{"+4791766750"},
		UID:         []string{"galieleo"},
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *LdapData
		wantErr bool
	}{
		{"NoDataToExtract", TestFields, args{ldap: nil}, nil, true},
		{"ExtractedGalileo", TestFields, args{ldap: newStuff}, &ExpectedObject, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notify := Notifier{
				client: tt.fields.client,
			}
			got, err := notify.ExtractLdapData(tt.args.ldap)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractLdapData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractLdapData() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotifier_NewEmailTemplate(t *testing.T) {
	type fields struct {
		client *http.Client
		cake   data.CakeModel
	}
	type args struct {
		bodyContent   string
		theRecipients ToRecipients
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *EmailMessage
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notify := Notifier{
				client: tt.fields.client,
			}
			if got := notify.NewEmailTemplate(tt.args.bodyContent, tt.args.theRecipients); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEmailTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotifier_SendMail(t *testing.T) {
	type fields struct {
		client *http.Client
		cake   data.CakeModel
	}
	type args struct {
		theEmailMessage *EmailMessage
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notify := Notifier{
				client: tt.fields.client,
			}
			notify.SendMail(tt.args.theEmailMessage)
		})
	}
}

func TestSetEmailRecipient(t *testing.T) {
	type args struct {
		Recipient string
	}

	ValidEmailShouldReturnRecipientStruct := make(ToRecipients, 1)
	ValidEmailShouldReturnRecipientStruct[0].EmailAddress.Address = "valid@email.com"
	ValidPhoneNumberEmail := make(ToRecipients, 1)
	ValidPhoneNumberEmail[0].EmailAddress.Address = "+4791766750@sms.gw"

	tests := []struct {
		name    string
		args    args
		want    ToRecipients
		wantErr bool
	}{
		{"InvalidEmailCheck", args{Recipient: "invalidemail"}, nil, true},
		{"ValidEmailShouldReturnRecipientStruct", args{Recipient: "valid@email.com"}, ValidEmailShouldReturnRecipientStruct, false},
		{"ValidPhoneNumberEmail", args{Recipient: "+4791766750@sms.gw"}, ValidPhoneNumberEmail, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SetEmailRecipient(tt.args.Recipient)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetEmailRecipient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetEmailRecipient() got = %v, want %v", got, tt.want)
			}
		})
	}
}
