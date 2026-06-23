package enums

type contextKey string

const (
	ContextKeyAuthId           contextKey = "auth_id"
	ContextKeyAuthUser         contextKey = "auth_user"
	ContextKeyImpersonatorId   contextKey = "impersonator_id"
	ContextKeyRole             contextKey = "role"
	ContextKeyPermissions      contextKey = "permissions"
	ContextKeySessionId        contextKey = "session_id"
	ContextKeyRequestId        contextKey = "request_id"
	ContextKeyRequestStartTime contextKey = "request_start_time"
	ContextKeyRequestEndTime   contextKey = "request_end_time"
	ContextKeyRequestDuration  contextKey = "request_duration"
)

func (c contextKey) String() string {
	return string(c)
}
