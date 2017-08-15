/**
 * Warning: Generated code! do not change!
 * Generated by: go/AppCtx.ftl
 */
package impl

import (
	"encoding/json"
	"net/http"

	"github.com/quintans/goSQL/db"
	"github.com/quintans/maze"
	"github.com/quintans/taskboard/go/service"
)

func NewAppCtx(
	w http.ResponseWriter,
	r *http.Request,
	filters []*maze.Filter,
	taskBoardService service.ITaskBoardService,
) *AppCtx {
	var this = new(AppCtx)
	this.Context = maze.NewContext(w, r, filters)
	this.taskBoardService = taskBoardService
	return this
}

var _ maze.IContext = &AppCtx{}

type AppCtx struct {
	*maze.Context

	Principal        *Principal
	Store            db.IDb
	taskBoardService service.ITaskBoardService
}

func (this *AppCtx) Proceed() error {
	return this.Next(this)
}

func (this *AppCtx) GetTaskBoardService() service.ITaskBoardService {
	return this.taskBoardService
}

func (this *AppCtx) SetTaskBoardService(taskBoardService service.ITaskBoardService) {
	this.taskBoardService = taskBoardService
}

func (this *AppCtx) BuildJsonRpcTaskBoardService(transaction func(ctx maze.IContext) error) *maze.JsonRpc {
	// JSON-RPC services
	var rpc = maze.NewJsonRpc(this.taskBoardService, transaction) // json-rpc builder
	rpc.SetActionFilters("WhoAmI", authorize("USER"))
	rpc.SetActionFilters("FetchBoardUsers", authorize("USER"))
	rpc.SetActionFilters("FetchBoardAllUsers", authorize("ADMIN"))
	rpc.SetActionFilters("FetchBoards", authorize("USER"))
	rpc.SetActionFilters("FetchBoardById", authorize("USER"))
	rpc.SetActionFilters("FullyLoadBoardById", authorize("USER"))
	rpc.SetActionFilters("SaveBoard", authorize("ADMIN"))
	rpc.SetActionFilters("DeleteBoard", authorize("ADMIN"))
	rpc.SetActionFilters("AddLane", authorize("ADMIN"))
	rpc.SetActionFilters("SaveLane", authorize("ADMIN"))
	rpc.SetActionFilters("DeleteLastLane", authorize("ADMIN"))
	rpc.SetActionFilters("SaveUser", authorize("ADMIN"))
	rpc.SetActionFilters("FetchUsers", authorize("ADMIN"))
	rpc.SetActionFilters("DeleteUser", authorize("ADMIN"))
	rpc.SetActionFilters("AddUserToBoard", authorize("ADMIN"))
	rpc.SetActionFilters("RemoveUserFromBoard", authorize("ADMIN"))
	rpc.SetActionFilters("SaveUserName", authorize("USER"))
	rpc.SetActionFilters("ChangeUserPassword", authorize("USER"))
	rpc.SetActionFilters("SaveTask", authorize("USER"))
	rpc.SetActionFilters("MoveTask", authorize("USER"))
	rpc.SetActionFilters("FetchNotifications", authorize("USER"))
	rpc.SetActionFilters("SaveNotification", authorize("USER"))
	rpc.SetActionFilters("DeleteNotification", authorize("USER"))
	return rpc
}

func (this *AppCtx) Reply(value interface{}) error {
	// JSON Vulnerability Protection.
	// AngularJS will automatically strip the prefix before processing it as JSON.
	result, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if _, err = this.Response.Write([]byte(")]}',\n")); err != nil {
		return err
	}
	_, err = this.Response.Write(result)
	return err
}

func authorize(roles ...string) func(ctx maze.IContext) error {
	return func(ctx maze.IContext) error {
		user := ctx.(*AppCtx).Principal
		for _, r := range roles {
			for _, role := range user.Roles {
				if r == string(role) {
					if err := ctx.Proceed(); err != nil {
						return err
					}
					return nil
				}
			}
		}
		http.Error(ctx.GetResponse(), "Unauthorized", http.StatusUnauthorized)
		return nil
	}
}
