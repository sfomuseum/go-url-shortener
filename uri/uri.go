package uri

type URI struct {
	Short   string `json:"short"`
	Source  string `json:"source"`
	Created int64  `json:"created"`
	Author  string `json:"author"`
}
