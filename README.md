# Username Monitor


**THIS TOOL NO LONGER WORKS, ALL ENDPOINTS ARE PATCHED**

This is a high-performance username monitor written in Go. It monitors target usernames and attempts to claim them the moment they become available using low-latency TLS 1.3 connections and optimized HTTP request handling.

## Features

- Username availability detection and auto-claiming  
- TLS 1.3 socket-based request engine  
- Proxy support using `fasthttpproxy`  
- Multi-session spoofing with randomized headers and cookies  
- Discord webhook integration for notifications  
- Persistent logging and session management  

## Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/yourusername/Username-Monitor
   cd Username-Monitor
   ```

2. **Build the binary**:
   ```bash
   go build -o autoclaimer
   ```

3. **Prepare the required files** inside `./data/`:

   - `sessions.txt`: One Instagram session per line, formatted like:
     ```
     username:password:email:emailpassword:sessionid
     ```

   - `usernames.txt`: List of target usernames to monitor (one per line).

   - `void_usernames.txt`: This will be updated automatically when voids are detected.

4. **Configure variables in the source code**:

   You must define the following constants in the source before building:
   - `instagramDatr`  
   - `InstagramCsrf`  
   - `DiscordAutoclaimedWebhook`  
   - `DiscordVoidWebhook`  
   - `DiscordMissedWebhook`

## Usage

Run the binary:
```bash
./autoclaimer
```

The tool will:

- Monitor usernames continuously  
- Attempt to claim them if available  
- Detect voided/unvoided states  
- Send claim/void/miss events to Discord  
- Write logs to `./data/logs/`

## How It Works

- **Session Rotation**: Uses pre-authenticated Instagram sessions from `sessions.txt` and rotates them to prevent detection or lockout.
- **Header Spoofing**: Every request includes randomized values for `IG-U-DS-USER-ID`, `IG-INTENDED-USER-ID`, and cookie fields (`csrftoken`, `datr`, etc.).
- **Request Engine**: Requests are issued using `fasthttp` for performance, with TLS 1.3 and compression enabled.
- **Proxying**: Proxies are consumed from a channel and injected into the client’s custom dialer using `fasthttpproxy`.
- **TLS Optimization**: A persistent TLS connection is opened and refreshed every 15 seconds to reduce handshake latency. A request is pre-written to accelerate further communication.
- **Claim Handling**: When a username is detected as claimable, a structured POST request is sent directly to `i.instagram.com/api/v1/accounts/update_profile_username/`.
- **Monitoring & Logging**:
  - Claims are logged to `./data/logs/<username>.log`
  - Voided/unvoided states are logged to their respective files
  - Webhook messages are dispatched using formatted templates

## Legal

This tool interacts with private Instagram endpoints and may violate their Terms of Service. It is provided for educational purposes only. Use at your own risk.
