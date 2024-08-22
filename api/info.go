package api

import (
	"fmt"
	"strings"
)

// GET "/api/info"
// GET "/api/info?sitename=sitename"
type SiteInfo struct {
	Name       string   `json:"sitename"`
	Hits       int      `json:"hits"`
	Created    string   `json:"created_at"` // RFC 2822
	LastUpdate string   `json:"updated_at"` // RFC 2822
	Domain     string   `json:"domain"`
	Tags       []string `json:"tags"`
}

// GET "/api/info"
// GET "/api/info?sitename=sitename"
func (c *Connection) Info(sitename string) (SiteInfo, error) {
	fmt.Println("----Info----")

	sitename = strings.TrimPrefix(sitename, "https://")
	sitename = strings.TrimPrefix(sitename, "http://")
	sitename = strings.TrimPrefix(sitename, "www.")
	sitename = strings.TrimSuffix(sitename, "/")
	sitename = strings.TrimSuffix(sitename, ".neocities.org")

	var params []string
	if sitename != "" {
		params = []string{"sitename=" + sitename}
	}
	response, err := c.Request(GET, "info", params, nil)
	if err != nil {
		return SiteInfo{}, err
	}
	return response.Info, nil
}
