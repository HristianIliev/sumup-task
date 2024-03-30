package apimodels

type Receiver struct {
	ID                 string               `json:"id,omitempty"`
	Email              string               `json:"email,omitempty"`
	Phone              string               `json:"phone,omitempty"`
	Endpoint           string               `json:"endpoint,omitempty"`
	ChannelPreferences []*ChannelPreference `json:"channelPreferences,omitempty"`
}

type ChannelPreference struct {
	Channel    string `json:"channel,omitempty"`
	Preference string `json:"preference,omitempty"`
}
