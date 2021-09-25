package types

import "strings"

type Toggle struct {
	Key string
	UsersDocument struct {
		Package string
		Name string
	}
}

func (t *Toggle) GetOpaEndpoint() string {
	return "/v1/data/" + strings.ReplaceAll(t.UsersDocument.Package, ".", "/") + "/" + t.UsersDocument.Name
}

