package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	az "main/pkg/msteams"
	"main/pkg/zoho_cliq"
	"os"
	"strconv"
)

var tenantID string
var clientID string
var clientSecret string

var TeamName string
var TeamDescription string
var TeamCreateDate string
var GuestAzObjectID string
var parallelImportMessages int

func main() {

	// Load global variables from .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	tenantID = os.Getenv("tenantID")
	clientID = os.Getenv("clientID")
	clientSecret = os.Getenv("clientSecret")

	TeamName = os.Getenv("TeamName")
	TeamDescription = os.Getenv("TeamDescription")
	TeamCreateDate = os.Getenv("TeamCreateDate")
	GuestAzObjectID = os.Getenv("GuestAzObjectID")
	parallelImportMessages, err = strconv.Atoi(os.Getenv("parallelImportMessages"))

	var stateApp AzTeam
	// write User list from messages
	countMessages, _ := zoho_cliq.CollectUniqueUsers()
	fmt.Printf("Count Messages: %v\n", countMessages)

	// Get Access token
	accessToken, err := az.GetAzureTokenSecrets(tenantID, clientID, clientSecret)
	if err != nil {
		fmt.Println("Error get token", err)
	}

	// Check state file if empty create new team and save state file
	_, err = os.Stat("files/output/app-state.json")
	if err != nil {
		//Create Team with migration mode + save state file
		stateApp.TeamId, err = az.CreateTeamMigrate(accessToken, TeamName, TeamDescription, TeamCreateDate)
		if err != nil {
			fmt.Println("Error create Team: ", err)
		}
		stateApp.TeamName = TeamName
		saveState(&stateApp)
		fmt.Printf("TeamName: %v | TeamID: %v \n", stateApp.TeamName, stateApp.TeamId)

		//Start load channels from XML and save to stateFile
		channels, err := zoho_cliq.ReadChannels()
		if err != nil {
			fmt.Println("Get channel err:", err)
		}

		for _, channel := range channels.Channels {

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
				stateApp.Channel = append(stateApp.Channel, AzChannel{
					ChannelId:            "",
					ChannelName:          channel.Name,
					ChannelDescription:   channel.Description,
					ImportMessagesStatus: false,
					DataDirectories:      append([]string{}, channel.DataDirectory),
				})
			}
			saveState(&stateApp)
		}

	} else {
		stateApp, err = loadState()
	}

	CreateChannelsAndImportMessagesToChannel(&stateApp, accessToken)

	// Close migrate channel
	//az.CompleateChannelMigrate(accessToken, azTeam.TeamName, channelID)
	//}
	// Close migrate Teams
	//az.CompleateTeamMigrate(accessToken, stateApp.TeamName)

}

func CreateChannelsAndImportMessagesToChannel(stateApp *AzTeam, accessToken string) {

	expiredToken, _ := az.IsTokenExpired(accessToken)
	if expiredToken {
		// Get access token
		accessToken, _ = az.GetAzureTokenSecrets(tenantID, clientID, clientSecret)
		fmt.Println("Token updated!")
	}

	// Create Channel and Import messages to channel
	respCreateChannelCount := make(map[int]int)

	for i, st := range stateApp.Channel {
		if st.ImportMessagesStatus {
			continue
		}
		// Create channel
		respCode, channelID, err := az.CreateChannelMigrate(accessToken, stateApp.TeamId, st.ChannelName, st.ChannelDescription, TeamCreateDate)
		if err != nil {
			fmt.Println("Error create channel:", err)
			break
		}
		if st.ChannelName == "general" {
			ch_tmp, _ := az.ListChannels(accessToken, stateApp.TeamId)
			for _, ch := range ch_tmp.Channels {
				if ch.DisplayName == "General" {
					channelID = ch.ID
				}
			}
		}
		fmt.Println("chid:", channelID)
		stateApp.Channel[i].ChannelId = channelID
		saveState(stateApp)

		respCreateChannelCount[respCode]++
		fmt.Println(st.ChannelId, st.ChannelName)
		fmt.Println(respCreateChannelCount)

		for _, dataDir := range st.DataDirectories {
			fmt.Printf("Try import file: %v\n", dataDir)
			if st.ImportMessagesStatus == false {
				ImportMessages(accessToken, stateApp.TeamId, channelID, dataDir)
			}
		}

		// Mark channel is success imported
		for i := range stateApp.Channel {
			if stateApp.Channel[i].ChannelName == st.ChannelName {
				stateApp.Channel[i].ImportMessagesStatus = true
				saveState(stateApp)
			}
		}

	}

}
