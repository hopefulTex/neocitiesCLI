package ui

import (
	"fmt"
	"neocitiesCli/api"
)

func InfoView(s api.SiteInfo) string {
	return fmt.Sprintf(
		`name: %s
hits: %d
created: %s
last update: %s
domain: %s
tags: %s`,
		s.Name, s.Hits, s.Created, s.LastUpdate, s.Domain, s.Tags)
}
