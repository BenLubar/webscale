//go:generate go run mkid.go Group

package model // import "github.com/BenLubar/webscale/model"

import "github.com/BenLubar/webscale/db"

type Group struct {
	ID   GroupID
	Name string
	Slug string
}

const groupFields = `g.id, g.name, g.slug`

func scanGroup(s scanner) (*Group, error) {
	var g Group
	if err := s.Scan(&g.ID, &g.Name, &g.Slug); err != nil {
		return nil, err
	}
	return &g, nil
}

var idGetGroup = db.Prepare(`select ` + groupFields + ` from groups as g where can_group($1::bigint, 'group-meta', $2::boolean, g.id) and g.id = $3::bigint order by g.id asc;`)
var idsGetGroup = db.Prepare(`select ` + groupFields + ` from groups as g where can_group($1::bigint, 'group-meta', $2::boolean, g.id) and array[g.id] <@ $3::bigint[] order by g.id asc;`)
