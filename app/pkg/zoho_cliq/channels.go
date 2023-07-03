package zoho_cliq

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type Channels struct {
	XMLName  xml.Name  `xml:"channels"`
	Channels []Channel `xml:"channel"`
}

type Channel struct {
	SysArchived   bool     `xml:"sys_archived"`
	Name          string   `xml:"name"`
	ID            string   `xml:"id"`
	Creator       string   `xml:"creator"`
	OpenToAll     bool     `xml:"open_to_all"`
	Archived      bool     `xml:"archived"`
	ChatID        string   `xml:"chat_id"`
	DataDirectory string   `xml:"data_directory"`
	Description   string   `xml:"description"`
	Scope         string   `xml:"scope"`
	Members       []Member `xml:"members>member"`
	Teams         []string `xml:"teams>team"`
}

type Member struct {
	Pinned bool   `xml:"pinned"`
	ID     string `xml:"id"`
}

func ReadChannels() (*Channels, error) {

	xmlFile, err := os.Open("files/import/channels.xml")
	if err != nil {
		fmt.Println("Error opening XML file:", err)
		return &Channels{}, err
	}
	defer xmlFile.Close()

	// Read the XML content
	xmlData, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		fmt.Println("Error reading XML data:", err)
		return &Channels{}, err
	}

	// Initialize the struct to hold the parsed data
	channels := Channels{}

	// Unmarshal the XML into the struct
	err = xml.Unmarshal(xmlData, &channels)
	if err != nil {
		fmt.Println("Error unmarshaling XML:", err)
		return &Channels{}, err
	}

	return &channels, err

}
