package main

import (
	"AlertMechanism/service/alert"
	"encoding/json"
	"fmt"
)

const (
	KafkaServer  = "localhost:9092"
	KafkaTopic   = "orders-v1-topic"
	KafkaGroupId = "user-service"
)

func main() {
	kc := alert.NewKafaConsumer("brkedId", "cleintId", "cosumerGrp", "cosumertopic")
	kreader := kc.GetKafkareader()

	// will keep reading incming messages and will create detections using factory
	for {
		msg, err := kreader.ReadMessage(-1)
		if err == nil {
			alfactory := NewAlertCreatorFactory()
			alert, err := alfactory.createAlert(string(msg.Value))
			if err != nil {
				fmt.Printf("Error decoding message: %v\n", err)
				continue
			}
			fmt.Printf("Received Order: %+v\n", alert)
		} else {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

type Alert interface {
	processNotification()
}

type AlertCreatorFactory struct {
	alertmsgkafka string
}

func NewAlertCreatorFactory() *AlertCreatorFactory {
	return &AlertCreatorFactory{}
}

func (altfactory *AlertCreatorFactory) createAlert(msgString string) (Alert, error) {
	var detection Detection
	err := json.Unmarshal([]byte(msgString), detection)
	if err != nil {
		return nil, err
	}
	if detection.event_type == "secret_detection" {
		return &VulnerabilityDetection{}, nil
	} else if detection.event_type == "vulnerability_detection" {
		return &SecretDetection{}, nil
	}
	return nil, fmt.Errorf("no format found")
}

type Prioroity int

const (
	HIGH Prioroity = iota + 1
	MEDDIUM
	LOW
)

type Detection struct {
	event_type string
	severity   Prioroity
	org_id     string
	metadata   map[string]interface{}
}

type VulnerabilityDetection struct {
	cloud_account_id string
	Detection
	minPriortyForMsgsebd Prioroity
}

func (vd *VulnerabilityDetection) processNotification() {

}

type SecretDetection struct {
	version_control_integration_id string
	Detection
	minPriortyForMsgsebd Prioroity
}

func (sd *SecretDetection) processNotification() {
	// form notification message
	notifitifireService := NotifitifireService{}
	var notifiers []NotificationSendpr
	if sd.metadata.Priority < sd.minPriortyForMsgsebd {
		notifiers = notifitifireService.createNotificationSenders(sd.metadata)
	}
	// send notification
	for _, notificationsebder := range notifiers {
		notificationsebder.sendNotification()
	}

	// collect notification sender list

}

type Notification interface {
	NotificationCreator
	NotificationSendpr
}

type MessageCreator interface {
	createMessage(metadata map[string]interface{}) string
}

type MessageCreatorfactpry struct {
	msgCreator MessageCreator
}

func (msgfactory MessageCreatorfactpry) createMessage(metadata map[string]interface{}) string {
	return msgfactory.createMessage(metadata)
}

type EmailMessageCreator struct {
}

func (sm EmailMessageCreator) createMessage(metadata map[string]interface{}) string {
	slackmsg := ""
	for _, mes := range metadata {
		// form slac
		slackmsg = mes.(string)
	}
	return slackmsg
}

type JIRAMessageCreator struct {
}

func (sm JIRAMessageCreator) createMessage(metadata map[string]interface{}) string {
	slackmsg := ""
	for _, mes := range metadata {
		// form slac
		slackmsg = mes.(string)
	}
	return slackmsg
}

type SlackMessageCreator struct {
}

func (sm SlackMessageCreator) createMessage(metadata map[string]interface{}) string {
	slackmsg := ""
	for _, mes := range metadata {
		// form slac
		slackmsg = mes.(string)
	}
	return slackmsg
}

type NotifitifireService struct {
}

func (ns NotifitifireService) createNotificationSenders(metadata map[string]interface{}) []NotificationSendpr {
	arr := make([]NotificationSendpr, 0)
	msgfactory := MessageCreatorfactpry{}

	for _, mes := range metadata {
		if mes == "slackenabled" {
			smmessage := SlackMessageCreator{}
			msgfactory.msgCreator = smmessage
			arr = append(arr, SlackNotification{
				sclakchannel: mes.channmel,
				slacktoken:   mes.token,
				message:      msgfactory.createMessage(mes),
			})
		} else if mes == "emailenabled" {
			smmessage := EmailMessageCreator{}
			msgfactory.msgCreator = smmessage
			arr = append(arr, EmailNotification{
				emailID: mes.emailid,
				message: msgfactory.createMessage(mes),
			})
		} else if mes == "jiraenabled" {
			smmessage := JIRAMessageCreator{}
			msgfactory.msgCreator = smmessage
			arr = append(arr, JIRANotification{
				jiraurl: mes.jiraurl,
				token:   mes.token,
				message: msgfactory.createMessage(mes),
			})
		}
	}
	return arr
}

type NotificationSendpr interface {
	sendNotification()
}

type SlackNotification struct {
	sclakchannel string
	slacktoken   string
	message      string
}

func (slnf SlackNotification) sendNotification() {
	slacklib.sendnotification(slnf.sclakchannel, slnf.slacktoken, slnf.message)
}

type EmailNotification struct {
	emailID string
	message string
}

func (emlf EmailNotification) sendNotification() {
	emaillib.sendnotification(emlf.emailID, emlf.message)
}

type JIRANotification struct {
	jiraurl string
	message string
	token   string
}

func (jiranf JIRANotification) sendNotification() {
	jiralib.sendnotification(jiranf.jiraurl, jiranf.token, jiranf.message)
}
