package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Структура для входящего JSON
type RequestBody struct {
	IDrequest int `json:"id"`
}

// Структура для ответа
type ResponseBody struct {
	Result string `json:"result"`
}

func main() {
	app := fiber.New()

	// Добавляем middleware для проверки заголовков
	app.Use(rpcMiddleware)

	app.Post("/rpc", handleRPC)

	log.Fatal(app.Listen("127.0.0.1:8080"))
}

// Проверка метода и заголовка
func rpcMiddleware(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodPost {
		return c.Status(http.StatusMethodNotAllowed).SendString("Method Not Allowed")
	}
	if c.Get("X-Service-Account") != "service-account" || c.Get("X-Service-Credentials") == "" {
		return c.Status(http.StatusUnauthorized).SendString("Unauthorized")
	}
	return c.Next()
}

// Основной обработчик RPC запросов
func handleRPC(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), time.Second)
	defer cancel()

	// Проверка размера тела
	if len(c.Body()) > 32 {
		return c.Status(http.StatusBadRequest).SendString("Request body too large")
	}

	// Парсинг JSON
	var req RequestBody
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON",
		})
	}

	// Обработка запроса
	method := c.Get("X-Rpc-Method")
	if method == "" || !isValidMethod(method) {
		return c.Status(http.StatusBadRequest).SendString("Invalid X-Rpc Method")
	}

	// Таймаут обработки
	done := make(chan struct{})
	var result string
	var err error

	go func() {
		result, err = handleRPCRequest(method, int(req.IDrequest))
		close(done)
	}()
	// Обработка времени, истечение таймаута
	select {
	case <-ctx.Done(): // Возвращаем ошибку о завершении времени
		return c.Status(http.StatusRequestTimeout).SendString("Request timeout")
	case <-done:
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(ResponseBody{
			Result: result,
		})
	}
}

// Проверка валидности метода
func isValidMethod(method string) bool {
	return method == "user.username" ||
		method == "user.profile" ||
		method == "user.fullname" ||
		method == "user.group"
}

// Обработчик RPC запросов
func handleRPCRequest(method string, IDrequest int) (string, error) {
	switch method {
	case "user.username":
		return "username", nil
	case "user.profile":
		return "profile", nil
	case "user.fullname":
		return "fullname", nil
	case "user.group":
		return "group", nil
	default:
		return "", fmt.Errorf("unknown method: %s", method)
	}
}
