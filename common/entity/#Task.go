/**
 * WARNING: Generated code! do not change!
 * Generated by: go/Entity.ftl
 */
package entity

import (
	"github.com/quintans/toolkit"
	"github.com/quintans/toolkit/ext"
	"github.com/quintans/toolkit/web/app"
)

var _ toolkit.Hasher = &Task{}

func NewTask() *Task {
	this := new(Task)
	return this
}

type Task struct {
	EntityAudit

	//ATTRIBUTES
	Title     *string `json:"title"`
	Detail    *string `json:"detail"`
	HeadColor *string `json:"headColor"`
	BodyColor *string `json:"bodyColor"`
	Position  *int64  `json:"position"`
	UserId    *int64  `json:"userId"`
	LaneId    *int64  `json:"laneId"`
	// ASSOCIATIONS
	// user
	User *User `json:"user"`
	// lane
	Lane *Lane `json:"lane"`
	// notifications
	Notifications []*Notification `json:"notifications"`
}

func (this *Task) Clone() interface{} {
	clone := NewTask()
	clone.Copy(this)
	return clone
}

func (this *Task) Copy(entity *Task) {
	if entity != nil {
		this.EntityAudit.Copy(entity.EntityAudit)
		// attributes
		this.Title = app.CopyString(entity.Title)
		this.Detail = app.CopyString(entity.Detail)
		this.HeadColor = app.CopyString(entity.HeadColor)
		this.BodyColor = app.CopyString(entity.BodyColor)
		this.Position = app.CopyInteger(entity.Position)
		// associations
		this.User = entity.User
		this.Lane = entity.Lane
		this.Notifications = make([]*Notification, len(entity.Notifications), cap(entity.Notifications))
		copy(this.Notifications, entity.Notifications)
	}
}

func (this *Task) String() string {
	sb := toolkit.NewStrBuffer()
	sb.Add("{Id: ", this.Id, ", Version: ", this.Version)
	sb.Add(", Title: ", this.Title)
	sb.Add(", Detail: ", this.Detail)
	sb.Add(", HeadColor: ", this.HeadColor)
	sb.Add(", BodyColor: ", this.BodyColor)
	sb.Add(", Position: ", this.Position)
	sb.Add(", UserId: ", this.UserId)
	sb.Add(", *LaneId: ", *this.LaneId)
	sb.Add("}")
	return sb.String()
}

func (this *Task) Equals(e interface{}) bool {
	if this == e {
		return true
	}

	switch t := e.(type) {
	case *Task:
		return this.Id != nil && t.Id != nil && *this.Id == *t.Id
	}
	return false
}

func (this *Task) HashCode() int {
	result := toolkit.HashType(toolkit.HASH_SEED, this)
	result = toolkit.HashLong(result, ext.DefInt64(this.Id, 0))
	return result
}
