package mail

type MailMessage struct {
	To          []string
	Subject     string
	ContentHTML string
}
