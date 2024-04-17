package settings

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"os"
)

// GetEnv retrieves environment variables or returns a default value
func GetEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

// GetNoReplySecret retrieves a specific secret from AWS Secrets Manager and extracts the "noreply" email
// address from the JSON data stored within the secret. This function is designed to securely access
// and parse configuration data required by applications.
//
// Parameters:
//   - ctx: A context.Context to propagate deadlines, cancellation signals, and other request-scoped values
//     across API boundaries and between processes.
//   - secretName: The name or identifier of the secret in AWS Secrets Manager. This secret should contain
//     a JSON object with a "noreply" key.
//
// Returns:
//   - string: The "noreply" email address extracted from the secret's JSON data.
//   - error: An error object that reports failure scenarios during the execution of the function. Potential
//     errors include failure to load AWS SDK configuration, failure to retrieve the secret, failure to
//     parse the JSON data, and the absence of the "noreply" key in the JSON data.
func GetNoReplySecret(ctx context.Context, secretName string) (string, error) {

	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", err
	}

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(sdkConfig)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"),
	}

	// Fetch the secret value.
	result, err := svc.GetSecretValue(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve secret: %w", err)
	}

	// Check if the secret comes as a string (standard case).
	if result.SecretString == nil {
		return "", fmt.Errorf("secret string is nil")
	}

	// Parse the JSON to extract the "noreply" field.
	var secretMap map[string]string
	err = json.Unmarshal([]byte(*result.SecretString), &secretMap)
	if err != nil {
		return "", fmt.Errorf("failed to parse secret string: %w", err)
	}

	// Extract the "noreply" email address from the map.
	email, ok := secretMap["noreply"]
	if !ok {
		return "", fmt.Errorf("email not found in secret")
	}

	return email, nil
}
