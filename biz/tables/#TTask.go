/**
 * Warning: Generated code! do not change!
 * Generated by: go/EntityDB.ftl
 */
package entity;

import (
	"github.com/quintans/goSQL/db"
)

var (
	TASK = db.TABLE("TASK")
	TASK_C_ID = TASK.KEY("ID")
	TASK_C_VERSION = TASK.VERSION("VERSION")
	// Audit
	TASK_C_CREATION = TASK.COLUMN("CREATION")
	TASK_C_MODIFICATION = TASK.COLUMN("MODIFICATION")
	// Atributos 
	// who is with this task
	TASK_C_RESPONSIBLE = TASK.COLUMN("RESPONSIBLE")
	TASK_C_TITLE = TASK.COLUMN("TITLE")
	TASK_C_DETAIL = TASK.COLUMN("DETAIL")
	TASK_C_HEAD_COLOR = TASK.COLUMN("HEAD_COLOR")
	TASK_C_BODY_COLOR = TASK.COLUMN("BODY_COLOR")
	TASK_C_POSITION = TASK.COLUMN("POSITION")
	TASK_C_LANE_ID = TASK.COLUMN("LANE").As("LaneId")
	//FK's
	TASK_A_LANE = TASK.
				ASSOCIATE(TASK_C_LANE_ID).TO(LANE_C_ID).
				As("Lane")
	TASK_A_NOTIFICATIONS = TASK.
				ASSOCIATE(TASK_C_ID).TO(NOTIFICATION_C_TASK_ID).
				As("Notifications")
)
