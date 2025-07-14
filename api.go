package main

import (
	"github.com/valyala/fasthttp"
	"github.com/valyala/fastjson"
	"strings"
	"fmt"
	"os"
)

func isSessionActive(sessionID string) bool {
	var request *fasthttp.Request = createRequest("HEAD", false)
	var response *fasthttp.Response = fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(request)
	defer fasthttp.ReleaseResponse(response)
	
	request.SetRequestURI("https://www.instagram.com/")

	request.Header.Set("Cookie", "sessionid=" + sessionID)

	fasthttp.Do(request, response)

	return response.StatusCode() == 200
}

func getUsernameID(username string) string {
	var request *fasthttp.Request = createRequest("GET", false)
	var response *fasthttp.Response = fasthttp.AcquireResponse()
	
	defer fasthttp.ReleaseRequest(request)
	defer fasthttp.ReleaseResponse(response)

	request.SetRequestURI(`https://i.instagram.com/graphql_www?doc_id=6881983411865519&lsd=AVrQgTMWncg&variables={"username":"` + username + `"}`)
	
	request.Header.Set("IG-U-DS-USER-ID", randomIntString(11))
	request.Header.Set("IG-INTENDED-USER-ID", randomIntString(11))
	request.Header.Set("Cookie", "s_user_id=" + randomIntString(11))

	fasthttp.Do(request, response)
	var body []byte = response.Body()
	
	if (!strings.Contains(string(body), `"user":`)) {
		fmt.Printf("ERROR: Unable to fetch new username ID %s \n\n", strings.Repeat(" ", 45))
		os.Exit(0)
	}
	if usernameID := fastjson.GetString(body, "data", "user", "id"); (usernameID != "") {
		return usernameID
	}
	return "Unavailable"
}

func sendDiscordWebhook(webhook, template string) {
	var request *fasthttp.Request = createRequest("POST", false)
	var response *fasthttp.Response = fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(request)
	defer fasthttp.ReleaseResponse(response)

	request.SetRequestURI(webhook)

	request.Header.SetContentType("application/json")
	request.SetBody([]byte(template))

	fasthttp.Do(request, response)
}