package adapters

import (
	"fmt"
	"notification-service/internal/models/dbmodels"
	"notification-service/pkg/models/apimodels"
)

func DbReceiverToApiReceiver(dbReceiver *dbmodels.Receiver) *apimodels.Receiver {
	fmt.Printf("%+v", dbReceiver)
	result := &apimodels.Receiver{
		ID:                 dbReceiver.ID,
		ChannelPreferences: DbPreferencesToApiPreferences(dbReceiver.ChannelPreferences),
		Email:              dbReceiver.Email,
		Phone:              dbReceiver.Phone,
		Endpoint:           dbReceiver.Endpoint,
	}

	return result
}

func DbPreferencesToApiPreferences(dbPreferences []*dbmodels.ChannelPreference) []*apimodels.ChannelPreference {
	result := []*apimodels.ChannelPreference{}
	for _, preference := range dbPreferences {
		result = append(result, &apimodels.ChannelPreference{
			Channel:    preference.Channel,
			Preference: preference.Preference,
		})
	}

	return result
}

func ApiReceiverToDbReceiver(apiReceiver *apimodels.Receiver) *dbmodels.Receiver {
	result := &dbmodels.Receiver{
		ID:                 apiReceiver.ID,
		ChannelPreferences: ApiPreferencesToDbPreferences(apiReceiver.ChannelPreferences),
		Email:              apiReceiver.Email,
		Phone:              apiReceiver.Phone,
		Endpoint:           apiReceiver.Endpoint,
	}

	return result
}

func ApiPreferencesToDbPreferences(apiPreferences []*apimodels.ChannelPreference) []*dbmodels.ChannelPreference {
	result := []*dbmodels.ChannelPreference{}
	for _, preference := range apiPreferences {
		result = append(result, &dbmodels.ChannelPreference{
			Channel:    preference.Channel,
			Preference: preference.Preference,
		})
	}

	return result
}
