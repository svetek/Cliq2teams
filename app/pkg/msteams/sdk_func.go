package msteams

import (
	"context"
	"fmt"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

func getTeams(graphClient *msgraphsdk.GraphServiceClient) {

	result, err := graphClient.Me().JoinedTeams().Get(context.Background(), nil)
	if err != nil {
		fmt.Printf("Error getting teams: %v\n", err)
	}
	if result == nil {
		fmt.Println("Result nil")
		return
	}

	for _, team := range result.GetValue() {
		teamID := *team.GetId()
		fmt.Println("Team Name:", *team.GetDisplayName(), "TeamID:", teamID)
	}
}

func getChannels(graphClient *msgraphsdk.GraphServiceClient, teamId string) {

	result, err := graphClient.Teams().ByTeamId(teamId).Channels().Get(context.Background(), nil)
	if err != nil {
		fmt.Printf("Error getting Channels: %v\n", err)
	}
	if result == nil {
		fmt.Println("Result nil")
		return
	}

	//19:aaa1f32ab3d7c4af29fbea89be6db1490@thread.tacv2

	for _, chanels := range result.GetValue() {
		fmt.Println("Channel Name:", *chanels.GetDisplayName(), "ChannelID:", *chanels.GetId())
	}
}

func getIncomingChannels(graphClient *msgraphsdk.GraphServiceClient, teamId string) {

	fmt.Println(teamId)
	result, err := graphClient.Teams().ByTeamId(teamId).IncomingChannels().Get(context.Background(), nil)
	if err != nil {
		fmt.Printf("Error getting Incoming channels: %v\n", err)
	}
	if result == nil {
		fmt.Println("Result nil")
		return
	}

	for _, chanels := range result.GetValue() {
		fmt.Println("Channel Name:", *chanels.GetDisplayName(), "TeamID:", *chanels.GetId())
	}
}

func sendMessages(graphClient *msgraphsdk.GraphServiceClient, teamId string) {

	requestBody := graphmodels.NewChatMessage()
	body := graphmodels.NewItemBody()
	content := "Hello World"
	body.SetContent(&content)
	requestBody.SetBody(body)

	result, err := graphClient.Teams().ByTeamId(teamId).Channels().ByChannelId("19:aaa1f32ab3d7c4af29fbea89be6db1490@thread.tacv2").Messages().Post(context.Background(), requestBody, nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if result == nil {
		fmt.Println("Result nil")
	}

}

func printOdataError(err error) {
	switch err.(type) {
	case *odataerrors.ODataError:
		typed := err.(*odataerrors.ODataError)
		fmt.Printf("error:", typed.Error())
		if terr := typed.GetError(); terr != nil {
			fmt.Printf("code: %s", *terr.GetCode())
			fmt.Printf("msg: %s", *terr.GetMessage())
		}
	default:
		fmt.Printf("%T > error: %#v", err, err)
	}
}
