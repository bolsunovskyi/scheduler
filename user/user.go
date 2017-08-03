package user

import (
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type Model struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Password  string
	Role      string `gorm:"column:user_role"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Model) TableName() string {
	return "user"
}

type Session struct {
	ID        int
	UserID    int
	Value     string
	IP        string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Session) TableName() string {
	return "session"
}

func createAdmin(db *gorm.DB, email, password string) error {
	pwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	var admin Model
	if err := db.Where("email = ?", email).First(&admin).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}

		return db.Save(&Model{
			Password:  string(pwd),
			Email:     email,
			UpdatedAt: time.Now(),
			CreatedAt: time.Now(),
			Role:      "admin",
		}).Error
	}

	admin.Password = string(pwd)
	admin.UpdatedAt = time.Now()

	return db.Save(&admin).Error
}
