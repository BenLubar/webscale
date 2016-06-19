//go:generate go run mkid.go Permission

package model // import "github.com/BenLubar/webscale/model"

import "github.com/BenLubar/webscale/db"

type Permission struct {
	ID   PermissionID
	Slug string
}

const permissionFields = `p.id, p.slug`

func scanPermission(s scanner) (*Permission, error) {
	var p Permission
	if err := s.Scan(&p.ID, &p.Slug); err != nil {
		return nil, err
	}
	return &p, nil
}

var idGetPermission = db.Prepare(`select ` + permissionFields + ` from permissions as p where ` + noopPermission + ` and p.id = $3::bigint order by p.id asc;`)
var idsGetPermission = db.Prepare(`select ` + permissionFields + ` from permissions as p where ` + noopPermission + ` and array[p.id] <@ $3::bigint[] order by p.id asc;`)
