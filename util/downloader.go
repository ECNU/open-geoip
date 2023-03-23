package util

import (
	"os"
	"time"

	"github.com/ECNU/open-geoip/g"
	geoip "github.com/pieterclaerhout/go-geoip/v2"
	"github.com/toolkits/pkg/logger"
)

const (
	DefaultDownloadTimeout = 3
	DefaultTargetFilePath  = "./"
)

func AutoDownloadMaxmindDatabase(config g.AutoDownloadConfig) (string, error) {
	if config.MaxmindLicenseKey == "" {
		config.MaxmindLicenseKey = os.Getenv("MAXMIND_LICENSE_KEY")
	}

	if config.Timeout == 0 {
		config.Timeout = DefaultDownloadTimeout
	}

	if config.TargetFilePath == "" {
		config.TargetFilePath = DefaultTargetFilePath
	}

	dbPath := config.TargetFilePath + "GeoLite2-City.mmdb"

	downloader := geoip.NewDatabaseDownloader(config.MaxmindLicenseKey, dbPath, time.Duration(config.Timeout)*time.Minute)

	logger.Debug("Checking if the database needs updating")

	localChecksum, err := downloader.LocalChecksum()
	if err != nil {
		return dbPath, err
	}

	remoteChecksum, err := downloader.RemoteChecksum()
	if err != nil {
		return dbPath, err
	}

	logger.Debug("Local checksum: ", localChecksum)
	logger.Debug("Remote checksum:", remoteChecksum)

	shouldDownload, err := downloader.ShouldDownload()
	if err != nil {
		return dbPath, err
	}

	if !shouldDownload {
		logger.Debug("Database is up-to-date, no download needed")
		return dbPath, nil
	}

	logger.Info("Database not found or outdated, downloading")

	if err := downloader.Download(); err != nil {
		return dbPath, err
	}

	logger.Info("Database downloaded succesfully")
	return dbPath, nil
}
