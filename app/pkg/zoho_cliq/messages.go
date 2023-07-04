package zoho_cliq

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Messages struct {
	XMLName xml.Name  `xml:"messages"`
	Message []Message `xml:"message"`
}

type Message struct {
	IsRead    string        `xml:"is_read"`
	ID        string        `xml:"id"`
	Time      string        `xml:"time"`
	Type      string        `xml:"type"`
	Revision  string        `xml:"revision"`
	Timestamp string        `xml:"timestamp"`
	Sender    Sender        `xml:"sender"`
	Content   string        `xml:"content>type,omitempty"`
	Text      string        `xml:"content>text,omitempty"`
	Data      string        `xml:"content>data>users" xml:"innerxml"`
	Source    MessageSource `xml:"message_source"`
	Mentions  string        `xml:"-"`
}

type Sender struct {
	Name  string `xml:"name"`
	ID    string `xml:"id"`
	Alias Alias  `xml:"alias"`
}

type Alias struct {
	Name  string `xml:"name"`
	Image string `xml:"image"`
}

//type Content struct {
//	Text string `xml:"text"`
//	Card Card   `xml:"card"`
//}

type Card struct {
	Theme int    `xml:"theme"`
	Title string `xml:"title"`
}

type MessageSource struct {
	StoreAppID int `xml:"store_appid"`
}

func ReadMessagesFromFile(filePath string) (*Messages, error) {

	fixXML(filePath)

	xmlFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening XML file:", err)
		return &Messages{}, err
	}
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)

	decoder := xml.NewDecoder(bytes.NewReader(byteValue))
	decoder.Strict = false
	decoder.AutoClose = xml.HTMLAutoClose
	decoder.Entity = xml.HTMLEntity

	var messages Messages
	err = decoder.Decode(&messages)
	if err != nil {
		fmt.Println("Error unmarshaling XML:", err)
		return &Messages{}, err
	}

	return &messages, err
}

func CollectUniqueUsers() (countMessages int, err error) {
	var users []User
	uniqueMap := make(map[string]bool)

	channels, err := ReadChannels()
	if err != nil {
		fmt.Println("Get channel err:", err)
	}

	for _, channel := range channels.Channels {
		//fmt.Println()
		//fmt.Printf("%v | %v | %v  \n", channel.Name, channel.Description, channel.DataDirectory)

		directory := fmt.Sprintf("files/import/messages/channels/%v/", channel.DataDirectory)

		err = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				// Handle the error, if any
				return err
			}
			// Skip directories
			if info.IsDir() {
				return nil
			}

			//fmt.Println(path)
			messages, err := ReadMessagesFromFile(path)

			for _, msg := range messages.Message {
				countMessages++
				if !uniqueMap[msg.Sender.ID] {
					user := User{
						ID:   msg.Sender.ID,
						Name: msg.Sender.Name,
					}
					// Append the user to the slice
					users = append(users, user)

					// Mark the ID as seen
					uniqueMap[msg.Sender.ID] = true
				}
				//fmt.Println(msg.Sender.Name, msg.Sender.ID)
			}

			return nil
		})
	}

	file, err := os.Create("files/output/listUniqUsersFromMessages.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}

	for _, user := range users {
		// Data to write
		data := fmt.Sprintf("%v - %v \n", user.Name, user.ID)

		// Write the data to the file
		_, err = file.WriteString(data)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}
	fmt.Println("Unique users saved to file: files/output/listUniqUsersFromMessages.txt")

	file.Close()
	return countMessages, nil

}
