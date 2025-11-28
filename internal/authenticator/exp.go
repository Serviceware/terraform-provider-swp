package authenticator

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ExpirationTime is a simple helper function that extracts the expiration
// time claim from jwt and returns it as [time.Time].
func ExpirationTime(token string) (time.Time, error) {
	_, payloadPart, _, err := parts(token)
	if err != nil {
		return time.Time{}, err
	}

	payload, err := base64.RawURLEncoding.DecodeString(payloadPart)
	if err != nil {
		return time.Time{}, err
	}

	parsedPayload := struct {
		Exp *int64 `json:"exp"`
	}{}
	if err = json.Unmarshal(payload, &parsedPayload); err != nil {
		return time.Time{}, err
	} else if parsedPayload.Exp == nil {
		return time.Time{}, fmt.Errorf("token has no exp claim")
	}

	return time.Unix(*parsedPayload.Exp, 0), nil
}

func parts(token string) (string, string, string, error) {
	header, payloadAndSignature, hasHeader := strings.Cut(token, ".")
	payload, signature, hasPayload := strings.Cut(payloadAndSignature, ".")
	dividesSignature := strings.Contains(signature, ".")
	if !hasHeader || !hasPayload || dividesSignature {
		return "", "", "", fmt.Errorf("token has != 3 parts")
	}
	return header, payload, signature, nil
}

// IsValidIn checks whether jwt is still valid in the future at
// now + delta. If the expiration time claim cannot be read
// from jwt, false is returned.
func IsValidIn(token string, delta time.Duration) bool {
	expirationTime, err := ExpirationTime(token)
	if err != nil {
		return false
	}
	return time.Now().Add(delta).Before(expirationTime)
}
