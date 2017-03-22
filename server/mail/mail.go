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
	Subject: "Reset your Basiq password",
	Sender: sender{
		Name:  "Basiq",
		Email: "support@basiq.io",
	},
	Message: `
		<p>Hi %s!<br/><br/>
		You have requested a password reset for your Basiq account. Follow the link below to set a new password:<br/><br/>
		<a href="https://dashboard.basiq.io/reset/%s">https://dashboard.basiq.io/reset/%s</a><br/><br/>
		If you don't wish to reset your password, disregard this email and no action will be taken.<br/><br/>
		Yours,<br/><br/>
		The Basiq Team`,
}

var Template = struct {
	Subject string
	Sender  sender
	Message string
}{
	Subject: "Confirm yout Basiq email address",
	Sender: sender{
		Name:  "Basiq",
		Email: "support@basiq.io",
	},
	Message: `
	<p>Welcome to Basiq!<br/><br/>
	Before you can start integrating with our APIs, you need to confirm your email address. To get started, just confirm your email address by clicking the link below:<br/><br/>
	<a href="https://dashboard.basiq.io/signup/accept/%s">https://dashboard.basiq.io/signup/accept/%s</a><br/><br/>
	Once you're ready to start integrating, we recommended taking a look at our docs:<br/><br/>
	<a href="http://basiq.io/api/">http://basiq.io/api/</a><br/><br/>
	You can view your setup, API request logs, and a variety of other information about your account right from dashboard:<br/><br/>
	<a href="https://dashboard.basiq.io">https://dashboard.basiq.io</a><br/><br/>
	We'll be here to help you with every step along the way. You can find answers to most questions and get in touch with us at <a href="https://support.basiq.io">https://support.basiq.io</a><br/><br/>
	Hope you enjoy getting up and running. We're excited to see what comes next!<br/><br/>
	Yours,<br/><br/>
	The Basiq Team`,
}
