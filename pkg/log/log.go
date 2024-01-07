package log

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const (
	deadlineTimeout = 5
)

var ZL *zerolog.Logger

func init() {

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	ZL = &logger

	// Add http logger writer to log collector
	if viper.GetBool("IS_HTTP_LOGGER") {
		httpWriter := HTTPWriter{}

		writer := diode.NewWriter(
			httpWriter,
			10000,
			10*time.Millisecond, func(missed int) {
				log.Printf("logger dropped %d messages", missed)
			})

		logger = zerolog.New(writer)
		ZL = &logger
	}

}

// HTTPWriter wraps Write method.
type HTTPWriter struct {
	io.Writer
}

func (HTTPWriter) Write(bs []byte) (n int, err error) {
	/*get Logstash address.*/
	url := fmt.Sprintf("https://%s:%s", viper.GetString("LOGSTASH_HOST"), viper.GetString("LOGSTASH_PORT"))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bs))
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: deadlineTimeout * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		log.Error().Msgf("%v", err)
		return 0, fmt.Errorf("failed send data to Logstash [%s]", err)
	}
	defer res.Body.Close()
	return len(bs), nil
}
