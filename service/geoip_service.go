package service

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	ip2loc "github.com/ip2location/ip2location-go"
	"github.com/portfolio-report/pr-api/models"
)

type geoipService struct {
	db    *ip2loc.DB
	Token string
}

const ip2locationDirectory = "ip2location"
const ip2locationFilename = "IP2LOCATION-LITE-DB1.IPV6.BIN"
const ip2locationApiFile = "DB1LITEBINIPV6"

func NewGeoipService(ip2locToken string) models.GeoipService {
	s := &geoipService{
		Token: ip2locToken,
	}

	go s.loadIp2LocationDb()

	return s
}

// Looks up country code for ip
func (s *geoipService) GetCountryFromIp(ip string) string {
	if s.db == nil {
		return ""
	}

	results, err := s.db.Get_country_short(ip)

	if err != nil {
		log.Printf("GeoIP: Failed to get country (ip: %s): %s\n", ip, err.Error())
		return ""
	}

	return results.Country_short
}

// Opens local database file, retrieves database file if necessary
func (s *geoipService) loadIp2LocationDb() {
	dbPath := path.Join(ip2locationDirectory, ip2locationFilename)

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Println("GeoIP: Database file not found, trying to retrieve...")
		err := s.retrieveIp2LocationDb()
		if err != nil {
			log.Println("GeoIP:", err.Error())
			log.Println("GeoIP: Will continue to work without GeoIP data.")
			return
		}
		log.Println("Geoip: Database retrieved.")
	}

	log.Println("GeoIP: Loading database...")
	db, err := ip2loc.OpenDB(dbPath)

	if err != nil {
		log.Printf("GeoIP: Could not load database: %s", err)
		return
	}

	s.db = db
	log.Println("GeoIP: Database loaded.")
}

// Downloads, unzips and saves database file
func (s *geoipService) retrieveIp2LocationDb() error {
	zipped, err := s.downloadIp2LocationDb()
	if err != nil {
		return fmt.Errorf("could not download: %w", err)
	}

	if err := os.MkdirAll(ip2locationDirectory, os.ModePerm); err != nil {
		return fmt.Errorf("could not create directory: %w", err)

	}

	if err := s.unzipIp2LocationDb(zipped); err != nil {
		return fmt.Errorf("could not unzip: %w", err)
	}

	return nil
}

// Downloads database file
func (s *geoipService) downloadIp2LocationDb() ([]byte, error) {
	if s.Token == "" {
		return nil, errors.New("no token provided")
	}

	url := "https://www.ip2location.com/download/?token=" + s.Token + "&file=" + ip2locationApiFile
	log.Println("GeoIP: Download zipped database file...")
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if len(bodyBytes) < 100000 {
		return nil, fmt.Errorf("body is too small (%d bytes): %.1000s", len(bodyBytes), bodyBytes)
	}

	return bodyBytes, nil
}

// Unzips database file
func (s *geoipService) unzipIp2LocationDb(zipped []byte) error {
	log.Println("GeoIP: Unzip database file...")

	zipReader, err := zip.NewReader(bytes.NewReader(zipped), int64(len(zipped)))
	if err != nil {
		return err
	}

	for _, f := range zipReader.File {
		if f.Name != ip2locationFilename {
			continue
		}

		dstFile, err := os.Create(path.Join(ip2locationDirectory, ip2locationFilename))
		if err != nil {
			return err
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return err
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return err
		}

		dstFile.Close()
		fileInArchive.Close()

		return nil
	}

	return fmt.Errorf("file %s not found in archive", ip2locationFilename)
}
