package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type AzTeam struct {
	TeamId   string
	TeamName string
	Channel  []AzChannel
}

type AzChannel struct {
	ChannelId            string
	ChannelName          string
	ChannelDescription   string
	ImportMessagesStatus bool
	DataDirectories      []string
}

func saveState(state *AzTeam) {

	jsonData, err := json.Marshal(state)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create("files/output/app-state.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("state is safe to JSON file:")

}

func loadState() (AzTeam, error) {
	file, err := os.Open("files/output/app-state.json")
	if err != nil {
		return AzTeam{}, err
	}
	defer file.Close()

	jsonData, err := ioutil.ReadAll(file)
	if err != nil {
		return AzTeam{}, err
	}

	var team AzTeam
	err = json.Unmarshal(jsonData, &team)
	if err != nil {
		return AzTeam{}, err
	}

	return team, nil
}
