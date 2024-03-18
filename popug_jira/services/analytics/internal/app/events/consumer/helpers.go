package consumer

import (
	"errors"
	"regexp"
	"strings"

	pbV1TaskCreated "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/taskevents/created/v1"
	pbV2TaskCreated "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/taskevents/created/v2"
	pbV1TxCredited "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/transactionevents/credited/v1"
	pbV1TxDebited "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/transactionevents/debited/v1"
	pbV1TxPaymented "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/transactionevents/paymented/v1"
	pbV1UserCreated "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/userevents/created/v1"
	pbV1UserUpdated "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/userevents/updated/v1"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func getUnmarshalledPbEventFromBinary(data []byte) protoreflect.ProtoMessage {
	// task event types
	taskCreatedV1 := &pbV1TaskCreated.Event{}
	taskCreatedV2 := &pbV2TaskCreated.Event{}
	// transactions event types
	txCreditedV1 := &pbV1TxCredited.Event{}
	txDebitedV1 := &pbV1TxDebited.Event{}
	txPaymentedV1 := &pbV1TxPaymented.Event{}
	// user event types
	userCreatedV1 := &pbV1UserCreated.Event{}
	userUpdatedV1 := &pbV1UserUpdated.Event{}

	// TODO: очень глупая ошибка, пока как костыль. Но надо исправить
	if err := proto.Unmarshal(data, taskCreatedV1); err == nil && taskCreatedV1.GetHeader().GetEventName() == "TaskCreated" && taskCreatedV1.GetHeader().GetEventVersion() == "1" {
		logrus.Debugf(
			"fetched event %q version %q:. Start processing %v",
			taskCreatedV1.GetHeader().GetEventName(),
			taskCreatedV1.GetHeader().GetEventVersion(),
			taskCreatedV1,
		)

		return taskCreatedV1
		// TODO: очень глупая ошибка, пока как костыль. Но надо исправить
	} else if err := proto.Unmarshal(data, taskCreatedV2); err == nil && taskCreatedV2.GetHeader().GetEventName() == "TaskCreated" && taskCreatedV2.GetHeader().GetEventVersion() == "2" {
		logrus.Debugf(
			"fetched event %q version %q:. Start processing %v",
			taskCreatedV1.GetHeader().GetEventName(),
			taskCreatedV1.GetHeader().GetEventVersion(),
			taskCreatedV1,
		)

		return taskCreatedV2
		// TODO: очень глупая ошибка, пока как костыль. Но надо исправить
	} else if err := proto.Unmarshal(data, userCreatedV1); err == nil && userCreatedV1.GetHeader().GetEventName() == "CreatedUser" {
		logrus.Debugf(
			"fetched event %q:. Start processing %v",
			userCreatedV1.GetHeader().GetEventName(),
			userCreatedV1,
		)

		return userCreatedV1
		// TODO: очень глупая ошибка, пока как костыль. Но надо исправить
	} else if err := proto.Unmarshal(data, userUpdatedV1); err == nil && userUpdatedV1.GetHeader().GetEventName() == "UpdatedUser" {
		logrus.Debugf(
			"fetched event %q:. Start processing %v",
			userUpdatedV1.GetHeader().GetEventName(),
			userUpdatedV1,
		)

		return userUpdatedV1
		// TODO: очень глупая ошибка, пока как костыль. Но надо исправить
	} else if err := proto.Unmarshal(data, txCreditedV1); err == nil && txCreditedV1.GetHeader().GetEventName() == "TransactionCredited" {
		logrus.Debugf(
			"fetched event %q:. Start processing %v",
			txCreditedV1.GetHeader().GetEventName(),
			txCreditedV1,
		)

		return txCreditedV1
		// TODO: очень глупая ошибка, пока как костыль. Но надо исправить
	} else if err := proto.Unmarshal(data, txDebitedV1); err == nil && txDebitedV1.GetHeader().GetEventName() == "TransactionDebited" {
		logrus.Debugf(
			"fetched event %q:. Start processing %v",
			txDebitedV1.GetHeader().GetEventName(),
			txDebitedV1,
		)

		return txDebitedV1
	} else if err := proto.Unmarshal(data, txPaymentedV1); err == nil && txPaymentedV1.GetHeader().GetEventName() == "TransactionPaymented" {
		logrus.Debugf(
			"fetched event %q:. Start processing %v",
			txPaymentedV1.GetHeader().GetEventName(),
			txPaymentedV1,
		)

		return txPaymentedV1
	} else {
		return nil
	}
}

func validateJiraIDAndTitle(jiraID, title string) error {
	if strings.ReplaceAll(jiraID, " ", "") == "" {
		return errors.New("jira id must be set")
	}

	if containsBracketSequence(title) {
		return errors.New("jira id is not be in title")
	}

	return nil
}

func containsBracketSequence(s string) bool {
	re := regexp.MustCompile(`.*\[.*\].*`)
	return re.MatchString(s)
}
