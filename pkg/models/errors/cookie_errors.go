package error

type (
	FlowError string
)

const (
	FlowError_AutoAccountCreationNotAllowed FlowError = "FlowError_AutoAccountCreationNotAllowed"
	FlowError_BadSession                    FlowError = "FlowError_BadSession"
	// ... more errors
)
