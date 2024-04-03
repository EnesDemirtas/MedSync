package tagdb

import (
	"fmt"

	"github.com/EnesDemirtas/medisync/business/core/crud/tagbus"
	"github.com/google/uuid"
)

type dbTag struct {
	ID		uuid.UUID `db:"tag_id"`
	Name	string	  `db:"name"`
}

func toDBTag(tag tagbus.Tag) dbTag {
	tagDB := dbTag{
		ID:		tag.ID,
		Name:	tag.Name,
	}

	return tagDB
}

func toCoreTag(dbTag dbTag) (tagbus.Tag, error) {
	tag := tagbus.Tag{
		ID:		dbTag.ID,
		Name:	dbTag.Name,
	}

	return tag, nil
}

func toCoreTagSlice(dbTags []dbTag) ([]tagbus.Tag, error) {
	tags := make([]tagbus.Tag, len(dbTags))

	for i, dbTag := range dbTags {
		var err error
		tags[i], err = toCoreTag(dbTag)
		if err != nil {
			return nil, fmt.Errorf("parse type: %w", err)
		}
	}

	return tags, nil
}
