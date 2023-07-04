package main

import (
	"fmt"
	az "main/pkg/msteams"
	"main/pkg/zoho_cliq"
)

var tenantID string = "e5fa314d-6060-4810-b263-abdcba14735e"
var clientID string = "02b48c2b-200c-411e-b8e8-d0c2eb709cb4"
var clientSecret string = ""

var TeamName string = "AMigration-06"
var TeamDescription string = "Some describtion"
var TeamCreateDate string = "2015-03-14T11:22:17.043Z"
var GuestAzObjectID string = "3575eae2-f2d0-4dcd-a8b6-ff8a5b50cc4a"

func main() {
	// write User list from messages
	countMessages, _ := zoho_cliq.CollectUniqueUsers()
	fmt.Printf("Count Messages: %v\n", countMessages)

	//Create Team and Channels with migration mode + save state file
	//stateApp, _ := createTeamsAndChannels()
	//fmt.Println(stateApp)
	//saveState(stateApp)

	accessToken, _ := az.GetAzureTokenSecrets(tenantID, clientID, clientSecret)
	stateApp, _ := loadState()

	for _, st := range stateApp.Channel {
		fmt.Println(st.ChannelId, st.ChannelName)
		for _, dataDir := range st.DataDirectories {
			fmt.Printf("Try import file: %v\n", dataDir)
			if st.ImportMessagesStatus == false {
				ImportMessages(accessToken, stateApp.TeamId, st.ChannelId, dataDir)
			}
		}

		// Mark channel is success imported
		for i := range stateApp.Channel {
			if stateApp.Channel[i].ChannelName == st.ChannelName {
				stateApp.Channel[i].ImportMessagesStatus = true
				saveState(stateApp)
			}
		}

		az.CompleateChannelMigrate(accessToken, stateApp.TeamId, st.ChannelId)
	}

	//os.Exit(1)
	//ImportMessages(accessToken, azTeam.TeamName, channelID, channel.DataDirectory)

	// Close migrate channel
	//az.CompleateChannelMigrate(accessToken, azTeam.TeamName, channelID)
	//}
	// Close migrate Teams
	//az.CompleateTeamMigrate(accessToken, stateApp.TeamName)

}

func createTeamsAndChannels() (AzTeam, error) {
	stateApp := AzTeam{}
	// Set up authentication
	accessToken, _ := az.GetAzureTokenSecrets(tenantID, clientID, clientSecret)
	fmt.Printf("%v, %v, %v", TeamName, TeamDescription, TeamCreateDate)

	//Start Migrate
	stateApp.TeamName = TeamName
	stateApp.TeamId, _ = az.CreateTeamMigrate(accessToken, stateApp.TeamName, TeamDescription, TeamCreateDate)
	fmt.Println("TeamID:", stateApp.TeamId)

	//zoho_cliq.ReadTeams()
	channels, err := zoho_cliq.ReadChannels()
	if err != nil {
		fmt.Println("Get channel err:", err)
		return stateApp, err
	}

	// Access the parsed data
	for _, channel := range channels.Channels {

		channelID, err := az.CreateChannelMigrate(accessToken, stateApp.TeamId, channel.Name, channel.Description, TeamCreateDate)
		if err != nil {
			fmt.Println("Error create channel:", err)
			break
		}

		// save data to State
		channelIsNotExistInState := true
		for i := range stateApp.Channel {
			if stateApp.Channel[i].ChannelName == channel.Name {
				fmt.Println("dublicate found")
				channelIsNotExistInState = false
				stateApp.Channel[i].DataDirectories = append(stateApp.Channel[i].DataDirectories, channel.DataDirectory)
			}
		}

		if channelIsNotExistInState {
			ch_tmp, _ := az.ListChannels(accessToken, stateApp.TeamId)
			for _, ch := range ch_tmp.Channels {
				if ch.DisplayName == "General" {
					channelID = ch.ID
				}
			}

			stateApp.Channel = append(stateApp.Channel, AzChannel{
				ChannelId:            channelID,
				ChannelName:          channel.Name,
				ImportMessagesStatus: false,
				DataDirectories:      append([]string{}, channel.DataDirectory),
			})
		}
		//fmt.Println(stateApp)
		fmt.Printf("%v | %v | %v | %v  \n", channelID, channel.Name, channel.Description, channel.DataDirectory)
	}

	return stateApp, nil
}
