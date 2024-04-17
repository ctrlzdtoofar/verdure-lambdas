package mdl

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type ConfirmationType string

const (
	NewUser       ConfirmationType = "NewUser"
	ResetPassword ConfirmationType = "ResetPassword"
)

// UserConfirmation matches the Rust UserConfirmation struct
type UserConfirmation struct {
	ConfirmationType ConfirmationType `json:"confirmation_type"`
	BaseUrl          string           `json:"base_url"`
	UserLoginID      int32            `json:"user_login_id"`
	Email            string           `json:"email"`
	Token            string           `json:"token"`
	ExpiresAtMillis  int64            `json:"expires_at_millis"`
}

// ExpiresAt Helper function to convert ExpiresAtMillis to a time.Time
func (uc *UserConfirmation) ExpiresAt() time.Time {
	return time.Unix(0, uc.ExpiresAtMillis*int64(time.Millisecond))
}

// JSON Creates a JSON string from a UserConfirmation object.
func (uc *UserConfirmation) JSON() string {
	b, err := json.Marshal(uc)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return ""
	}
	return string(b)
}

// ConfirmUrL Creates the URL used with email confirmation.
func (uc *UserConfirmation) ConfirmUrL() string {
	return fmt.Sprintf("%s/confirm/%s/%d/%s",
		uc.BaseUrl,
		strings.ToLower(string(uc.ConfirmationType)),
		uc.UserLoginID,
		uc.Token)
}

// UserConfirmationFromJson Unmarshalls a JSON string into the UserConfirmation instance.
func UserConfirmationFromJson(jsonStr string) (*UserConfirmation, error) {

	// Convert string to byte slice
	jsonData := []byte(jsonStr)

	// Declare a variable of type DbConnect
	var conf UserConfirmation

	// Decode the JSON data into the struct
	err := json.Unmarshal(jsonData, &conf)
	if err != nil {
		fmt.Printf("error decoding JSON, %v", err)
		return nil, err
	}

	return &conf, nil
}
