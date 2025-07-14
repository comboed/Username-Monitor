package main

import (
	"github.com/valyala/fasthttp"
	"os/signal"
	_"net/http"
	"strings"
	"syscall"
	"time"
	"fmt"
	"os"
)

func main() {
	fmt.Printf("Instagram Autoclaimer \n")

	usernames = openFile("./data/usernames.txt")
	voidUsernames = openFile("./data/void_usernames.txt")
	fmt.Printf("Usernames: %s \n", formatNumber(int64(len(usernames)) + int64(len(voidUsernames))))

	sessions = openFile("./data/sessions.txt")
	fmt.Printf("Sessions: %s \n", formatNumber(int64(len(sessions))))

	var proxies []string = openFile("./data/proxies.txt")
	fmt.Printf("Proxies: %s \n", formatNumber(int64(len(proxies))))

	var goroutines int
	fmt.Printf("Goroutines: ")
	fmt.Scanln(&goroutines)
	fmt.Println()

	var kill chan os.Signal = make(chan os.Signal, 1)
	signal.Notify(kill, os.Interrupt, syscall.SIGTERM)

	// Begin building data
	buildUsernameGroups()
	buildRequests()

	var usernameChannel chan []string = make(chan []string)
	var voidUsernameChannel chan []string = make(chan []string)
	var requestChannel chan []byte = make(chan []byte)
	var proxyChannel chan string = make(chan string)
	
	var client *fasthttp.Client = createClient(proxyChannel)
	for i := 0; i < goroutines; i++ {
		go autoClaimer(client, usernameChannel)
	}

	for i := 0; i < 10; i++ {
		go voidMonitor(client, voidUsernameChannel)
	}

	for i := 0; i < 30; i++ { // 30 test first
		go usernameSpammer(requestChannel)
	}

	go func() {
		for {
			for _, group := range usernameGroups {
				usernameChannel <- group
			}
		}
	}()

	go func() {
		for {
			if len(requests) > 0 { // I'm unsure if this fixes the crash at start
				for _, request := range requests {
					requestChannel <- request
				}
			}
		}	
	}()

	go func() {
		for {
			for _, proxy := range proxies {
				proxyChannel <- proxy
			}
		}
	}()

	go func() {
		for len(voidUsernames) > 0 {
			for _, username := range voidUsernames {
				var parts []string = strings.Split(username, ":")
				var url string = `\",\"credential_type\":\"none\",\"token\":\"\"}]}}"}`	
				if (rotateEndpoint) {
					url = fmt.Sprintf(InstagramBloksGraphQLUrl, parts[1]) + `\",\"credential_type\":\"none\",\"token\":\"\"}]}"}}`
				}
				voidUsernameChannel <- []string{url, parts[0], parts[1]}
			}
		}
	}()

	go func() {
		for {
			if len(voidUsernames) > 0 {
				for _, username := range voidUsernames {
					var parts []string = strings.Split(username, ":")
					if (rotateEndpoint) {
						voidUsernameChannel <- []string{fmt.Sprintf(InstagramBloksGraphQLUrl, parts[1]) + `\",\"credential_type\":\"none\",\"token\":\"\"}]}"}}`, parts[0], parts[1]}
					} else {
						voidUsernameChannel <- []string{fmt.Sprintf(InstagramWBloksFetchAsyncUrl, InstagramLsd, parts[1]) + `\",\"credential_type\":\"none\",\"token\":\"\"}]}}"}`, parts[0], parts[1]}
					}
				}	
			}
		}
	}()

	
	go func () {
		for {
			time.Sleep(time.Hour)
			verifySessions()
		}
	}()

	var epr int64
	go func() {
		for {
			time.Sleep(time.Hour)
			rotateEndpoint = !rotateEndpoint
			epr += 1
			if (synchronizeA) {
				time.Sleep(time.Second * 15)
			}
			buildUsernameGroups()
		}
	}()

	var rps int64
	go func() {
		for {
			var before int64 = attempts
			time.Sleep(time.Second * 1)
			rps = attempts - before
		}
	}()

	go func() {
		for {
			for _, spinner := range []string{"|", "/", "-", "\\", "|", "/", "-", "\\"} {
				fmt.Printf("[%s] Autoclaiming - Attempts: %s | RPS: %s | CPS: %s | EPR: %s | RLs: %s %s \r", spinner, formatNumber(attempts), formatNumber(rps),  formatNumber(rps * int64(BatchSize)), formatNumber(epr), formatNumber(rl), strings.Repeat(" ", 5))
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()

	go func() {
		<- kill
		fmt.Printf("A/C stopped, exiting after %s attempts... %s \r\n", formatNumber(attempts), strings.Repeat(" ", 60))
		os.Exit(0)
	}()

	select {}
}

