package consumer

import (
	"errors"
	"math/rand"
	"regexp"
	"strings"

	pbV1TaskAssigned "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/taskevents/assigned/v1"
	pbV1TaskCompleted "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/taskevents/completed/v1"
	pbV1TaskCreated "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/taskevents/created/v1"
	pbV2TaskCreated "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/taskevents/created/v2"
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
	taskAssignedV1 := &pbV1TaskAssigned.Event{}
	taskCompletedV1 := &pbV1TaskCompleted.Event{}
	// user event types
	userCreatedV1 := &pbV1UserCreated.Event{}
	userUpdatedV1 := &pbV1UserUpdated.Event{}

	// TODO: очень глупая ошибка, пока как костыль. Но надо исправить
	if err := proto.Unmarshal(data, taskCreatedV1); err == nil && taskCreatedV1.GetHeader().GetEventName() == "TaskCreated" && taskCreatedV1.GetHeader().GetEventVersion() == "1" {
		logrus.Debugf(
			"fetched event %q:. Start processing %v",
			taskCreatedV1.GetHeader().GetEventName(),
			taskCreatedV1,
		)

		return taskCreatedV1
		// TODO: очень глупая ошибка, пока как костыль. Но надо исправить
	} else if err := proto.Unmarshal(data, taskCreatedV2); err == nil && taskCreatedV2.GetHeader().GetEventName() == "TaskCreated" && taskCreatedV2.GetHeader().GetEventVersion() == "2" {
		logrus.Debugf(
			"fetched event %q:. Start processing %v",
			taskCreatedV2.GetHeader().GetEventName(),
			taskCreatedV2,
		)

		return taskCreatedV2
		// TODO: очень глупая ошибка, пока как костыль. Но надо исправить
	} else if err := proto.Unmarshal(data, taskAssignedV1); err == nil && taskAssignedV1.GetHeader().GetEventName() == "TaskAssigned" {
		logrus.Debugf(
			"fetched event %q:. Start processing %v",
			taskAssignedV1.GetHeader().GetEventName(),
			taskAssignedV1,
		)

		return taskAssignedV1
		// TODO: очень глупая ошибка, пока как костыль. Но надо исправить
	} else if err := proto.Unmarshal(data, taskCompletedV1); err == nil && taskCompletedV1.GetHeader().GetEventName() == "TaskCompleted" {
		logrus.Debugf(
			"fetched event %q:. Start processing %v",
			taskCompletedV1.GetHeader().GetEventName(),
			taskCompletedV1,
		)

		return taskCompletedV1
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
	} else {
		return nil
	}
}

func generateTaskCosts() (int, int) {
	assignCost := rand.Intn(11) + 10
	completeCost := rand.Intn(21) + 20

	return assignCost, completeCost
}

func parseTaskDescription(desc string) (string, string) {
	re := regexp.MustCompile(`\[(.*?)\](.*)`)
	matches := re.FindStringSubmatch(desc)

	if len(matches) == 3 {
		return matches[1], strings.TrimSpace(matches[2])
	}

	return "", strings.TrimSpace(desc)
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
