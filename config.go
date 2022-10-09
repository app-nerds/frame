package frame

import (
	"github.com/app-nerds/configinator"
	"github.com/sirupsen/logrus"
)

type Config struct {
	AppName string
	Version string

	AutoSSLEmail       string `flag:"autosslemail" env:"AUTO_SSL_EMAIL" default:"" description:"Email address to use for Lets Encrypt"`
	AutoSSLWhitelist   string `flag:"autosslwhitelist" env:"AUTO_SSL_WHITELIST" default:"" description:"Comma-seperated list of domains for SSL"`
	Debug              bool   `flag:"debug" evn:"DEBUG" default:"true" description:"True to turn on debug mode."`
	DSN                string `flag:"dsn" env:"DSN" default:"host=localhost user=postgres password=password dbname=frame port=5432" description:"DSN string to connect to a database"`
	FireplaceURL       string `flag:"fireplaceurl" env:"FIREPLACE_URL" default:"" description:"URL to a Fireplace logging server"`
	FireplacePassword  string `flag:"fireplacepassword" env:"FIREPLACE_PASSWORD" default:"" description:"Password to the Fireplace logging server"`
	GoogleClientID     string `flag:"googleclientid" env:"GOOGLE_CLIENT_ID" default:"" description:"Google OAuth2 client ID"`
	GoogleClientSecret string `flag:"googleclientsecret" env:"GOOGLE_CLIENT_SECRET" default:"" description:"Google OAuth2 client secret"`
	GoogleRedirectURI  string `flag:"googleredirecturi" env:"GOOGLE_REDIRECT_URI" default:"http://localhost:8080/auth/google/callback" description:"Google OAuth2 redirect URI"`
	LogLevel           string `flag:"loglevel" env:"LOG_LEVEL" default:"info" description:"Minimum log level to report"`
	ServerHost         string `flag:"serverhost" env:"SERVER_HOST" default:"localhost:8080" description:"Host and port to bind to"`
	SessionKey         string `flag:"sessionkey" env:"SESSION_KEY" default:"my-secret-key" description:"Key used to encrypt sessions"`
	ServerIdleTimeout  int    `flag:"serveridletimeout" env:"SERVER_IDLE_TIMEOUT" default:"30" description:"Timeout for HTTP idle"`
	ServerReadTimeout  int    `flag:"serverreadtimeout" env:"SERVER_READ_TIMEOUT" default:"60" description:"Timeout for HTTP reads"`
	ServerWriteTimeout int    `flag:"serverwritetimeout" env:"SERVER_WRITE_TIMEOUT" default:"30" description:"Timeout for HTTP writes"`
}

func (fa *FrameApplication) setupConfig() *Config {
	result := Config{}
	configinator.Behold(&result)

	result.AppName = fa.appName
	result.Version = fa.version

	return &result
}

func (c *Config) GetLogLevel() logrus.Level {
	var (
		err      error
		loglevel logrus.Level
	)

	if loglevel, err = logrus.ParseLevel(c.LogLevel); err != nil {
		panic("invalid log level '" + c.LogLevel + "'")
	}

	return loglevel
}
