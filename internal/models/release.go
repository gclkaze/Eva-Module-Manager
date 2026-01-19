package models

import "time"

type Release struct {
	ID          uint       `json:"id" binding:"required"`
	Version     string     `json:"version" binding:"required"`
	ReleasedAt  *time.Time `json:"released_at" binding:"required"`
	Description string     `json:"description" binding:"required"`
	DiskSize    int64      `json:"diskSize" binding:"required"`
	Status      string     `json:"status" binding:"required"`
	Keywords    []Keyword  `json:"keywords" binding:"required"`
}
