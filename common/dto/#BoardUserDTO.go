/**
 * WARNING: Generated code! do not change!
 * Generated by: go/DTO.ftl
 */
package dto

import (
)

func NewBoardUserDTO() BoardUserDTO {
	this := BoardUserDTO{}
	return this
}

type BoardUserDTO struct {
	IdVersionDTO
	
	Name *string `json:"name"`
	Belongs bool `json:"belongs"`
}