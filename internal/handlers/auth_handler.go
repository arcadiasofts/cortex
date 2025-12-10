package handlers

import (
	"backend/internal/services"
	"backend/server/pb"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/protobuf/proto"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(s *services.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) RequestChallenge(c *fiber.Ctx) error {
	type Request struct {
		did string `json:"did"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	nonce, err := h.service.CreateChallenge(c.Context(), req.did)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"nonce":      nonce,
		"expires_in": 120,
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req pb.AuthLogin
	if err := proto.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(400).SendString("Invalid Protobuf")
	}

	token, err := h.service.VerifyAndLogin(c.Context(), req.Did, req.Nonce, req.Signature)
	if err != nil {
		return c.Status(401).SendString("Authentication failed: " + err.Error())
	}

	resp := &pb.TokenResponse{
		AccessToken: token,
		ExpiresIn:   3600,
	}

	data, _ := proto.Marshal(resp)
	c.Set("Content-Type", "application/x-protobuf")
	return c.Send(data)
}
