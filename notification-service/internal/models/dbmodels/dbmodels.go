package dbmodels

import (
	"notification-service/internal/models/enums"
)

type Receiver struct {
	ID                 string
	ChannelPreferences []*ChannelPreference `dynamodbav:"channelPreferences"`
	Email              string
	Phone              string
	Endpoint           string
	IsDeleted          bool `dynamodbav:"isDeleted"`
}

type ChannelPreference struct {
	Channel    enums.Channel
	Preference enums.Preference
}
