package commands

import (
	"context"
	"fmt"
	"net/mail"
	"time"

	"github.com/EnesDemirtas/medisync/business/data/sqldb"
	"github.com/EnesDemirtas/medisync/business/domain/userbus"
	"github.com/EnesDemirtas/medisync/business/domain/userbus/stores/userdb"
	"github.com/EnesDemirtas/medisync/foundation/logger"
)

// UserAdd adds new users into the database.
func UserAdd(log *logger.Logger, cfg sqldb.Config, name, email, password string) error {
	if name == "" || email == "" || password == "" {
		fmt.Println("help: useradd <name> <email> <password>")
		return ErrHelp
	}

	db, err := sqldb.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userBus := userbus.NewCore(log, nil, userdb.NewStore(log, db))

	addr, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("parsing email: %w", err)
	}

	nu := userbus.NewUser{
		Name:            name,
		Email:           *addr,
		Password:        password,
		PasswordConfirm: password,
		Roles:           []userbus.Role{userbus.RoleAdmin, userbus.RoleUser},
	}

	usr, err := userBus.Create(ctx, nu)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	fmt.Println("user id:", usr.ID)
	return nil
}
