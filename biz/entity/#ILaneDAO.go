/**
 * WARNING: Generated code! do not change!
 * Generated by: go/IEntityDAO.ftl
 */
package entity

import (
	"github.com/quintans/taskboard/common/entity"
)

type ILaneDAO interface {
	Delete(entity *entity.Lane) error
    DeleteById(id int64) (bool, error)
    DeleteByIdAndVersion(id int64, version int64) error
    Save(entity *entity.Lane) error
    FindById(id int64) (*entity.Lane, error)
    FindAll() ([]*entity.Lane, error)
    FindByPos(boardId *int64, position *int64) (*entity.Lane, error)
}
