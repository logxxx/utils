package commproto

type ReportReq struct {
	ReportTime int64                  `json:"report_time,omitempty"`
	DeviceID   string                 `json:"device_id,omitempty"`
	Action     string                 `json:"action,omitempty"`
	Payload    map[string]interface{} `json:"payload"`
	WorkDir    string                 `json:"work_dir"`
	HostName   string                 `json:"host_name"`
	HostUser   string                 `json:"host_user"`
	Pid        int                    `json:"pid"`
	IP         string                 `json:"ip"`
	*Statistic
}

type CommonReq struct {
	Body       []byte `json:"body,omitempty"`
	CreateTime int64  `json:"create_time,omitempty"`
}
