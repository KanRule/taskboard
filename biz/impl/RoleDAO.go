/**
 * Initialy Generated by: go/EntityDAO.ftl
 */
package impl

import (
	dao "github.com/quintans/taskboard/biz/entity"
)

var _ dao.IRoleDAO = &RoleDAO{}

func NewRoleDAO(appCtx *AppCtx) dao.IRoleDAO {
	this := RoleDAO{}
	this.Context = appCtx
	return this
}

type RoleDAO struct {
	RoleDAOBase
}
