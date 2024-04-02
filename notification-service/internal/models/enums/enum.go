package enums

type Channel string

const (
	Email Channel = "email"
	SMS   Channel = "sms"
	Slack Channel = "slack"
)

var Channels = []Channel{Email, SMS, Slack}

type Preference string

const (
	PreferenceNever  Preference = "never"
	PreferenceAlways Preference = "always"
)
