package httphandler

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/repositories"
	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog/log"
)

const SignatureHeader = "x-signature"
const TimestampHeader = "x-timestamp"

type HmacAuthenticator struct {
	repo *repositories.HooksRepo
}

func (h *HmacAuthenticator) Authenticate(ctx *echo.Context) error {
	// for hamc authentication first we check if headers are present
	headersMap := ctx.Request().Header
	// first we check required headers are present or not
	for _, header := range []string{SignatureHeader, TimestampHeader} {
		if headersMap.Get(header) == "" {
			log.Error().Msgf("required header is not present: %s", header)
			return echo.NewHTTPError(http.StatusUnprocessableEntity, "required header is not present:"+header)
		}
	}
	// get timestamp from header and check it is not older than 5 minutes
	timestamp := headersMap.Get(TimestampHeader)
	if timestamp == "" {
		log.Error().Msg("missing timestamp header")
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "invalid timestamp header")
	}

	timestampInt, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		log.Error().Msgf("invalid timestamp header: %s", timestamp)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "invalid timestamp header")
	}

	requestTime := time.UnixMilli(timestampInt)
	diff := time.Since(requestTime)
	if diff < 0 {
		diff = -diff
	}

	if diff > 5*time.Minute {
		log.Error().Msgf("request timestamp out of allowed skew: %s", diff)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "request timestamp is too old or too far in the future")
	}
	// Now we get the path parameeter and fetch the secret from database
	hooksname := ctx.Param("id")

	hook, err := h.repo.FetchHook(hooksname)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "hook not found")
	}
	secret := []byte(hook.Secret)
	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		log.Error().Msgf("invalid request body: %s", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	// validate the hash
	calculatedSignature := hmac.New(sha256.New, secret)
	calculatedSignature.Write([]byte(timestamp))
	calculatedSignature.Write(body)

	if hex.EncodeToString(calculatedSignature.Sum(nil)) != headersMap.Get(SignatureHeader) {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "invalid signature")
	}

	ctx.Request().Body = io.NopCloser(bytes.NewBuffer(body))
	return nil

}

func NewHmacAuthenticator(dbclient *client.PostgresClient) *HmacAuthenticator {
	return &HmacAuthenticator{
		repo: repositories.NewHooksRepo(dbclient),
	}
}
