package main

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	Password string `json:"password"`
	Username string `json:"email"`
}

func NewApp() *fiber.App {
	instance := fiber.New()

	return instance
}

func NewDB() *gorm.DB {
	instance, err := gorm.Open(
		sqlite.Open("test.sqlite"),
		&gorm.Config{},
	)

	if err != nil {
		panic("failed to connect database")
	}

	return instance
}

func Migrate(
	db *gorm.DB,
) {
	db.AutoMigrate(&User{})
}

func Seed(
	db *gorm.DB,
) {
	db.Create(
		&User{
			Username: "admin",
			Password: "admin",
		},
	)
}

func SetupRoutes(
	app *fiber.App,
	db *gorm.DB,
) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Get("/users", func(c *fiber.Ctx) error {
		var users []User

		db.Find(&users)

		return c.JSON(users)
	})
}

func StartApp(
	app *fiber.App,
	lc fx.Lifecycle,
) {
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				fmt.Println("Starting server on port :3000")
				go app.Listen(":3000")
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return app.Shutdown()
			},
		},
	)
}

func main() {
	fx.New(
		fx.Provide(
			NewDB,
			NewApp,
		),
		fx.Invoke(
			Migrate,
			Seed,
			SetupRoutes,
			StartApp,
		),
	).Run()
}
