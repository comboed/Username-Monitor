package main

const (
	BatchSize int = 21

	InstagramLsd string = "AVrrNwPTtyg"
	InstagramCsrf string = "OF-fzaKjh9mBpFwxPnx6h_"
	instagramDatr string = "eRybZ6Nt2KD8QA8Dye8PMT_d"
	InstagramSessionID string = "67921106280%3Abt828QpakP70s2%3A6%3AAYcQvMr1FDz_oXfUnVsO9Ae1S-G8_bY_9NJ0ysf9gw"
	InstagramFBDtsg string = "NAcOTZ7igcSHDAnqVtVfSJpYuYzjHz7YIOY8NsaFrMdwZAG7wZ-yMPw:17865379441060568:1738206894"

	InstagramWBloksFetchAsyncUrl string = `https://www.instagram.com/async/wbloks/fetch/?__a=1&lsd=%s&appid=com.bloks.www.bloks.caa.login.process_client_data_and_redirect&__bkv=6ffa76cd52f4fb43c58effe983d6ad7dd8ead63a7d2e4f5b41b54fa30b5ce39d&type=app&params={"params":"{\"server_params\":{\"account_list\":[{\"uid\":\"%s`
	InstagramBloksGraphQLUrl string = `https://i.instagram.com/graphql_www?doc_id=7765850536785467&variables={"input":{"appid":"com.bloks.www.bloks.caa.login.process_client_data_and_redirect","bloks_versioning_id":"6ffa76cd52f4fb43c58effe983d6ad7dd8ead63a7d2e4f5b41b54fa30b5ce39d","params":"{\"account_list\":[{\"uid\":\"%s`

	DiscordAutoclaimedWebhook string = "https://discord.com/api/webhooks/1182116405383004270/Sp4xSu2hQSgSsIX0M7vYgwdyuzOy5KBg_SR3x0LUt7BSIbDxL1XmmfGq8X2SC91uC_62"
	DiscordMissedWebhook string = "https://discord.com/api/webhooks/1182546747948535870/Vp5jbcryjB82Z7I5y9foigFKkvtwP_ClKQfgiPdi0AdCnKRLX3E3LYXgmhnuawFqiwdD"
	DiscordVoidWebhook string = "https://discord.com/api/webhooks/1182706975151231137/_6lhzScC4WU-TedpM8_dezEEXP9C6_OU_T_cJoAbAIOyz6kBy7W3a2rzcu1ugkgvLenF"

	DiscordAutoclaimedTemplate string = `{"content":null,"embeds":[{"color":null,"author":{"name":"Autoclaimed @%s","url":"https://www.instagram.com/%s","icon_url":"https://static-00.iconduck.com/assets.00/high-voltage-emoji-1304x2048-a4e802ha.png"},"timestamp":"%s"}],"attachments":[]}`
	DiscordMissedTemplate string = `{"content":null,"embeds":[{"color":5814783,"fields":[{"name":"Previous ID","value":"%s","inline":true},{"name":"New ID","value":"%s","inline":true}],"author":{"name":"Missed @%s","url":"https://www.instagram.com/%s","icon_url":"https://lh6.googleusercontent.com/proxy/SRLBfptaxlBv9uR1hjqGguwLy4bRXyKFZpwOUM2-yyf1P_LI3v1rtZIpinEAarDh2T7TBEgaVtcMsxBabIheDgQ64RbMo48JPnk05ifwsApKhevWbsVWpBQrzrlNuqwTjNSdsW25"},"timestamp":"%s"}],"attachments":[]}`
	DiscordVoidTemplate string = `{"content":null,"embeds":[{"color":5059879,"fields":[{"name":"ID","value":"%s"}],"author":{"name":"@%s Voided","url":"https://www.instagram.com/%s","icon_url":"https://imgproxy.attic.sh/unsafe/rs:fit:768:768:1:1/t:1:FF00FF:false:false/pngo:false:true:256/aHR0cHM6Ly9hdHRp/Yy5zaC82MDcxd3I1/M2NyYzkyb3F1ODF2/Y3NuaTVmbWMz.png"},"timestamp":"%s"}],"attachments":[]}`
	DiscordUnvoidTemplate string = `{"content":null,"embeds":[{"color":16753776,"fields":[{"name":"ID","value":"%s"}],"author":{"name":"@%s Returned","url":"https://www.instagram.com/%s","icon_url":"https://imgproxy.attic.sh/unsafe/rs:fit:768:768:1:1/t:1:FF00FF:false:false/pngo:false:true:256/aHR0cHM6Ly9hdHRp/Yy5zaC82MDcxd3I1/M2NyYzkyb3F1ODF2/Y3NuaTVmbWMz.png"},"timestamp":"%s"}],"attachments":[]}`
)

var (
	usernames []string
	voidUsernames []string
	sessions []string

	usernameGroups [][]string
	requestMap map[string][][]byte
	requests [][]byte

	rotateEndpoint bool
	synchronizeA bool
	synchronizeB bool

	attempts int64
	rl int64
)