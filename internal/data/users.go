package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
)

type User struct {
	Id               string `json:"-"`
	CreatedTimestamp int64  `json:"-"`
	Username         string `json:"username"`
	Enabled          bool   `json:"-"-`
	Totp             bool   `json:"-" -`
	EmailVerified    bool   `json:"-" -`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	Email            string `json:"email,omitempty"`
	FederationLink   string `json:"-"-`
	Attributes       struct {
		LDAPENTRYDN     []string `json:"-" - `
		CreateTimestamp []string `json:"-"-`
		ModifyTimestamp []string `json:"-"-`
		LDAPID          []string `json:"-"-`
	} `json:"-"`
	DisableableCredentialTypes []interface{} `json:"-"-`
	RequiredActions            []interface{} `json:"-"-`
	NotBefore                  int           `json:"-"-`
	Access                     struct {
		ManageGroupMembership bool `json:"-"-`
		View                  bool `json:"-"-`
		MapRoles              bool `json:"-"-`
		Impersonate           bool `json:"-"-`
		Manage                bool `json:"-"-`
	} `json:"-"-`
}

type KeycloakClient struct {
	Config oauth2.Config
}

type UserModel struct {
	KeycloakClient KeycloakClient
	DB             *sql.DB
}

func (c KeycloakClient) GetToken() (*oauth2.Token, error) {

	ctx := context.Background()
	token, err := c.Config.PasswordCredentialsToken(ctx, "admin", "admin")
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (c KeycloakClient) NewClient() *http.Client {

	token, err := c.GetToken()
	if err != nil {
		fmt.Println(err)
	}
	ctx := context.Background()

	return c.Config.Client(ctx, token)

}

func (u UserModel) Get() {

}

func (u UserModel) GetAll() ([]User, error) {

	Users := []User{}
	client := u.KeycloakClient.NewClient()

	req, err := client.Get("http://localhost:8080/admin/realms/cakebot/users")
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(req.Body).Decode(&Users)
	if err != nil {
		return nil, err
	}
	return Users, nil

}

func (u UserModel) Delete() {

}
func (u UserModel) Insert() {

}

func (u UserModel) Create() {

}

func (u UserModel) Update() {

}

func (u UserModel) SyncDB() {

}
