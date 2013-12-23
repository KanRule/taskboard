/// <reference path="typings/angularjs/angular.d.ts"/>
/// <reference path="toolkit.ts"/>
var __extends = this.__extends || function (d, b) {
    for (var p in b) if (b.hasOwnProperty(p)) d[p] = b[p];
    function __() { this.constructor = d; }
    __.prototype = b.prototype;
    d.prototype = new __();
};
/**
* WARNING: Generated code. Changes will be overwritten.
* Generated by: ts/Entities.ftl
*/
var taskboard;
(function (taskboard) {
    var EntityBase = (function () {
        function EntityBase() {
        }
        EntityBase.prototype.copy = function (e) {
            this.id = e.id;
            this.version = e.version;
            this.creation = e.creation;
            this.modification = e.modification;
            return this;
        };
        return EntityBase;
    })();
    taskboard.EntityBase = EntityBase;

    var Board = (function (_super) {
        __extends(Board, _super);
        function Board() {
            _super.apply(this, arguments);
        }
        Board.prototype.clone = function () {
            var o = new Board();
            o.copy(this);
            return o;
        };

        Board.prototype.copy = function (e) {
            _super.prototype.copy.call(this, e);
            this.name = e.name;
            this.description = e.description;
            this.lanes = e.lanes;
            return this;
        };
        return Board;
    })(EntityBase);
    taskboard.Board = Board;

    var Lane = (function (_super) {
        __extends(Lane, _super);
        function Lane() {
            _super.apply(this, arguments);
        }
        Lane.prototype.clone = function () {
            var o = new Lane();
            o.copy(this);
            return o;
        };

        Lane.prototype.copy = function (e) {
            _super.prototype.copy.call(this, e);
            this.name = e.name;
            this.position = e.position;
            this.boardId = e.boardId;
            this.board = e.board;
            this.tasks = e.tasks;
            this.notifications = e.notifications;
            return this;
        };
        return Lane;
    })(EntityBase);
    taskboard.Lane = Lane;

    var Task = (function (_super) {
        __extends(Task, _super);
        function Task() {
            _super.apply(this, arguments);
        }
        Task.prototype.clone = function () {
            var o = new Task();
            o.copy(this);
            return o;
        };

        Task.prototype.copy = function (e) {
            _super.prototype.copy.call(this, e);
            this.responsible = e.responsible;
            this.title = e.title;
            this.detail = e.detail;
            this.headColor = e.headColor;
            this.bodyColor = e.bodyColor;
            this.position = e.position;
            this.laneId = e.laneId;
            this.lane = e.lane;
            this.notifications = e.notifications;
            return this;
        };
        return Task;
    })(EntityBase);
    taskboard.Task = Task;

    var Notification = (function (_super) {
        __extends(Notification, _super);
        function Notification() {
            _super.apply(this, arguments);
        }
        Notification.prototype.clone = function () {
            var o = new Notification();
            o.copy(this);
            return o;
        };

        Notification.prototype.copy = function (e) {
            _super.prototype.copy.call(this, e);
            this.email = e.email;
            this.taskId = e.taskId;
            this.laneId = e.laneId;
            this.task = e.task;
            this.lane = e.lane;
            return this;
        };
        return Notification;
    })(EntityBase);
    taskboard.Notification = Notification;
})(taskboard || (taskboard = {}));

