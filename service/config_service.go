package service

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/portfolio-report/pr-api/db"
)

// Config holds the parameters read from environment variables
type Config struct {
	Db                    db.Config
	MailerTransport       string
	CacheMaxAge           time.Duration
	ContactRecipientEmail string
	SessionTimeout        time.Duration
	Ip2locToken           string
	SearchMaxResults      int
	AwsAccessKeyID        string
	AwsSecretAccessKey    string
	AwsRegion             string
	AwsLogoBucket         string
	AwsLogoBucketURL      string
}

func defaultAtoi(s string, def int) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}

// ReadConfig reads and parses configuration from environment variables.
// It also sets default values where applicable.
func ReadConfig() *Config {
	c := Config{}
	c.Db = db.Config{}

	c.MailerTransport = os.Getenv("MAILER_TRANSPORT")
	c.ContactRecipientEmail = os.Getenv("CONTACT_RECIPIENT_EMAIL")

	c.Db.Url = os.Getenv("DATABASE_URL")
	if c.Db.Url == "" {
		log.Fatalln("DATABASE_URL not set")
	}
	c.Db.MaxOpenConn = defaultAtoi(os.Getenv("DATABASE_MAX_OPEN_CONN"), 25)
	c.Db.MaxIdleConn = defaultAtoi(os.Getenv("DATABASE_MAX_IDLE_CONN"), 25)
	c.Db.ConnMaxLife = time.Duration(
		defaultAtoi(os.Getenv("DATABASE_CONN_MAX_LIFE"), 5*60), // 5mins
	) * time.Second

	c.CacheMaxAge = time.Duration(
		defaultAtoi(os.Getenv("CACHE_MAX_AGE"), 10*60), // 10min
	) * time.Second
	c.SessionTimeout = time.Duration(
		defaultAtoi(os.Getenv("SESSION_TIMEOUT"), 15*60), // 15min
	) * time.Second

	c.Ip2locToken = os.Getenv("IP2LOCATION_TOKEN")
	c.SearchMaxResults = defaultAtoi(os.Getenv("SECURITIES_SEARCH_MAX_RESULTS"), 10)

	c.AwsAccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	c.AwsSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	c.AwsRegion = os.Getenv("AWS_REGION")
	c.AwsLogoBucket = os.Getenv("AWS_LOGO_BUCKET")
	c.AwsLogoBucketURL = os.Getenv("AWS_LOGO_BUCKET_URL")

	return &c
}
