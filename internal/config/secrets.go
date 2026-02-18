package config

import (
	"context"
	"encoding/base64"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// Secrets ...
type Secrets struct {
	svc *secretsmanager.Client
}

// NewSecrets ...
func NewSecrets(svc *secretsmanager.Client) *Secrets {
	return &Secrets{
		svc: svc,
	}
}

// GoogleAdminEmail ...
func (s *Secrets) GoogleAdminEmail(secretArn string) (string, error) {
	if len([]rune(secretArn)) == 0 {
		return s.getSecret(context.Background(), "SSOSyncGoogleAdminEmail")
	}
	return s.getSecret(context.Background(), secretArn)
}

// SCIMAccessToken ...
func (s *Secrets) SCIMAccessToken(secretArn string) (string, error) {
	if len([]rune(secretArn)) == 0 {
		return s.getSecret(context.Background(), "SSOSyncSCIMAccessToken")
	}
	return s.getSecret(context.Background(), secretArn)
}

// SCIMEndpointURL ...
func (s *Secrets) SCIMEndpointURL(secretArn string) (string, error) {
	if len([]rune(secretArn)) == 0 {
		return s.getSecret(context.Background(), "SSOSyncSCIMEndpointURL")
	}
	return s.getSecret(context.Background(), secretArn)
}

// GoogleCredentials ...
func (s *Secrets) GoogleCredentials(secretArn string) (string, error) {
	if len([]rune(secretArn)) == 0 {
		return s.getSecret(context.Background(), "SSOSyncGoogleCredentials")
	}
	return s.getSecret(context.Background(), secretArn)
}

// Region ...
func (s *Secrets) Region(secretArn string) (string, error) {
	if len([]rune(secretArn)) == 0 {
		return s.getSecret(context.Background(), "SSOSyncRegion")
	}
	return s.getSecret(context.Background(), secretArn)
}

// IdentityStoreID ...
func (s *Secrets) IdentityStoreID(secretArn string) (string, error) {
	if len([]rune(secretArn)) == 0 {
		return s.getSecret(context.Background(), "IdentityStoreID")
	}
	return s.getSecret(context.Background(), secretArn)
}

func (s *Secrets) getSecret(ctx context.Context, secretKey string) (string, error) {
	r, err := s.svc.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretKey),
		VersionStage: aws.String("AWSCURRENT"),
	})

	if err != nil {
		return "", err
	}

	var secretString string

	if r.SecretString != nil {
		secretString = *r.SecretString
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(r.SecretBinary)))
		l, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, r.SecretBinary)
		if err != nil {
			return "", err
		}
		secretString = string(decodedBinarySecretBytes[:l])
	}

	return secretString, nil
}
