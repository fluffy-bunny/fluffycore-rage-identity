package error

type (
	FlowError string
)

const (
	FlowError_UserNotFound        FlowError = "FlowError_UserNotFound"
	FlowError_BadSession          FlowError = "FlowError_BadSession"
	FlowError_InvalidRequest      FlowError = "FlowError_InvalidRequest"
	FlowError_EmailSendFailed     FlowError = "FlowError_EmailSendFailed"
	FlowError_LinkAccountFailed   FlowError = "FlowError_LinkAccountFailed"
	FlowError_CreateAccountFailed FlowError = "FlowError_CreateAccountFailed"
	FlowError_InternalError       FlowError = "FlowError_InternalError"
)
