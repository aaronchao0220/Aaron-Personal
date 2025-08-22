package model

// ScrubbedEvent represents a telemetry event with sensitive information removed.
type ScrubbedEvent struct {
	Id                 string
	SpecVersion        string
	TenantId           string
	UserId             string
	SessionId          string
	Source             string
	Type               string
	Time               string
	Host               string
	OriginIp           string
	OwnerId            string
	TopLevelResourceId string
	SpaceId            string
	ClientId           string
	Reason             string
	Data               map[string]any
}
