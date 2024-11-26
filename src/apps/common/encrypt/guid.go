package encrypt

import "github.com/rs/xid"

func NewGUID() string {
	return xid.New().String()
}
