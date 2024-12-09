package main

import (
	"Program/database"
	"Program/generateToken"
	"Program/refresh"
	"Program/utilits"
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"sync"
)

func main() {
	// Подключение к базе данных
	tokenMap := &sync.Map{}

	connStr := "user=postgres password=efim dbname=token sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	// Закртыие базы данных
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	// Создание таблицы с refresh токенами
	database.CreateTableRefreshTokens(db)

	privateKey := []byte("secret-key")

	app := fiber.New()

	app.Get("user/getToken/:guid", func(c *fiber.Ctx) error {
		userGUID, err := utilits.ParseGUID(c)
		accessToken, refreshToken, err := generateToken.GenerateAndSaveTokens(c, userGUID, db, privateKey, tokenMap)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		})
	})

	app.Get("user/refreshToken/:refreshToken", func(c *fiber.Ctx) error {
		refreshToken := c.Params("refreshToken")
		hasher := sha512.New()
		_, err := hasher.Write([]byte(refreshToken))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		hashedToken := []byte(hex.EncodeToString(hasher.Sum(nil)))

		userGUID, ipList, err := utilits.IpToTokens(db, hashedToken)
		ip := c.IP()
		useIP := false
		for _, oldIP := range ipList {
			if oldIP == ip {
				useIP = true
				break
			}
		}

		if !useIP {
			// Отправляем сообщение на email
			println("New ip!!!")
		}

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		err = refresh.DeleteRefreshToken(db, hashedToken, tokenMap)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		accessToken, refreshToken, err := generateToken.GenerateAndSaveTokens(c, userGUID, db, privateKey, tokenMap)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		})
	})

	err = app.Listen(":3000")
	if err != nil {
		return
	}
}
