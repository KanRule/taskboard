/**
 * Initialy Generated by: go/EntityDAO.ftl
 */
package impl

import (
	dao "github.com/quintans/taskboard/biz/entity"
)

var _ dao.INotificationDAO = &NotificationDAO{}

func NewNotificationDAO(appCtx *AppCtx) dao.INotificationDAO {
	this := NotificationDAO{}
	this.Context = appCtx
	return this
}

type NotificationDAO struct {
	NotificationDAOBase
}
