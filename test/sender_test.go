package test

import (
	"testing"

	"github.com/eugeniopolito/gobetemplate/mail"
	"github.com/eugeniopolito/gobetemplate/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	// Skip send email when run test (make test)
	if testing.Short() {
		t.Skip()
	}

	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := mail.NewEmailConfigSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword, config.SmtpAuthAddress, config.SmtpServerAddress)

	subject := "A test email"
	content := `
	<h1>Hello world</h1>
	<p>This is a test message from GO BE Template</p>
	`
	to := []string{"eugenio.polito@gmail.com"}
	attachFiles := []string{"../go.mod"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
