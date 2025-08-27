package dto

type CreateSecurityEventRequest struct {
	EventType   string `json:"event_type" binding:"required"`
	Description string `json:"description" binding:"required"`
	IPAddress   string `json:"ip_address"`
	UserID      *uint  `json:"user_id"`
	UserAgent   string `json:"user_agent"`
	Metadata    string `json:"metadata"`
}

type CreateRateLimitRuleRequest struct {
	RuleType    string `json:"rule_type" binding:"required"`
	Identifier  string `json:"identifier" binding:"required"`
	MaxRequests int64  `json:"max_requests" binding:"required"`
	WindowSize  int64  `json:"window_size" binding:"required"`
}

type UpdateRateLimitRuleRequest struct {
	RuleType    string `json:"rule_type"`
	Identifier  string `json:"identifier"`
	MaxRequests int64  `json:"max_requests"`
	WindowSize  int64  `json:"window_size"`
}
