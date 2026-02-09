package handler

import "encoding/json"

type webhookPayload struct {
	Action       string               `json:"action"`
	Installation *webhookInstallation `json:"installation"`
	Sender       *webhookSender       `json:"sender"`
}

type webhookInstallation struct {
	ID          int64           `json:"id"`
	Account     *webhookAccount `json:"account"`
	SuspendedAt *string         `json:"suspended_at"`
}

type webhookAccount struct {
	AvatarURL *string `json:"avatar_url"`
	ID        int64   `json:"id"`
	Login     string  `json:"login"`
	Type      string  `json:"type"`
}

type webhookSender struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
}

func parseWebhookPayload(body []byte) (*webhookPayload, error) {
	var payload webhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}
	return &payload, nil
}
