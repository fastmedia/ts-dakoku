package app

import (
	"net/http"
	"testing"
	"time"

	"golang.org/x/oauth2"
	gock "gopkg.in/h2non/gock.v1"
)

func TestGetSalesforceOAuthCallbackURL(t *testing.T) {
	app := createMockApp()
	req, _ := http.NewRequest(http.MethodGet, "https://example.com/test", nil)
	ctx := app.createContext(req)
	Test{"https://example.com/oauth/salesforce/callback", ctx.getSalesforceOAuthCallbackURL()}.Compare(t)
}

func TestGetSalesforceAuthenticateURL(t *testing.T) {
	app := createMockApp()
	req, _ := http.NewRequest(http.MethodGet, "https://example.com/test", nil)
	ctx := app.createContext(req)
	Test{"https://example.com/oauth/salesforce/authenticate/foo", ctx.getSalesforceAuthenticateURL("foo")}.Compare(t)
}

func TestGetSlackOAuthCallbackURL(t *testing.T) {
	app := createMockApp()
	req, _ := http.NewRequest(http.MethodGet, "https://example.com/test", nil)
	ctx := app.createContext(req)
	Test{"https://example.com/oauth/slack/callback", ctx.getSlackOAuthCallbackURL()}.Compare(t)
}

func TestGetSlackeAuthenticateURL(t *testing.T) {
	app := createMockApp()
	req, _ := http.NewRequest(http.MethodGet, "https://example.com/test", nil)
	ctx := app.createContext(req)
	Test{"https://example.com/oauth/slack/authenticate/foo/bar", ctx.getSlackAuthenticateURL("foo", "bar")}.Compare(t)
}

func TestSetAndGetSalesforceAccessToken(t *testing.T) {
	token := &oauth2.Token{
		AccessToken:  "foo",
		RefreshToken: "bar",
		TokenType:    "Bearer",
		Expiry:       time.Now().Add(-10 * time.Hour),
	}
	app := createMockApp()
	app.CleanRedis()
	req, _ := http.NewRequest(http.MethodGet, "https://example.com/test", nil)
	ctx := app.createContext(req)
	err := ctx.setSalesforceAccessToken(token)
	Test{true, err != nil}.Compare(t)
	ctx.UserID = "FOO"
	err = ctx.setSalesforceAccessToken(token)
	Test{false, err != nil}.Compare(t)
	token = ctx.getSalesforceAccessTokenForUser()
	for _, test := range []Test{
		{"foo", token.AccessToken},
		{"bar", token.RefreshToken},
		{"Bearer", token.TokenType},
	} {
		test.Compare(t)
	}
	ctx = app.createContext(req)
	ctx.UserID = "BAR"
	token = ctx.getSalesforceAccessTokenForUser()
	Test{true, token == nil}.Compare(t)
}

func TestSetAndGetSalesforceOAuthClient(t *testing.T) {
	defer gock.Off()
	newExpiry := time.Now().Add(2 * time.Hour).Truncate(time.Second)
	oldExpiry := time.Now().Add(-10 * time.Hour).Truncate(time.Second)
	resExpiry, _ := time.Parse("2016-01-02T15:04:05Z", "0001-01-01T00:00:00Z")
	gock.New("https://login.salesforce.com").
		Post("/services/oauth2/token").
		Reply(200).
		JSON(oauth2.Token{
			AccessToken:  "foo2",
			RefreshToken: "bar2",
			TokenType:    "Bearer2",
			Expiry:       resExpiry,
		})
	token := &oauth2.Token{
		AccessToken:  "foo",
		RefreshToken: "bar",
		TokenType:    "Bearer",
		Expiry:       oldExpiry,
	}
	app := createMockApp()
	app.CleanRedis()
	req, _ := http.NewRequest(http.MethodGet, "https://example.com/test", nil)
	ctx := app.createContext(req)
	ctx.UserID = "FOO"
	ctx.TimeoutDuration = 2 * time.Hour
	err := ctx.setSalesforceAccessToken(token)
	token = ctx.getSalesforceAccessTokenForUser()
	for _, test := range []Test{
		{false, token == nil},
		{oldExpiry.String(), token.Expiry.String()},
		{"bar", token.RefreshToken},
		{"foo", token.AccessToken},
		{"Bearer", token.TokenType},
		{false, token.Expiry.IsZero()},
		{true, err == nil},
	} {
		test.Compare(t)
	}
	client := ctx.getSalesforceOAuth2Client()
	token = ctx.getSalesforceAccessTokenForUser()
	for _, test := range []Test{
		{false, client == nil},
		{newExpiry.String(), token.Expiry.String()},
		{"bar2", token.RefreshToken},
		{"foo2", token.AccessToken},
		{"Bearer2", token.TokenType},
	} {
		test.Compare(t)
	}
}

func TestSetAndGetSlackAccessToken(t *testing.T) {
	app := createMockApp()
	app.CleanRedis()
	req, _ := http.NewRequest(http.MethodGet, "https://example.com/test", nil)
	ctx := app.createContext(req)
	err := ctx.setSlackAccessToken("foo")
	Test{true, err != nil}.Compare(t)
	ctx.UserID = "FOO"
	err = ctx.setSlackAccessToken("foo")
	Test{false, err != nil}.Compare(t)
	token := ctx.getSlackAccessTokenForUser()
	Test{"foo", token}.Compare(t)
	ctx = app.createContext(req)
	ctx.UserID = "BAR"
	token = ctx.getSlackAccessTokenForUser()
	Test{"", token}.Compare(t)
}
