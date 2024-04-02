package apimodels

import "notification-service/internal/models/enums"

type Receiver struct {
	ID                 string               `json:"id,omitempty"`
	Email              string               `json:"email,omitempty"`
	Phone              string               `json:"phone,omitempty"`
	Endpoint           string               `json:"endpoint,omitempty"`
	ChannelPreferences []*ChannelPreference `json:"channelPreferences,omitempty"`
}

type ChannelPreference struct {
	Channel    enums.Channel    `json:"channel,omitempty"`
	Preference enums.Preference `json:"preference,omitempty"`
}

type Notification struct {
	ReceiverID string          `json:"receiverID,omitempty"`
	Message    string          `json:"message,omitempty"`
	Channels   []enums.Channel `json:"channels,omitempty"`
}

type Message struct {
	Content  string          `json:"content,omitempty"`
	Email    string          `json:"email,omitempty"`
	Phone    string          `json:"phone,omitempty"`
	Endpoint string          `json:"endpoint,omitempty"`
	Channels []enums.Channel `json:"channels,omitempty"`
}

func (m *Message) FilterOutChannel(channel enums.Channel) {
	indexToRemove := -1
	for i, c := range m.Channels {
		if c == channel {
			indexToRemove = i
			break
		}
	}

	channels := make([]enums.Channel, 0, len(m.Channels)-1)
	channels = append(channels, m.Channels[:indexToRemove]...)
	channels = append(channels, m.Channels[indexToRemove+1:]...)

	m.Channels = channels
}

func (m *Message) Clone() *Message {
	return &Message{
		Content:  m.Content,
		Email:    m.Email,
		Phone:    m.Phone,
		Endpoint: m.Endpoint,
		Channels: m.Channels,
	}
}
