package msteams

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/golang-jwt/jwt/v4"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"log"
	"time"
)

func GetAzureTokenSecrets(tenantID string, clientID string, clientSecret string) (string, error) {

	cred, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)

	if err != nil {
		fmt.Printf("Error creating credentials: %v\n", err)
		return "", err
	}

	tokenRequest := policy.TokenRequestOptions{
		Scopes:   []string{"https://graph.microsoft.com/.default"},
		TenantID: tenantID,
	}

	token, err := cred.GetToken(context.Background(), tokenRequest)
	if err != nil {
		log.Fatalf("Failed to get token: %v", err)
	}

	return token.Token, err

}

func IsTokenExpired(tokenString string) (bool, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return false, fmt.Errorf("error parsing token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, fmt.Errorf("invalid token claims")
	}

	expirationTime, ok := claims["exp"].(float64)
	if !ok {
		return false, fmt.Errorf("invalid expiration time in token")
	}

	expiration := time.Unix(int64(expirationTime), 0)
	currentTime := time.Now().Add(-2 * time.Minute)

	//fmt.Println(expiration, " | ", currentTime)
	return expiration.Before(currentTime), nil
}

func GetGraphClient(tenantID string, clientID string, clientSecret string) (*msgraphsdk.GraphServiceClient, string, error) {

	cred, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)

	if err != nil {
		fmt.Printf("Error creating credentials: %v\n", err)
		return nil, "", err
	}

	tokenRequest := policy.TokenRequestOptions{
		Scopes:   []string{"https://graph.microsoft.com/.default"},
		TenantID: tenantID,
	}

	token, err := cred.GetToken(context.Background(), tokenRequest)
	if err != nil {
		log.Fatalf("Failed to get token: %v", err)
	}

	accessToken := token.Token

	graphClient, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, nil)
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return nil, "", err
	}

	return graphClient, accessToken, nil
}
