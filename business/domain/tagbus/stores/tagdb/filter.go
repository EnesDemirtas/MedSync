package tagdb

import (
	"bytes"
	"strings"

	"github.com/EnesDemirtas/medisync/business/core/crud/tagbus"
)

func (s *Store) applyFilter(filter tagbus.QueryFilter, data map[string]interface{}, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["id"] = *filter.ID
		wc = append(wc, "tag_id = :tag_id")
	}

	if filter.Name != nil {
		data["name"] = *filter.Name
		wc = append(wc, "name = :name")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}