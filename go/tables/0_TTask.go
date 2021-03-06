/**
 * Warning: Generated code! do not change!
 * Generated by: go/EntityDB.ftl
 */
package tables
import (
	"github.com/quintans/goSQL/db"
	"github.com/quintans/taskboard/go/entity"
	. "github.com/quintans/toolkit/ext"
)

var (
	TASK = db.TABLE("TASK")
	TASK_C_ID = TASK.KEY("ID")
	TASK_C_VERSION = TASK.VERSION("VERSION")
	// Audit
	TASK_C_CREATION = TASK.COLUMN("CREATION")
	TASK_C_MODIFICATION = TASK.COLUMN("MODIFICATION")
	TASK_C_USER_CREATION = TASK.COLUMN("USER_CREATION").As("UserCreationId")
	TASK_C_USER_MODIFICATION = TASK.COLUMN("USER_MODIFICATION").As("UserModificationId")
	// Atributos 
	TASK_C_TITLE = TASK.COLUMN("TITLE")
	TASK_C_DETAIL = TASK.COLUMN("DETAIL")
	TASK_C_HEAD_COLOR = TASK.COLUMN("HEAD_COLOR")
	TASK_C_BODY_COLOR = TASK.COLUMN("BODY_COLOR")
	TASK_C_POSITION = TASK.COLUMN("POSITION")
	TASK_C_REFERENCE = TASK.COLUMN("REFERENCE")
	TASK_C_SPENT = TASK.COLUMN("SPENT")
	TASK_C_REMAINING = TASK.COLUMN("REMAINING")
	// who is with this task
	TASK_C_USER_ID = TASK.COLUMN("USER").As("UserId")
	TASK_C_LANE_ID = TASK.COLUMN("LANE").As("LaneId")
	//FK's
	// who is with this task
	TASK_A_USER = TASK.
				ASSOCIATE(TASK_C_USER_ID).TO(USER_C_ID).
				As("User")
	TASK_A_LANE = TASK.
				ASSOCIATE(TASK_C_LANE_ID).TO(LANE_C_ID).
				As("Lane")
	TASK_A_NOTIFICATIONS = TASK.
				ASSOCIATE(TASK_C_ID).TO(NOTIFICATION_C_TASK_ID).
				As("Notifications")
	TASK_A_USER_CREATION = TASK.ASSOCIATE(TASK_C_USER_CREATION).TO(USER_C_ID).As("UserCreation")
	TASK_A_USER_MODIFICATION = TASK.ASSOCIATE(TASK_C_USER_MODIFICATION).TO(USER_C_ID).As("UserModification")
)

func init() {
	TASK.PreInsertTrigger = func(ins *db.Insert) {
		ins.Set(TASK_C_VERSION, 1)
		ins.Set(TASK_C_CREATION, NOW())
		uid, ok := ins.GetDb().GetAttribute(entity.ATTR_USERID)
		if ok {
			ins.Set(TASK_C_USER_CREATION, uid.(int64))
		}
	}
	TASK.PreUpdateTrigger = func(upd *db.Update) {
		upd.Set(TASK_C_MODIFICATION, NOW())
		uid, ok := upd.GetDb().GetAttribute(entity.ATTR_USERID)
		if ok {
			upd.Set(TASK_C_USER_MODIFICATION, uid.(int64))
		}
	}
}
