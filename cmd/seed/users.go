package main

import (
	"errors"

	"github.com/eckertalex/improved-fiesta/internal/data"
)

type user struct {
	Username  string
	Email     string
	Password  string
	Activated bool
	Role      string
}

func (app *application) seedUsers() {
	admin := user{
		Username:  "admin",
		Email:     "admin@improved-fiesta.go",
		Password:  "admin123",
		Activated: true,
		Role:      data.AdminRole,
	}

	app.logger.Info("seeding admin user...")
	app.seedUser(&admin)
	app.logger.Info("done seeding admin user")

	activatedUser := user{
		Username:  "activated",
		Email:     "activated@improved-fiesta.go",
		Password:  "activated123",
		Activated: true,
		Role:      data.UserRole,
	}

	app.logger.Info("seeding activated user...")
	app.seedUser(&activatedUser)
	app.logger.Info("done seeding activated user")

	unactivatedUser := user{
		Username:  "unactivated",
		Email:     "unactivated@improved-fiesta.go",
		Password:  "unactivated123",
		Activated: false,
		Role:      data.UserRole,
	}

	app.logger.Info("seeding unactivated user...")
	app.seedUser(&unactivatedUser)
	app.logger.Info("done seeding unactivated user")
}

func (app *application) seedUser(user *user) {
	domainUser := &data.User{
		Username:  user.Username,
		Email:     user.Email,
		Activated: user.Activated,
		Role:      user.Role,
	}

	err := domainUser.Password.Set(user.Password)
	if err != nil {
		app.logger.Error(err.Error(), "failed to set password for user", domainUser.Username)
		return
	}

	err = app.models.Users.Insert(domainUser)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			app.logger.Error(err.Error(), "a user with this email address already exists", domainUser.Email)
		case errors.Is(err, data.ErrDuplicateUsername):
			app.logger.Error(err.Error(), "a user with this username already exists", domainUser.Username)
		default:
			app.logger.Error(err.Error(), "error inserting user", domainUser.Username)
		}
		return
	}

	app.logger.Info("successfully seeded user", "email", domainUser.Email)
}
