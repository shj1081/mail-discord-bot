package mail

import (
	"strings"

	"github.com/emersion/go-imap"
)

// DomainFilter represents a filter for email domains
type DomainFilter struct {
	allowedDomains []string
	unreadOnly     bool // filter unread emails only
}

// NewDomainFilter creates a new DomainFilter with the specified allowed domains
func NewDomainFilter(domains []string, unreadOnly bool) *DomainFilter {
	return &DomainFilter{
		allowedDomains: domains,
		unreadOnly:     unreadOnly,
	}
}

// FilterMessages filters messages based on the sender's domain and read status
func (f *DomainFilter) FilterMessages(messages []*imap.Message) []*imap.Message {
	var filtered []*imap.Message
	for _, msg := range messages {
		if len(msg.Envelope.From) == 0 {
			continue
		}

		// check if the email domain is allowed
		senderAddress := msg.Envelope.From[0].Address()
		if !f.isAllowedDomain(senderAddress) {
			continue
		}

		// check if the email is read
		if f.unreadOnly {
			isRead := false
			for _, flag := range msg.Flags {
				if flag == imap.SeenFlag {
					isRead = true
					break
				}
			}
			if isRead {
				continue
			}
		}

		filtered = append(filtered, msg)
	}
	return filtered
}

// isAllowedDomain checks if the email address belongs to an allowed domain
func (f *DomainFilter) isAllowedDomain(email string) bool {
	for _, domain := range f.allowedDomains {
		if strings.HasSuffix(strings.ToLower(email), "@"+strings.ToLower(domain)) {
			return true
		}
	}
	return false
}
