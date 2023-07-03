package zoho_cliq

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type Team struct {
	Name        string   `xml:"name"`
	Description string   `xml:"description"`
	Members     []string `xml:"members>member"`
}

type Teams struct {
	XMLName xml.Name `xml:"teams"`
	Teams   []Team   `xml:"team"`
}

func ReadTeams() (*Teams, error) {

	xmlFile, err := os.Open("files/import/teams.xml")
	if err != nil {
		fmt.Println("Error opening XML file:", err)
		return &Teams{}, err
	}
	defer xmlFile.Close()

	// Read the XML content
	xmlData, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		fmt.Println("Error reading XML data:", err)
		return &Teams{}, err
	}

	// Initialize the struct to hold the parsed data
	teams := Teams{}

	// Unmarshal the XML into the struct
	err = xml.Unmarshal(xmlData, &teams)
	if err != nil {
		fmt.Println("Error unmarshaling XML:", err)
		return &Teams{}, err
	}

	// Access the parsed data
	for _, team := range teams.Teams {
		fmt.Println("Team Name:", team.Name)
		fmt.Println("Description:", team.Description)
		fmt.Println("Members:")
		for _, member := range team.Members {
			fmt.Println("-", member)
		}
		fmt.Println()
	}

	return &teams, err

}
