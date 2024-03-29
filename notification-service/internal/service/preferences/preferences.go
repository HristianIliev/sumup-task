package preferences

import "github.com/aws/aws-sdk-go-v2/service/dynamodb"

type PreferencesService struct {
}

func New(db *dynamodb.Client) *PreferencesService {
	return &PreferencesService{}
}
