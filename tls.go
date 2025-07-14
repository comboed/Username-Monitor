package main

import (
	"github.com/valyala/fasthttp"
	"crypto/tls"
	"strings"
	"strconv"
	"bufio"
	"time"
	"net"
	"fmt"
)

func createTCPConnection() *net.TCPConn {
	for {
		var address, _ = net.ResolveTCPAddr("tcp4", "31.13.80.52:443")
		var connection, error = net.DialTCP("tcp4", nil, address)
		if (error == nil) {
			connection.SetLinger(0)
			connection.SetNoDelay(true)
			connection.SetKeepAlive(true)
			return connection
		}
		fmt.Printf("ERROR: Failed to dial to 31.13.80.52:443 %s \n\n", strings.Repeat(" ", 45))
		time.Sleep(time.Millisecond * 100)
	}
}

func createTLSConnection() *tls.Conn {
	for {
		var connection *tls.Conn = tls.Client(createTCPConnection(), &tls.Config {
			InsecureSkipVerify: true,
			MinVersion: tls.VersionTLS13,
			MaxVersion: tls.VersionTLS13,
		})
		if (connection.Handshake() == nil) {
			var response *fasthttp.Response = fasthttp.AcquireResponse()
			connection.Write(buildRequest("", "")) // This is a test to see if pre-writing to a connection will connection information to make sending faster
			response.Read(bufio.NewReaderSize(connection, 8096))
			return connection
		}
		fmt.Printf("ERROR: Failed to create tls connection %s \n\n", strings.Repeat(" ", 45))
		time.Sleep(time.Millisecond * 100)
	}
}

func connectionRefresher(connection **tls.Conn) {
    for {
        time.Sleep(time.Second * 15)
        if (synchronizeA) {
            time.Sleep(time.Second * 10)
        }
        var oldConnection = *connection
        *connection = createTLSConnection()
        oldConnection.Close()
    }
}

func buildRequest(sessionid, username string) []byte {
	return []byte("POST https://i.instagram.com/api/v1/accounts/update_profile_username/ HTTP/1.1\r\n" +
		   "Host: i.instagram.com\r\n" +
		   "Connection: keep-alive\r\n" +
		   "Cookie: sessionid=" + sessionid + "\r\n" +
		   "Content-Type: application/x-www-form-urlencoded\r\n" +
		   "User-Agent: Instagram 336.0.0.35.90 Android (30/11; 280dpi; 720x1411; samsung; SM-A115F; a11q; qcom; en_US; 671551917)\r\n" +
		   "Content-Length: " + strconv.Itoa(9 + len(username)) + "\r\n\r\n" +
		   "username=" + username)
}