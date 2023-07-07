package zoho_cliq

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type Users struct {
	XMLName xml.Name `xml:"users"`
	Users   []User   `xml:"user"`
}

type User struct {
	ID           string `xml:"id"`
	Email        string `xml:"email"`
	Name         string `xml:"name"`
	DisplayName  string `xml:"display_name"`
	IsAdmin      bool   `xml:"is_admin"`
	IsSuperAdmin bool   `xml:"is_super_admin"`
	IsDeleted    bool   `xml:"is_deleted"`
	Timezone     string `xml:"tz"`
	AzHash       string `xml:"az_hash"`
}

func ReadUsers() (*Users, error) {

	xmlFile, err := os.Open("files/import/users.xml")
	if err != nil {
		fmt.Println("Error opening XML file:", err)
		return &Users{}, err
	}
	defer xmlFile.Close()

	// Read the XML content
	xmlData, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		fmt.Println("Error reading XML data:", err)
		return &Users{}, err
	}

	// Initialize the struct to hold the parsed data
	users := Users{}

	// Unmarshal the XML into the struct
	err = xml.Unmarshal(xmlData, &users)
	if err != nil {
		fmt.Println("Error unmarshaling XML:", err)
		return &Users{}, err
	}

	return &users, err

}
