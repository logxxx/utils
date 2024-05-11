package commproto

type Statistic struct {
	StartupTime          int64        `json:"startup_time,omitempty"`
	DownloadRecord       int          `json:"download_record,omitempty"`
	DownloadImage        int          `json:"download_image,omitempty"`
	DownloadVideo        int          `json:"download_video,omitempty"`
	DownloadVideoLive    int          `json:"download_video_live,omitempty"`
	DownloadBytes        int64        `json:"download_bytes,omitempty"`
	DownloadBytesForShow string       `json:"download_bytes_for_show,omitempty"`
	LastDownload         LastDownload `json:"last_download"`
	LastErrMsgs          []string     `json:"last_err_msgs,omitempty"`
	ClientVersion        string       `json:"client_version,omitempty"`
	ActiveStatus         string       `json:"active_status,omitempty"`
	BulletinBoard        string       `json:"bulletin_board,omitempty"`
}

type LastDownload struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	UserName     string `json:"user_name"`
	DownloadTo   string `json:"download_to"`
	Size         int64  `json:"size"`
	DownloadTime int64  `json:"download_time"`
}
