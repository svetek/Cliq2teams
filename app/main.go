package main

import (
	"fmt"
	az "main/pkg/msteams"
	"main/pkg/zoho_cliq"
)

var tenantID string = "e5fa314d-6060-4810-b263-abdcba14735e"
var clientID string = "02b48c2b-200c-411e-b8e8-d0c2eb709cb4"
var clientSecret string = "MX28Q~aA0QFOe80zpUIEYTX5jhpIfIcK2obKwdmq"

var TeamName string = "AMigration-01"
var TeamDescription string = "Some describtion"
var TeamCreateDate string = "2015-03-14T11:22:17.043Z"
var GuestAzObjectID string = "3575eae2-f2d0-4dcd-a8b6-ff8a5b50cc4a"

func main() {

	// write User list from messages
	zoho_cliq.CollectUniqueUsers()

	// Set up authentication
	accessToken, _ := az.GetAzureTokenSecrets(tenantID, clientID, clientSecret)
	fmt.Printf("%v, %v, %v", TeamName, TeamDescription, TeamCreateDate)

	//Start Migrate
	teamID, _ := az.CreateTeamMigrate(accessToken, TeamName, TeamDescription, TeamCreateDate)

	//zoho_cliq.ReadTeams()
	channels, err := zoho_cliq.ReadChannels()
	if err != nil {
		fmt.Println("Get channel err:", err)
	}

	// Access the parsed data
	for _, channel := range channels.Channels {
		fmt.Println()
		fmt.Printf("%v | %v | %v  \n", channel.Name, channel.Description, channel.DataDirectory)
		channelID, err := az.CreateChannelMigrate(accessToken, teamID, channel.Name, channel.Description, TeamCreateDate)
		if err != nil {
			fmt.Println("Error create channel:", err)
			break
		}

		ImportMessages(accessToken, teamID, channelID, channel.DataDirectory)

		// Close migrate channel
		az.CompleateChannelMigrate(accessToken, teamID, channelID)
	}
	// Close migrate Teams
	az.CompleateTeamMigrate(accessToken, teamID)

}
