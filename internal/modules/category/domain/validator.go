package domain

import "strings"

func (c *Category) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return ErrNameRequired
	}
	return nil
}
