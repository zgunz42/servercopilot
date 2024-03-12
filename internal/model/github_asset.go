// Asset.go

package model

type Asset struct {
	URL                string `json:"url"`
	ID                 int64  `json:"id"`
	NodeID             string `json:"node_id"`
	Name               string `json:"name"`
	Label              string `json:"label"`
	Uploader           Author `json:"uploader"`
	ContentType        string `json:"content_type"`
	State              string `json:"state"`
	Size               int64  `json:"size"`
	DownloadCount      int64  `json:"download_count"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
	BrowserDownloadURL string `json:"browser_download_url"`
}
