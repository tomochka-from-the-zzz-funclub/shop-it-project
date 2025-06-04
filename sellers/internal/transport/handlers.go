package transport

import (
	"fmt"
	"net/http"

	"market/internal/config"
	"market/internal/logger"
	"market/internal/models/request"
	"market/internal/services/auth"
	sellers "market/internal/services/sellers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type HandlersBuilder struct {
	srv         sellers.ISellers
	authService *auth.AuthService
	app         *fiber.App
	log         *logger.Logger
}

func StartServer(cfg *config.Options, srv sellers.ISellers, authService *auth.AuthService, log *logger.Logger) {
	hb := HandlersBuilder{
		srv:         srv,
		authService: authService,
		app:         fiber.New(fiber.Config{Immutable: true}),
		log:         log,
	}

	// Start Prometheus metrics server
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":8090", nil); err != nil {
			log.LogErrorf("Failed to start metrics server: %v", err)
		}
	}()

	// Configure middleware
	hb.app.Use(cors.New(cors.Config{
		AllowOrigins:  "*",
		AllowMethods:  "GET,POST,PUT,DELETE",
		AllowHeaders:  "Authorization,Content-Type,X-Request-Id",
		ExposeHeaders: "X-Request-Id",
	}))
	hb.app.Use(requestid.New(requestid.Config{
		Header: "X-Request-Id",
	}))

	// Authorization middleware
	hb.app.Use(func(c *fiber.Ctx) error {
		if c.Path() == "/api/v1/auth/register" || c.Path() == "/api/v1/auth/login" {
			return c.Next()
		}
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing authorization header"})
		}
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}
		sellerID, err := authService.ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}
		c.Locals("sellerID", sellerID)
		return c.Next()
	})

	// Define routes
	api := hb.app.Group("/api/v1")
	api.Post("/auth/register", hb.HandleRegister)
	api.Post("/auth/login", hb.HandleLogin)
	api.Put("/sellers/update/:uuid", hb.HandleUpdateSeller)
	api.Delete("/sellers/delete/:uuid", hb.HandleDeleteSeller)

	// Start server
	if err := hb.app.Listen(":8084"); err != nil {
		log.LogErrorf("Failed to start server: %v", err)
	}
}

func (hb *HandlersBuilder) HandleRegister(c *fiber.Ctx) error {
	hb.log.LogDebugf("Start func HandleRegister")
	var req struct {
		Email      string         `json:"email"`
		Password   string         `json:"password"`
		SellerInfo request.Seller `json:"seller_info"`
	}
	if err := c.BodyParser(&req); err != nil {
		hb.log.LogErrorf("Invalid request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	sellerID, err := hb.authService.Register(c.Context(), req.Email, req.Password, req.SellerInfo)
	if err != nil {
		hb.log.LogErrorf("Failed to register seller with email %s: %v", req.Email, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to register seller"})
	}

	hb.log.LogInfof("Successfully registered seller %s", sellerID)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": sellerID.String()})
}

func (hb *HandlersBuilder) HandleLogin(c *fiber.Ctx) error {
	hb.log.LogDebugf("Start func HandleLogin")
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		hb.log.LogErrorf("Invalid request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	token, err := hb.authService.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		hb.log.LogErrorf("Failed to login seller with email %s: %v", req.Email, err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	hb.log.LogInfof("Successfully logged in seller with email %s", req.Email)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}

func (hb *HandlersBuilder) HandleUpdateSeller(c *fiber.Ctx) error {
	hb.log.LogDebugf("Start func HandleUpdateSeller")
	var seller request.Seller
	if err := c.BodyParser(&seller); err != nil {
		hb.log.LogErrorf("Invalid request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	idStr := c.Params("uuid")
	sellerID, err := uuid.Parse(idStr)
	if err != nil {
		hb.log.LogErrorf("Invalid seller ID: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid seller ID"})
	}

	// Check if the authenticated seller matches the requested ID
	authSellerID, ok := c.Locals("sellerID").(uuid.UUID)
	if !ok || authSellerID != sellerID {
		hb.log.LogErrorf("Seller %s attempted to update seller %s", authSellerID, sellerID)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	if err := hb.srv.UpdateSeller(c.Context(), sellerID, seller); err != nil {
		hb.log.LogErrorf("Failed to update seller %s: %v", sellerID, err)
		if err.Error() == "seller not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("seller with ID %s not found", sellerID)})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update seller"})
	}

	hb.log.LogInfof("Successfully updated seller %s", sellerID)
	return c.Status(fiber.StatusNoContent).Send(nil)
}

func (hb *HandlersBuilder) HandleDeleteSeller(c *fiber.Ctx) error {
	hb.log.LogDebugf("Start func HandleDeleteSeller")
	idStr := c.Params("uuid")
	sellerID, err := uuid.Parse(idStr)
	if err != nil {
		hb.log.LogErrorf("Invalid seller ID: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid seller ID"})
	}

	// Check if the authenticated seller matches the requested ID
	authSellerID, ok := c.Locals("sellerID").(uuid.UUID)
	if !ok || authSellerID != sellerID {
		hb.log.LogErrorf("Seller %s attempted to delete seller %s", authSellerID, sellerID)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	if err := hb.srv.DeleteSeller(c.Context(), sellerID); err != nil {
		hb.log.LogErrorf("Failed to delete seller %s: %v", sellerID, err)
		if err.Error() == "seller not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("seller with ID %s not found", sellerID)})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete seller"})
	}

	hb.log.LogInfof("Successfully deleted seller %s", sellerID)
	return c.Status(fiber.StatusNoContent).Send(nil)
}
