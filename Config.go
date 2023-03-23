package frame

import (
	"github.com/app-nerds/configinator"
	"github.com/sirupsen/logrus"
)

type Config struct {
	AppName string
	Version string

	AdminSessionKey    string `flag:"adminsessionkey" env:"ADMIN_SESSION_KEY" default:"my-secret-key" description:"Key used to encrypt admin sessions"`
	AdminSessionMaxAge int    `flag:"adminsessionmaxage" env:"ADMIN_SESSION_MAX_AGE" default:"86400" description:"Number of seconds a session is valid for"`
	AdminSessionName   string `flag:"adminsessionname" env:"ADMIN_SESSION_NAME" default:"" description:"Name of cookie sessions"`
	AutoSSLEmail       string `flag:"autosslemail" env:"AUTO_SSL_EMAIL" default:"" description:"Email address to use for Lets Encrypt"`
	AutoSSLWhitelist   string `flag:"autosslwhitelist" env:"AUTO_SSL_WHITELIST" default:"" description:"Comma-seperated list of domains for SSL"`
	Debug              bool   `flag:"debug" evn:"DEBUG" default:"true" description:"True to turn on debug mode."`
	DSN                string `flag:"dsn" env:"DSN" default:"host=localhost user=postgres password=password dbname=frame port=5432" description:"DSN string to connect to a database"`
	FireplaceURL       string `flag:"fireplaceurl" env:"FIREPLACE_URL" default:"" description:"URL to a Fireplace logging server"`
	FireplacePassword  string `flag:"fireplacepassword" env:"FIREPLACE_PASSWORD" default:"" description:"Password to the Fireplace logging server"`
	GobucketURL        string `flag:"gobucketurl" env:"GOBUCKET_URL" default:"" description:"URL to a Gobucket Server"`
	GobucketClientCode string `flag:"gobucketclientcode" env:"GOBUCKET_CLIENT_CODE" default:"" description:"Client Code to a Gobucket Server"`
	GobucketAppKey     string `flag:"gobucketappkey" env:"GOBUCKET_APP_KEY" default:"" description:"App Key token to connect to a Gobucket Server"`
	GoogleClientID     string `flag:"googleclientid" env:"GOOGLE_CLIENT_ID" default:"" description:"Google OAuth2 client ID"`
	GoogleClientSecret string `flag:"googleclientsecret" env:"GOOGLE_CLIENT_SECRET" default:"" description:"Google OAuth2 client secret"`
	GoogleRedirectURI  string `flag:"googleredirecturi" env:"GOOGLE_REDIRECT_URI" default:"http://localhost:8080/auth/google/callback" description:"Google OAuth2 redirect URI"`
	LogLevel           string `flag:"loglevel" env:"LOG_LEVEL" default:"debug" description:"Minimum log level to report"`
	PageSize           int    `flag:"pagesize" env:"PAGE_SIZE" default:"25" description:"Size of pages for results"`
	RootUserName       string `flag:"rootusername" env:"ROOT_USER_NAME" default:"root" description:"root user name for admin"`
	RootUserPassword   string `flag:"rootUserPassword" env:"ROOT_USER_PASSWORD" default:"password" description:"Password to the root admin user"`
	ServerHost         string `flag:"serverhost" env:"SERVER_HOST" default:"localhost:8080" description:"Host and port to bind to"`
	SessionKey         string `flag:"sessionkey" env:"SESSION_KEY" default:"my-secret-key" description:"Key used to encrypt sessions"`
	SessionMaxAge      int    `flag:"sessionmaxage" env:"SESSION_MAX_AGE" default:"86400" description:"Number of seconds a session is valid for"`
	SessionName        string `flag:"sessionname" env:"SESSION_NAME" default:"" description:"Name of cookie sessions"`
	ServerIdleTimeout  int    `flag:"serveridletimeout" env:"SERVER_IDLE_TIMEOUT" default:"30" description:"Timeout for HTTP idle"`
	ServerReadTimeout  int    `flag:"serverreadtimeout" env:"SERVER_READ_TIMEOUT" default:"60" description:"Timeout for HTTP reads"`
	ServerWriteTimeout int    `flag:"serverwritetimeout" env:"SERVER_WRITE_TIMEOUT" default:"30" description:"Timeout for HTTP writes"`
}

func NewConfig(appName, version string) *Config {
	result := Config{}
	configinator.Behold(&result)

	result.AppName = appName
	result.Version = version

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
