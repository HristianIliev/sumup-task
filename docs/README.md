# Setup

## Notification-Service

Create DynamoDB table with needed Partition Key and Sort Key
Create 


API_PORT=8081 TABLE_NAME="notification-receivers" SNS_TOPIC="arn:aws:sns:eu-central-1:176099643795:notifications.fifo" go run ./cmd/notification-service/
QUEUE_URL="https://sqs.eu-central-1.amazonaws.com/176099643795/sms" WORKERS=10 POLL_DELAY=10 go run ./cmd/notification-processor/