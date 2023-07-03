package msteams

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"log"
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

func GetGraphClient(tenantID string, clientID string, clientSecret string) (*msgraphsdk.GraphServiceClient, string, error) {
	//cred, err := azidentity.NewDeviceCodeCredential(&azidentity.DeviceCodeCredentialOptions{
	//	TenantID: tenantID,
	//	ClientID: clientID,
	//	UserPrompt: func(ctx context.Context, message azidentity.DeviceCodeMessage) error {
	//		fmt.Println(message.Message)
	//		return nil
	//	},
	//})

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