/**
* WARNING: Generated code. Changes will be overwritten.
* Generated by: ts/DTO.ftl
*/
var taskboard;
(function (taskboard) {
    var BoardSearchDTO = (function (_super) {
        __extends(BoardSearchDTO, _super);
        function BoardSearchDTO() {
            _super.apply(this, arguments);
        }
        BoardSearchDTO.prototype.clone = function () {
            var o = new BoardSearchDTO();
            o.copy(this);
            return o;
        };

        BoardSearchDTO.prototype.copy = function (dto) {
            _super.prototype.copy.call(this, dto);
            this.name = dto.name;
            this.nif = dto.nif;
            return this;
        };

        BoardSearchDTO.prototype.reset = function () {
            this.page = 1;
            this.name = null;
            this.nif = null;
        };
        return BoardSearchDTO;
    })(toolkit.Criteria);
    taskboard.BoardSearchDTO = BoardSearchDTO;

    var IdVersionDTO = (function () {
        function IdVersionDTO() {
            this.version = 1;
        }
        IdVersionDTO.prototype.clone = function () {
            var o = new IdVersionDTO();
            o.copy(this);
            return o;
        };

        IdVersionDTO.prototype.copy = function (dto) {
            this.id = dto.id;
            this.version = dto.version;
            return this;
        };
        return IdVersionDTO;
    })();
    taskboard.IdVersionDTO = IdVersionDTO;

    var NotificationSearchDTO = (function (_super) {
        __extends(NotificationSearchDTO, _super);
        function NotificationSearchDTO() {
            _super.apply(this, arguments);
        }
        NotificationSearchDTO.prototype.clone = function () {
            var o = new NotificationSearchDTO();
            o.copy(this);
            return o;
        };

        NotificationSearchDTO.prototype.copy = function (dto) {
            _super.prototype.copy.call(this, dto);
            this.taskId = dto.taskId;
            this.email = dto.email;
            return this;
        };

        NotificationSearchDTO.prototype.reset = function () {
            this.page = 1;
            this.taskId = null;
            this.email = null;
        };
        return NotificationSearchDTO;
    })(toolkit.Criteria);
    taskboard.NotificationSearchDTO = NotificationSearchDTO;
})(taskboard || (taskboard = {}));

/**
* WARNING: Generated code. Changes will be overwritten.
* Generated by: ts/ServiceOperationIn.ftl
*/
var taskboard;
(function (taskboard) {
    var MoveTaskIn = (function () {
        function MoveTaskIn() {
        }
        return MoveTaskIn;
    })();
    taskboard.MoveTaskIn = MoveTaskIn;
})(taskboard || (taskboard = {}));

