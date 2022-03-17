package service

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/portfolio-report/pr-api/db"
)

type configService struct {
	Db                    db.Config
	MailerTransport       string
	ContactRecipientEmail string
	SessionTimeout        time.Duration
	Ip2locToken           string
	SearchMaxResults      int
}

func DefaultAtoi(s string, def int) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}

// ReadConfig reads and parses configuration from environment variables.
// It also sets default values where applicable.
func ReadConfig() *configService {
	c := configService{}
	c.Db = db.Config{}

	c.MailerTransport = os.Getenv("MAILER_TRANSPORT")
	c.ContactRecipientEmail = os.Getenv("CONTACT_RECIPIENT_EMAIL")

	c.Db.Url = os.Getenv("DATABASE_URL")
	if c.Db.Url == "" {
		log.Fatalln("DATABASE_URL not set")
	}
	c.Db.MaxOpenConn = DefaultAtoi(os.Getenv("DATABASE_MAX_OPEN_CONN"), 25)
	c.Db.MaxIdleConn = DefaultAtoi(os.Getenv("DATABASE_MAX_IDLE_CONN"), 25)
	c.Db.ConnMaxLife = time.Duration(
		DefaultAtoi(os.Getenv("DATABASE_CONN_MAX_LIFE"), 5*60), // 5mins
	) * time.Second

	c.SessionTimeout = time.Duration(
		DefaultAtoi(os.Getenv("SESSION_TIMEOUT"), 15*60), // 15min
	) * time.Second

	c.Ip2locToken = os.Getenv("IP2LOCATION_TOKEN")
	c.SearchMaxResults = DefaultAtoi(os.Getenv("SECURITIES_SEARCH_MAX_RESULTS"), 10)

	return &c
}
