package mail

type Receiver struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Type  string `json:"type"`
}

type Receivers []Receiver

type Message struct {
	Html      string    `json:"html"`
	Text      string    `json:"text"`
	Subject   string    `json:"subject"`
	FromEmail string    `json:"from_email"`
	FromName  string    `json:"from_name"`
	To        Receivers `json:"to"`
}

type Mail struct {
	Key     string  `json:"key"`
	Message Message `json:"message"`
	Async   bool    `json:"async"`
}

type sender struct {
	Name  string
	Email string
}

var PasswordResetTemplate = struct {
	Subject string
	Sender  sender
	Message string
}{
	Subject: "Reset your password",
	Sender: sender{
		Name:  "Symphony",
		Email: "support@symphony.is",
	},
	Message: `
		<p>Hello!<br/><br/>
		You have requested a password reset for your account. Follow the link below to set a new password:<br/><br/>
		<a href="https://symphonyGopher.symphony.is/reset/%s">https://symphonyGopher.symphony.is/reset/%s</a><br/><br/>
		If you don't wish to reset your password, disregard this email and no action will be taken.<br/><br/>
		Yours,<br/><br/>
		The Symphony Team`,
}

var Template = struct {
	Subject string
	Sender  sender
	Message string
}{
	Subject: "Confirm yout Symphony email address",
	Sender: sender{
		Name:  "Symphony",
		Email: "support@Symphony.is",
	},
	Message: `
	<p>Welcome to SymphonyGopher!<br/><br/>
	Before you can start, you need to confirm your email address. To get started, just confirm your email address by clicking the link below:<br/><br/>
	<a href="https://symphonyGopher.symphony.is/signup/accept/%s">https://symphonyGopher.symphony.is/signup/accept/%s</a><br/><br/>
	Yours,<br/><br/>
	The Symphony Team`,
}
