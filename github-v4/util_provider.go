package github

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"github.com/google/go-github/v28/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

const (
	PROVIDER_BASE_URL            = "base_url"
	PROVIDER_ORGANIZATION        = "organization"
	PROVIDER_TOKEN               = "token"
	PROVIDER_APP                 = "app"
	PROVIDER_APP_PEM             = "pem"
	PROVIDER_APP_ID              = "id"
	PROVIDER_APP_INSTALLATION_ID = "installation"
)

type Config struct {
	BaseURL        string
	Organization   string
	Token          string
	Pem            string
	AppID          string
	InstallationID string
}

type Organization struct {
	Name          string
	GraphQLClient *githubv4.Client
	RESTClient    *github.Client
	StopContext   context.Context
}

type TokenResponse struct {
	Token string `json:"token"`
}

func (c *Config) Clients() (interface{}, error) {
	var org Organization

	token := c.Token
	if token == "" && c.InstallationID != "" {
		t, err := newAppToken(c)
		if err != nil {
			return nil, fmt.Errorf("error returning GitHub App installation token: %w", err)
		}
		token = t
	}

	httpClient := oauth2.NewClient(
		context.Background(),
		oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		),
	)

	// Create GraphQL Client
	uGQL, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, err
	}
	uGQL.Path = path.Join(uGQL.Path, "graphql")
	graphQLClient := githubv4.NewEnterpriseClient(uGQL.String(), httpClient)

	// Create Rest Client
	uREST, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, err
	}
	if uREST.String() != "https://api.github.com/" {
		uREST.Path = path.Join(uREST.Path, "v3")
	}
	restClient, err := github.NewEnterpriseClient(uREST.String(), "", httpClient)
	if err != nil {
		return nil, err
	}

	org.GraphQLClient = graphQLClient
	org.RESTClient = restClient
	org.Name = c.Organization
	return &org, nil
}

func newAppToken(c *Config) (string, error) {
	c.Pem = strings.ReplaceAll(c.Pem, "\\n", "\n")
	rsaPrivate, err := crypto.ParseRSAPrivateKeyFromPEM([]byte(c.Pem))
	if err != nil {
		return "", err
	}

	claims := jws.Claims{}
	claims.SetIssuedAt(time.Now())
	claims.SetExpiration(time.Now().Add(time.Duration(10) * time.Second))
	claims.SetIssuer(c.AppID)
	jwt := jws.NewJWT(claims, crypto.SigningMethodRS256)

	bearer, err := jwt.Serialize(rsaPrivate)
	if err != nil {
		return "", err
	}

	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return "", err
	}
	tokenURL := fmt.Sprintf("%s://%s/app/installations/%s/access_tokens", u.Scheme, u.Host, c.InstallationID)
	req, err := http.NewRequest("POST", tokenURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearer))
	req.Header.Set("Accept", "application/vnd.github.machine-man-preview+json")

	res, err := http.DefaultClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return "", err
	}
	if res.StatusCode != 201 {
		return "", fmt.Errorf("status code returned (%d) is not 201", res.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	tokenRes := TokenResponse{}
	err = json.Unmarshal([]byte(string(bodyBytes)), &tokenRes)
	if err != nil {
		return "", err
	}

	return tokenRes.Token, nil
}
