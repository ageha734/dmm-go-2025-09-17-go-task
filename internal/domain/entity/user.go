package entity

import (
	"time"
)

type User struct {
	ID        uint
	Name      string
	Email     string
	Age       int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(name, email string, age int) *User {
	now := time.Now()
	return &User{
		Name:      name,
		Email:     email,
		Age:       age,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (u *User) UpdateProfile(name string, age int) {
	if name != "" {
		u.Name = name
	}
	if age > 0 {
		u.Age = age
	}
	u.UpdatedAt = time.Now()
}

func (u *User) IsValidAge() bool {
	return u.Age >= 0 && u.Age <= 150
}

func (u *User) IsValidEmail() bool {
	return u.Email != "" && len(u.Email) > 0
}
