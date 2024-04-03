package service

import (
	"context"
	"errors"
	"log"
	"notification-service/internal/models/dbmodels"
	"notification-service/internal/models/enums"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/google/uuid"
)

var defaultPreferences = []*dbmodels.ChannelPreference{
	{Channel: "email", Preference: enums.PreferenceAlways},
	{Channel: "sms", Preference: enums.PreferenceAlways},
	{Channel: "slack", Preference: enums.PreferenceAlways},
}

type ReceiverService struct {
	db        *dynamodb.Client
	tableName string
}

func NewReceiverService(db *dynamodb.Client, tableName string) *ReceiverService {
	return &ReceiverService{
		db:        db,
		tableName: tableName,
	}
}

var ErrReceiverNotFound = errors.New("receiver not found")

func (r *ReceiverService) GetReceiver(id string) (*dbmodels.Receiver, error) {
	output, err := r.db.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
		TableName: aws.String(r.tableName),
	})
	if err != nil {
		return nil, err
	}

	if output.Item == nil {
		return nil, ErrReceiverNotFound
	}

	receiver, err := unmarshalMap(output.Item)
	if err != nil {
		return nil, err
	}

	if receiver.ChannelPreferences == nil || len(receiver.ChannelPreferences) == 0 {
		receiver.ChannelPreferences = defaultPreferences
	}
	return receiver, nil
}

func (r *ReceiverService) CreateReceiver(body *dbmodels.Receiver) (*dbmodels.Receiver, error) {
	UUID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	newReceiver := &dbmodels.Receiver{
		ID:                 UUID.String(),
		Email:              body.Email,
		Endpoint:           body.Endpoint,
		Phone:              body.Phone,
		ChannelPreferences: body.ChannelPreferences,
		IsDeleted:          false,
	}

	item, err := marshalMap(newReceiver)
	if err != nil {
		log.Println("Error marshalling map:", err)
		return nil, err
	}

	conditionExpression := "attribute_not_exists(#ID)"
	expressionAttributeNames := map[string]string{
		"#ID": "ID",
	}
	input := &dynamodb.PutItemInput{
		Item:                     item,
		TableName:                aws.String(r.tableName),
		ConditionExpression:      &conditionExpression,
		ExpressionAttributeNames: expressionAttributeNames,
	}

	_, err = r.db.PutItem(context.TODO(), input)
	if err != nil {
		var ccf *types.ConditionalCheckFailedException
		if errors.As(err, &ccf) {
			log.Println("Item with the same ID already exists")
			return nil, err
		}

		log.Println("Error putting item:", err)
		return nil, err
	}

	return newReceiver, nil
}

func marshalMap(input interface{}) (map[string]types.AttributeValue, error) {
	av, err := attributevalue.MarshalMap(input)
	if err != nil {
		return nil, err
	}

	return av, nil
}

func unmarshalMap(item map[string]types.AttributeValue) (*dbmodels.Receiver, error) {
	var result dbmodels.Receiver
	err := attributevalue.UnmarshalMap(item, &result)
	if err != nil {
		log.Println("Error unmarshalling item:", err)
		return nil, err
	}

	return &result, nil
}
