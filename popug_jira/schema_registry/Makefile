.PHONY: all target1 target2 target3

all: buf-lint buf-generate-auth buf-generate-task-tracker

buf-lint:
	buf lint

buf-generate-auth:
	buf generate --template=./buf.gen.auth.yaml . \
		--exclude-path=./task-tracker-api/tasktracker \
		--exclude-path=./task-tracker-api/task \
		--exclude-path=./cud-events/taskevents \
		-o ../services/auth/pkg/grpc

buf-generate-task-tracker:
	buf generate --template=./buf.gen.task.yaml . \
		--exclude-path=./auth-api/auth \
		--exclude-path=./users-api/users \
		-o ../services/task-tracker/pkg/grpc
