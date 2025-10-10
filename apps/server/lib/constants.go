package lib

const (
	AccessTokenCookieName  = "access_token"
	RefreshTokenCookieName = "refresh_token"
)

const (
	RoleAdmin   = "admin"
	RoleTeacher = "teacher"
	RoleStudent = "student"
)

const (
	TableUsers           = "public.users"
	TableFiles           = "files"
	TableFolders         = "folders"
	TableSubjects        = "subjects"
	TableUserOAuthTokens = "user_oauth_tokens"
	TableUserSubjects    = "user_subjects"
	TableAuditLogs       = "audit_logs"
	TableHealthLogs      = "health_logs"
)
