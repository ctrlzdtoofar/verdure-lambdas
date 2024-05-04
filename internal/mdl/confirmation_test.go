package mdl

import (
	"testing"
	"time"
)

func TestUserConfirmationFromJson(t *testing.T) {
	expectedUserLoginID := int32(123)
	expectedBaseUrl := "http://localhost:3211"
	expectedEmail := "tst123@example.com"
	expectedLang := "en"
	expectedToken := "abc123"
	expectedExpiresAtMillis := int64(1672531200000) // Assuming a specific time for testing

	jsonStr := `{"confirmation_type":"NewUser","base_url":"http://localhost:3211","user_login_id":123,"email":"tst123@example.com","lang":"en","token":"abc123","expires_at_millis":1672531200000}`

	confirmation, err := UserConfirmationFromJson(jsonStr)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if confirmation.BaseUrl != expectedBaseUrl {
		t.Errorf("Expected BaseUrl %s, got %s", expectedBaseUrl, confirmation.BaseUrl)
	}

	if confirmation.UserLoginID != expectedUserLoginID {
		t.Errorf("Expected UserLoginID %d, got %d", expectedUserLoginID, confirmation.UserLoginID)
	}

	if confirmation.Email != expectedEmail {
		t.Errorf("Expected Email %s, got %s", expectedEmail, confirmation.Email)
	}

	if confirmation.Lang != expectedLang {
		t.Errorf("Expected Lang %s, got %s", expectedLang, confirmation.Lang)
	}

	if confirmation.Token != expectedToken {
		t.Errorf("Expected Token %s, got %s", expectedToken, confirmation.Token)
	}

	if confirmation.ExpiresAtMillis != expectedExpiresAtMillis {
		t.Errorf("Expected ExpiresAtMillis %d, got %d", expectedExpiresAtMillis, confirmation.ExpiresAtMillis)
	}

	expectedExpiresAt := time.Unix(0, expectedExpiresAtMillis*int64(time.Millisecond))
	if confirmation.ExpiresAt() != expectedExpiresAt {
		t.Errorf("Expected ExpiresAt %v, got %v", expectedExpiresAt, confirmation.ExpiresAt())
	}
}

// TestUserConfirmationFromJson_Fail tests an expected failure due to incorrect JSON format
func TestUserConfirmationFromJson_Fail(t *testing.T) {
	jsonStr := `{"confirmation_type":"NewUser","user_login_id":"shouldafaildueToIncorrectType","email":"example@example.com","token":"abc123","expires_at_millis":1672531200000}`
	confirmation, err := UserConfirmationFromJson(jsonStr)
	if err == nil {
		t.Errorf("Expected an error due to incorrect JSON format, but did not get one")
	}

	if confirmation != nil {
		t.Errorf("Expected no user confirmation due to incorrect JSON format, but got one")
	}
}
