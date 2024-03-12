package events

import pbV1UserRoleV1 "github.com/lyckety/async_arch/popug_jira/schema-registry/pkg/pb/userevents/userrole/v1"

const producerName = "auth"

func convertStringRoleToRPCRoleV1(role string) (pbV1UserRoleV1.UserRole, error) {
	switch role {
	case "administrator":
		return pbV1UserRoleV1.UserRole_USER_ROLE_ADMIN, nil
	case "manager":
		return pbV1UserRoleV1.UserRole_USER_ROLE_MANAGER, nil
	case "bookkeeper":
		return pbV1UserRoleV1.UserRole_USER_ROLE_BOOKKEEPER, nil
	case "worker":
		return pbV1UserRoleV1.UserRole_USER_ROLE_WORKER, nil
	default:
		return 0, ErrUnknownUserRole
	}
}