/**
* WARNING: Generated code. Changes will be overwritten.
* Generated by: ts/Services.ftl
*/
var taskboard;
(function (taskboard) {
    function defaultErrorHandler(fail, status) {
        if (status == 500) {
            var msg = fail.error + ": " + fail.message;
            toolkit.stickyError("Server Error", msg);
        }
    }

    var TaskBoardService = (function () {
        function TaskBoardService(http) {
            this.http = http;
        }
        TaskBoardService.prototype.fetchBoards = function (criteria, successCallback, errorCallback) {
            var payload;
            if (criteria != null)
                payload = criteria;
else
                payload = "null";

            var promisse = this.http.post("rest/taskboard/FetchBoards", payload);
            if (successCallback != null) {
                promisse.success(successCallback);
            }
            promisse.error(errorCallback != null ? errorCallback : defaultErrorHandler);
        };

        TaskBoardService.prototype.fetchBoardById = function (id, successCallback, errorCallback) {
            var payload;
            if (id != null)
                payload = "" + id;
else
                payload = "null";

            var promisse = this.http.post("rest/taskboard/FetchBoardById", payload);
            if (successCallback != null) {
                promisse.success(successCallback);
            }
            promisse.error(errorCallback != null ? errorCallback : defaultErrorHandler);
        };

        TaskBoardService.prototype.fullyLoadBoardById = function (id, successCallback, errorCallback) {
            var payload;
            if (id != null)
                payload = "" + id;
else
                payload = "null";

            var promisse = this.http.post("rest/taskboard/FullyLoadBoardById", payload);
            if (successCallback != null) {
                promisse.success(successCallback);
            }
            promisse.error(errorCallback != null ? errorCallback : defaultErrorHandler);
        };

        TaskBoardService.prototype.saveBoard = function (board, successCallback, errorCallback) {
            var payload;
            if (board != null)
                payload = board;
else
                payload = "null";

            var promisse = this.http.post("rest/taskboard/SaveBoard", payload);
            if (successCallback != null) {
                promisse.success(successCallback);
            }
            promisse.error(errorCallback != null ? errorCallback : defaultErrorHandler);
        };

        TaskBoardService.prototype.deleteBoard = function (id, successCallback, errorCallback) {
            var payload;
            if (id != null)
                payload = "" + id;
else
                payload = "null";

            var promisse = this.http.post("rest/taskboard/DeleteBoard", payload);
            if (successCallback != null) {
                promisse.success(successCallback);
            }
            promisse.error(errorCallback != null ? errorCallback : defaultErrorHandler);
        };

        TaskBoardService.prototype.addLane = function (boardId, successCallback, errorCallback) {
            var payload;
            if (boardId != null)
                payload = "" + boardId;
else
                payload = "null";

            var promisse = this.http.post("rest/taskboard/AddLane", payload);
            if (successCallback != null) {
                promisse.success(successCallback);
            }
            promisse.error(errorCallback != null ? errorCallback : defaultErrorHandler);
        };

        TaskBoardService.prototype.saveLane = function (lane, successCallback, errorCallback) {
            var payload;
            if (lane != null)
                payload = lane;
else
                payload = "null";

            var promisse = this.http.post("rest/taskboard/SaveLane", payload);
            if (successCallback != null) {
                promisse.success(successCallback);
            }
            promisse.error(errorCallback != null ? errorCallback : defaultErrorHandler);
        };

        TaskBoardService.prototype.deleteLastLane = function (boardId, successCallback, errorCallback) {
            var payload;
            if (boardId != null)
                payload = "" + boardId;
else
                payload = "null";

            var promisse = this.http.post("rest/taskboard/DeleteLastLane", payload);
            if (successCallback != null) {
                promisse.success(successCallback);
            }
            promisse.error(errorCallback != null ? errorCallback : defaultErrorHandler);
        };

        TaskBoardService.prototype.saveTask = function (task, successCallback, errorCallback) {
            var payload;
            if (task != null)
                payload = task;
else
                payload = "null";

            var promisse = this.http.post("rest/taskboard/SaveTask", payload);
            if (successCallback != null) {
                promisse.success(successCallback);
            }
            promisse.error(errorCallback != null ? errorCallback : defaultErrorHandler);
        };

        TaskBoardService.prototype.moveTask = function (input, successCallback, errorCallback) {
            var promisse = this.http.post("rest/taskboard/MoveTask", input);
            if (successCallback != null) {
                promisse.success(successCallback);
            }
            promisse.error(errorCallback != null ? errorCallback : defaultErrorHandler);
        };

        TaskBoardService.prototype.fetchNotifications = function (criteria, successCallback, errorCallback) {
            var payload;
            if (criteria != null)
                payload = criteria;
else
                payload = "null";

            var promisse = this.http.post("rest/taskboard/FetchNotifications", payload);
            if (successCallback != null) {
                promisse.success(successCallback);
            }
            promisse.error(errorCallback != null ? errorCallback : defaultErrorHandler);
        };

        TaskBoardService.prototype.saveNotification = function (notification, successCallback, errorCallback) {
            var payload;
            if (notification != null)
                payload = notification;
else
                payload = "null";

            var promisse = this.http.post("rest/taskboard/SaveNotification", payload);
            if (successCallback != null) {
                promisse.success(successCallback);
            }
            promisse.error(errorCallback != null ? errorCallback : defaultErrorHandler);
        };

        TaskBoardService.prototype.deleteNotification = function (id, successCallback, errorCallback) {
            var payload;
            if (id != null)
                payload = "" + id;
else
                payload = "null";

            var promisse = this.http.post("rest/taskboard/DeleteNotification", payload);
            if (successCallback != null) {
                promisse.success(successCallback);
            }
            promisse.error(errorCallback != null ? errorCallback : defaultErrorHandler);
        };
        return TaskBoardService;
    })();
    taskboard.TaskBoardService = TaskBoardService;

    angular.module("remoteServices", []).factory("taskBoardService", [
        "$http",
        function ($http) {
            return new TaskBoardService($http);
        }
    ]);
})(taskboard || (taskboard = {}));
