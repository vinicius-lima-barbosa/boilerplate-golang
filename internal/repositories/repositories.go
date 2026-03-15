package repositories

import "gorm.io/gorm"

type Repositories struct {
	UserRepository *UserRepository
}

func New(db *gorm.DB) *Repositories {
	return &Repositories{
		UserRepository: &UserRepository{db: db},
	}
}
