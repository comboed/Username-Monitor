package main

import (
	"github.com/valyala/fasthttp"
	"crypto/tls"
	"strings"
	"bufio"
	"time"
	"fmt"
	"log"
)

func autoClaimer(client *fasthttp.Client, channel chan []string) {
	var request *fasthttp.Request = createRequest("POST", true)
	var response *fasthttp.Response = fasthttp.AcquireResponse()
	for {
		var group []string = <- channel
		request.SetRequestURI(group[0])
		request.Header.Set("IG-U-DS-USER-ID", randomIntString(11))
		request.Header.Set("IG-INTENDED-USER-ID", randomIntString(11))
		request.Header.Set("Cookie", "datr=" + instagramDatr + "; csrftoken=" + InstagramCsrf + "; ds_user_id=" + randomIntString(11))

		client.Do(request, response)
		var body, _ = response.BodyGunzip()
		var bodyString string = string(body)

		if (strings.Contains(bodyString, "credential_type")) {
			for i := 1; i < len(group); i += 2 {
				if (!strings.Contains(bodyString, `"` + group[i] + `\`) && strings.Count(bodyString, `"` + group[i + 1] + `\"`) > 2 && !synchronizeA) {
					synchronizeA = true
					requests = requestMap[group[i]]
					time.Sleep(time.Second * 3)

					requests = [][]byte{}
					isUsernameClaimed(group[i])
					time.Sleep(time.Second * 2)
					synchronizeA = false
				}
			}
			attempts += 1
		} else if len(bodyString) > 0 {
			rl += 1
		}
		request.Header.DelAllCookies()
	}
}

func voidMonitor(client *fasthttp.Client, channel chan []string) {
	var request *fasthttp.Request = createRequest("POST", true)
	var response *fasthttp.Response = fasthttp.AcquireResponse()
	for {
		var username []string = <- channel
		request.SetRequestURI(username[0])
		request.Header.Set("IG-U-DS-USER-ID", randomIntString(11))
		request.Header.Set("IG-INTENDED-USER-ID", randomIntString(11))
		request.Header.Set("Cookie", "datr=" + instagramDatr + "; csrftoken=" + InstagramCsrf + "; ds_user_id=" + randomIntString(11))

		client.Do(request, response)
		var body, _ = response.BodyGunzip()
		var bodyString string = string(body)

		if (strings.Contains(bodyString, `"` + username[1] + `\`) && strings.Count(bodyString, `"` + username[2] + `\"`) > 2 && !synchronizeB) {
			synchronizeB = true
			usernameUnvoided(username[1:])

			time.Sleep(time.Second * 2)
			synchronizeB = false
		}
		request.Header.DelAllCookies()
	}
}

func usernameSpammer(channel chan []byte) {
	var response *fasthttp.Response = fasthttp.AcquireResponse()
	var connection *tls.Conn = createTLSConnection()
	go connectionRefresher(&connection)
	for {
		connection.Write(<- channel)
		response.Read(bufio.NewReaderSize(connection, 8096))
		break
	}
}

func isUsernameClaimed(username string) {
	var usernameID, previousUsernameID string = getUsernameID(username), getPreviousUsernameID(username)
	for _, session := range sessions {
		if (strings.Contains(session, usernameID)) {
			claimedUsername(username, session)
			break
		}
	}
	checkUsernameStatus(username, usernameID, previousUsernameID)
}

func claimedUsername(username, session string) {
	log.Printf("Autoclaimed \x1b[32m@%s\x1b[39m after \x1b[32m%s\x1b[39m attempts! %s \n\n", username, formatNumber(attempts), strings.Repeat(" ", 35))
	sendDiscordWebhook(DiscordAutoclaimedWebhook, fmt.Sprintf(DiscordAutoclaimedTemplate, username, username, time.Now().Format(time.RFC3339)))

	var parts []string = strings.Split(session, ":")
	appendFile("./data/logs/" + username + ".log", fmt.Sprintf("Username: %s\nPassword: %s\nEmail: %s\nEmail Password: %s\nSession ID: %s\nTimestamp: %s", username, parts[1], parts[2], parts[3], parts[4], time.Now().Format(time.RFC1123)))

	sessions = removeString(sessions, session)
	createFile("./data/sessions.txt", sessions)
	buildRequests()
}

func checkUsernameStatus(username, usernameID, previousUsernameID string) {
	if (usernameID == "Unavailable") {
		log.Printf("@%s has been voided %s \n\n", username, strings.Repeat(" ", 45))
		sendDiscordWebhook(DiscordVoidWebhook, fmt.Sprintf(DiscordVoidTemplate, "`" + previousUsernameID + "`", username, username, time.Now().Format(time.RFC3339)))
		voidUsernames = append(voidUsernames, username + ":" + previousUsernameID)
		createFile("./data/void_usernames.txt", voidUsernames)
	} else {
		log.Printf("@%s has been swapped [Old ID: %s | New ID: %s] %s \n\n", username, previousUsernameID, usernameID, strings.Repeat(" ", 25))
		sendDiscordWebhook(DiscordMissedWebhook, fmt.Sprintf(DiscordMissedTemplate, "`" + previousUsernameID + "`", "`" + usernameID + "`", username, username, time.Now().Format(time.RFC3339)))
		usernames = append(usernames, username + ":" + usernameID)
	}
	usernames = removeString(usernames, username + ":" + previousUsernameID)
	createFile("./data/usernames.txt", usernames)
	buildUsernameGroups()
}

func usernameUnvoided(username []string) {
	log.Printf("@%s has been unvoided %s \n\n", username[0], strings.Repeat(" ", 45))
	sendDiscordWebhook(DiscordVoidWebhook, fmt.Sprintf(DiscordUnvoidTemplate, "`" + username[1] + "`", username[0], username[0], time.Now().Format(time.RFC3339)))

	usernames = append(usernames, username[0] + ":" + getUsernameID(username[0]))
	createFile("./data/usernames.txt", usernames)

	voidUsernames = removeString(voidUsernames, username[0] + ":" + username[1])
	createFile("./data/void_usernames.txt", voidUsernames)

	buildUsernameGroups()
	buildRequests()
}