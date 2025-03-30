package mail

import (
	"fmt"
	"strconv"
	"time"

	"scg-mail-discord-bot/internal/config"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

// check emails
func CheckEmails() ([]*imap.Message, error) {
	client, err := connectToIMAP()
	if err != nil {
		return nil, fmt.Errorf("IMAP 연결 실패: %w", err)
	}
	defer client.Logout()

	messages, err := fetchRecentMessages(client)
	if err != nil {
		return nil, fmt.Errorf("메시지 가져오기 실패: %w", err)
	}

	filter := NewDomainFilter(config.App.Mail.AllowedDomains, true)
	return filter.FilterMessages(messages), nil
}

// connect to IMAP
func connectToIMAP() (*client.Client, error) {
	cfg := config.App.Mail
	imapClient, err := client.DialTLS(cfg.Host+":"+strconv.Itoa(cfg.Port), nil)
	if err != nil {
		return nil, err
	}

	if err := imapClient.Login(cfg.Username, cfg.Password); err != nil {
		return nil, err
	}

	return imapClient, nil
}

// fetch recent messages
func fetchRecentMessages(client *client.Client) ([]*imap.Message, error) {
	// Select the INBOX
	if _, err := client.Select("INBOX", false); err != nil {
		return nil, err
	}

	// search emails since 1 hour ago
	criteria := imap.NewSearchCriteria()
	criteria.Since = time.Now().Add(-1 * time.Hour)

	uids, err := client.Search(criteria)
	if err != nil {
		return nil, err
	}

	if len(uids) == 0 {
		return nil, nil
	}

	// fetch emails
	seqset := new(imap.SeqSet)
	seqset.AddNum(uids...)

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- client.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope, imap.FetchBody, imap.FetchFlags}, messages)
	}()

	// result
	var result []*imap.Message
	for msg := range messages {
		result = append(result, msg)
	}

	if err := <-done; err != nil {
		return nil, err
	}

	return result, nil
}
