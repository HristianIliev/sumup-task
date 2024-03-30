package dbmodels

import "errors"

type Receiver struct {
	ID                 string
	ChannelPreferences []*ChannelPreference `dynamodbav:"channelPreferences"`
	Email              string
	Phone              string
	Endpoint           string
	IsDeleted          bool `dynamodbav:"isDeleted"`
}

type ChannelPreference struct {
	Channel    string
	Preference Preference
}

type Preference string

const (
	PreferenceNever  Preference = "never"
	PreferenceAlways Preference = "always"
)

var ErrInvalidPreference = errors.New("invalid preference name")

var dMap = map[string]Preference{
	"never":  PreferenceNever,
	"always": PreferenceAlways,
}

func NewPreferenceEnum(s string) (Preference, error) {
	value, found := dMap[s]
	if !found {
		return "", ErrInvalidPreference
	}

	return value, nil
}
