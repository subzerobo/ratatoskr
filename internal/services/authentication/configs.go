package authentication

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

// Config Holds required configuration for Authentication
type Config struct {
	OAuthProviders map[string]OAUTH `yaml:"OAUTH_PROVIDERS"`
	JWT            JWTConfig        `yaml:"JWT"`
}

func (c Config) GetOAuthProviders() map[string]oauth2.Config {
	// TODO: Add other providers
	out := make(map[string]oauth2.Config)
	for k, v := range c.OAuthProviders {
		switch k {
		case "google":
			out["google"] = oauth2.Config{
				ClientID:     v.ClientID,
				ClientSecret: v.ClientSecret,
				Endpoint:     google.Endpoint,
				RedirectURL:  v.RedirectURL,
				Scopes:       v.Scopes,
			}
		case "facebook":
			out["facebook"] = oauth2.Config{
				ClientID:     v.ClientID,
				ClientSecret: v.ClientSecret,
				Endpoint:     facebook.Endpoint,
				RedirectURL:  v.RedirectURL,
				Scopes:       v.Scopes,
			}
		}
	}
	return out
}

type OAUTH struct {
	ClientID     string   `yaml:"CLIENT_ID" envconfig:"CLIENT_ID"`
	ClientSecret string   `yaml:"CLIENT_SECRET" envconfig:"CLIENT_SECRET"`
	RedirectURL  string   `yaml:"REDIRECT_URL" envconfig:"REDIRECT_URL"`
	Scopes       []string `yaml:"SCOPES" envconfig:"SCOPES"`
}

// JWTConfig Stores JWT configurations
type JWTConfig struct {
	Audience   string `yaml:"AUDIENCE" envconfig:"JWT_AUDIENCE"`
	Issuer     string `yaml:"ISSUER" envconfig:"JWT_ISSUER"`
	Secret     string `yaml:"SECRET" envconfig:"JWT_SECRET"`
	Version    string `yaml:"VERSION" envconfig:"JWT_VERSION"`
	ExpireDays int    `yaml:"EXPIRE" envconfig:"JWT_EXPIRE_DAYS"`
}
