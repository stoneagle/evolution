package services

import "evolution/backend/time/models"

type General interface {
	One(int) (models.Task, error)
}
