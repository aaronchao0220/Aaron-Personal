package model

// CloudEvent struct for CloudEvent 1.0
type CloudEvent struct {
	Id                 string         `json:"id" bson:"event_id,omitempty"`
	SpecVersion        string         `json:"specversion" bson:"specversion,omitempty"`
	TenantId           string         `json:"tenantid" bson:"tenant_id,omitempty"`
	UserId             string         `json:"userid" bson:"user_id,omitempty"`
	SessionId          string         `json:"sessionid" bson:"session_id,omitempty"`
	Source             string         `json:"source" bson:"source,omitempty"`
	EventType          string         `json:"type" bson:"event_type,omitempty"`
	Time               string         `json:"time" bson:"event_time,omitempty"`
	Host               string         `json:"host" bson:"host,omitempty"`
	OriginIp           string         `json:"originip" bson:"origin_ip,omitempty"`
	OwnerId            string         `json:"ownerid" bson:"owner_id,omitempty"`
	TopLevelResourceId string         `json:"toplevelresourceid" bson:"top_level_resource_id,omitempty"`
	SpaceId            string         `json:"spaceid" bson:"space_id,omitempty"`
	ClientId           string         `json:"clientid" bson:"client_id,omitempty"`
	Reason             string         `json:"reason" bson:"reason,omitempty"`
	Data               map[string]any `json:"data" bson:"data,omitempty"`
}
