package api

import "fmt"

// GET "/api/list"
// GET "/api/list?path=path/to/dir"
type ListItem struct {
	Path       string `json:"path"`
	IsDir      bool   `json:"is_directory"`
	Size       int    `json:"size"`
	LastUpdate string `json:"updated_at"` // RFC 2822
	Hash       string `json:"hash"`       // SHA1
}

func (l ListItem) View() string {
	return fmt.Sprintf("%s\t%t\t%d\t%s", l.Path, l.IsDir, l.Size, l.LastUpdate)
}

// GET "/api/list"
// GET "/api/list?path=path/to/dir"
func (c *Connection) List(path string) ([]ListItem, error) {
	var params []string
	if path != "" {
		params = append(params, path)
	}
	response, err := c.Request(GET, "list", params, nil)
	if err != nil {
		return nil, err
	}

	return response.Files, nil
}
