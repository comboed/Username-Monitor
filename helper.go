package main

import (
	"math/rand"
	"strings"
	"sync"
	"time"
	"fmt"
)

func getPreviousUsernameID(username string) string {
	for i := range usernames {
		if parts := strings.Split(usernames[i], ":"); parts[0] == username {
			return parts[1]
		}
	}
	return ""
}

func seperateUsernames() ([]string, []string) {
	var seperatedUsernames, seperatedIDs []string
	for i := range usernames {
		var parts []string = strings.Split(usernames[i], ":")
		seperatedUsernames, seperatedIDs = append(seperatedUsernames, parts[0]), append(seperatedIDs, parts[1])
	}
	return seperatedUsernames, seperatedIDs
}

func formatEndpoint(chunkedIDs []string) string {
	var joinedIDs string = strings.Join(chunkedIDs, `\",\"credential_type\":\"none\",\"token\":\"\"},{\"uid\":\"`)
	if (rotateEndpoint) {
		return fmt.Sprintf(InstagramBloksGraphQLUrl, joinedIDs) + `\",\"credential_type\":\"none\",\"token\":\"\"}]}"}}`
	}
	return fmt.Sprintf(InstagramWBloksFetchAsyncUrl, InstagramLsd, joinedIDs) + `\",\"credential_type\":\"none\",\"token\":\"\"}]}}"}`	
}

func buildUsernameGroups() {	
	var tempUsernameGroups [][]string
	var seperatedUsernames, seperatedIDs []string = seperateUsernames()
	var chunkedUsernames, chunkedIDs [][]string = chunkSlice(seperatedUsernames, BatchSize), chunkSlice(seperatedIDs, BatchSize)

	for i := range chunkedIDs {
		var group []string = []string{formatEndpoint(chunkedIDs[i])}
		for j := range chunkedIDs[i] {
			group = append(group, chunkedUsernames[i][j], chunkedIDs[i][j])
		}
		tempUsernameGroups = append(tempUsernameGroups, group)
	}
	usernameGroups = tempUsernameGroups
}

func buildRequests() {
	requests = [][]byte{}
	requestMap = make(map[string][][]byte, len(usernames))
	rand.New(rand.NewSource(time.Now().UnixNano())).Shuffle(len(sessions), func(i, j int) { sessions[i], sessions[j] = sessions[j], sessions[i] })

	var mutex sync.Mutex
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(usernames))

	for _, username := range usernames {
		go func(username string) {
			defer waitGroup.Done()
			
			var requests [][]byte
			for _, session := range sessions {
				requests = append(requests, buildRequest(strings.Split(session, ":")[4], username))
			}
			mutex.Lock()
			requestMap[username] = requests
			mutex.Unlock()

		}(strings.Split(username, ":")[0])
	}
	waitGroup.Wait()
}

func verifySessions() {
	for _, session := range sessions {
		if (synchronizeA) {
			time.Sleep(time.Second * 10)
		}
		if (!isSessionActive(strings.Split(session, ":")[4])) {
			sessions = removeString(sessions, session)
		}
	}
	createFile("./data/sessions.txt", sessions)
	buildRequests()
}