package main

import (
	"fmt"
	"log"
	az "main/pkg/msteams"
	"main/pkg/zoho_cliq"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

func convertDateTimeForImportFormat(datetimeStr string) (string, error) {
	// Input datetime string
	//datetimeStr := "2023-03-15 14:02:04"

	// Parse the datetime string
	layout := "2006-01-02 15:04:05"
	datetime, err := time.Parse(layout, datetimeStr)
	if err != nil {
		fmt.Println("Error parsing datetime:", err)
		return "", err
	}

	// Format the datetime in the desired output format
	outputFormat := "2006-01-02T15:04:05.000Z"
	outputStr := datetime.UTC().Format(outputFormat)

	return outputStr, err // Output: "2023-03-15T14:02:04.000Z"
}

func FindUserById(users *zoho_cliq.Users, id string) zoho_cliq.User {
	var userFound zoho_cliq.User
	foundUser := false

	for _, user := range users.Users {
		//fmt.Println("Users:", user.ID, " | ", id)
		if user.ID == id {
			foundUser = true
			userFound = user
			break
		}
	}

	if foundUser {
		return userFound
	} else {
		userFound := zoho_cliq.User{
			ID:          "0",
			Email:       "",
			Name:        "Guest",
			DisplayName: "Guest",
			AzHash:      GuestAzObjectID,
		}
		return userFound
	}
}

func ImportMessages(accessToken string, teamID string, channelID string, dataDir string, stateApp *AzTeam) ([]StatusImportedMessages, error) {

	var respCodes []StatusImportedMessages
	logFile, err := os.OpenFile("files/output/import-message.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	users, err := zoho_cliq.ReadUsers()
	if err != nil {
		fmt.Println("Error Read Users")
	}

	directory := fmt.Sprintf("files/import/messages/channels/%v/", dataDir)

	err = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Handle the error, if any
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Process the file
		fmt.Println(path)
		messages, err := zoho_cliq.ReadMessagesFromFile(path)

		var wg sync.WaitGroup
		var counter int
		var msgCount int

		var mutex sync.Mutex
		responseCounts := make(map[int]int)

		// Access the converted struct data
		for _, msg := range messages.Message {
			var textMessage string = " "
			var imageBase64 string = ""
			hostedContents := []map[string]interface{}{}

			expiredToken, _ := az.IsTokenExpired(accessToken)
			if expiredToken {
				// Get access token
				accessToken, _ = az.GetAzureTokenSecrets(tenantID, clientID, clientSecret)
				fmt.Println("Token updated when start try import message")
			}

			dateTimeMessage, _ := convertDateTimeForImportFormat(msg.Timestamp)
			if msg.Type == "text" || msg.Type == "file" {
				sendUser := FindUserById(users, msg.Sender.ID)

				if msg.Text != "" {
					textMessage = msg.Text
				}
				if msg.Type == "file" {
					bytes, _ := strconv.ParseInt(msg.FileSize, 10, 64)
					if msg.FileType == "image/jpeg" || msg.FileType == "image/png" {
						imageBase64, _ = az.DownloadUrlToBase64(msg.FileUrl)
						hostedContents = []map[string]interface{}{
							{
								"@microsoft.graph.temporaryId": "1",
								"contentBytes":                 imageBase64,
								"contentType":                  msg.FileType,
							},
						}
						textMessage = fmt.Sprintf("%v <a href=\"%v\">%v [%vMb]</a><br><img src=\"../hostedContents/1/$value\" style=\"vertical-align:bottom; width:600px; height: auto;\"> \n", msg.Comment, msg.FileUrl, msg.FileName, fmt.Sprintf("%.2f", float64(bytes)/(1024*1024)))
					} else {
						textMessage = fmt.Sprintf("%v <a href=\"%v\">%v [%vMb]</a> \n", msg.Comment, msg.FileUrl, msg.FileName, fmt.Sprintf("%.2f", float64(bytes)/(1024*1024)))

					}
				}

				//Type of user. Possible values are: aadUser, onPremiseAadUser, anonymousGuest, federatedUser, personalMicrosoftAccountUser, skypeUser, phoneUser, unknownFutureValue and emailUser.
				payload := map[string]interface{}{

					"body": map[string]interface{}{
						"contentType": "html",
						"content":     textMessage,
					},
					"from": map[string]interface{}{
						"user": map[string]interface{}{
							"id":               sendUser.AzHash,
							"displayName":      sendUser.Name,
							"userIdentityType": "aadUser",
						},
					},
					"createdDateTime": dateTimeMessage,
				}
				// add image file as attachments
				if len(hostedContents) != 0 {
					payload["hostedContents"] = hostedContents
					//fmt.Println(payload)
				}
				// Run Import Messages and check status code
				wg.Add(1)
				counter++
				msgCount++
				go func() {
					condition := false
					for ok := true; ok; ok = condition {
						respCode, respCause, err := az.PushMessageMigrate(accessToken, teamID, channelID, payload)

						mutex.Lock()
						responseCounts[respCode]++
						mutex.Unlock()

						if respCode == 429 {
							condition = true
							time.Sleep(3 * time.Second)
						} else if respCode == 201 {
							condition = false
							//} else if respCode == 403 {
							//	condition = true
							//	time.Sleep(1 * time.Second)
						} else if respCode == 401 {
							os.Exit(1)
							//condition = true
						} else if respCode == 405 {
							fmt.Println("Got 405 error, please restart it")
							time.Sleep(3 * time.Second)
							os.Exit(1)
							//condition = true
						} else if respCode == 412 {
							fmt.Println("Got 412 error: ", payload)
							condition = false
						} else {
							condition = false
						}

						logLine := fmt.Sprintf("%v | %v | %v | %v | %v \n %v | %v", respCode, teamID, channelID, payload, err, dataDir, respCause)
						_, err = logFile.WriteString(logLine)

					}
					wg.Done()
				}()

				if counter == parallelImportMessages {
					counter = 0
					wg.Wait()

					fmt.Println("Run count GoRoutine: ", parallelImportMessages, " msg loaded: ", msgCount)
					fmt.Println("Import messages response:", responseCounts)
					time.Sleep(3 * time.Second)
				}
			}
		}
		// Save response code
		tmpResp := StatusImportedMessages{
			FileName:   path,
			RespStatus: responseCounts,
		}
		for i := range stateApp.Channel {
			if stateApp.Channel[i].ChannelId == channelID {
				stateApp.Channel[i].ImportedFiles = append(stateApp.Channel[i].ImportedFiles, tmpResp)
				saveState(stateApp)
			}
		}

		logFile.Close()
		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
	}
	return respCodes, nil
}
