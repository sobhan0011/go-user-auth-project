package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"

	authuc "dekamond/internal/usecase/auth"
)

type AuthHandler struct {
	authUsecase *authuc.AuthUsecase
}

func NewAuthHandler(authUsecase *authuc.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

func (h *AuthHandler) RequestOTP(w http.ResponseWriter, req *http.Request) {
	var body struct{ Phone string `json:"phone"` }
	if !h.decodeAndValidate(w, req, &body) || !isValidE164(body.Phone) {
		WriteJSON(w, http.StatusBadRequest, ApiResponse{Error: "invalid_phone"})
		return
	}
	code, err := h.authUsecase.RequestOTP(req.Context(), body.Phone)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ApiResponse{Error: err.Error()})
		return
	}
	log.Printf("OTP for %s: %s", body.Phone, code)
	WriteJSON(w, http.StatusOK, ApiResponse{Message: "otp_sent"})
}

func (h *AuthHandler) VerifyOTP(w http.ResponseWriter, req *http.Request) {
	var body struct{ Phone, Code string }
	if !h.decodeAndValidate(w, req, &body) || !isValidE164(body.Phone) || strings.TrimSpace(body.Code) == "" {
		WriteJSON(w, http.StatusBadRequest, ApiResponse{Error: "invalid_payload"})
		return
	}
	token, user, err := h.authUsecase.VerifyOTPAndIssueToken(req.Context(), body.Phone, body.Code)
	if err != nil {
		WriteJSON(w, http.StatusUnauthorized, ApiResponse{Error: err.Error()})
		return
	}
	WriteJSON(w, http.StatusOK, ApiResponse{Data: map[string]any{"token": token, "user": user}})
}

func (h *AuthHandler) decodeAndValidate(_ http.ResponseWriter, req *http.Request, body any) bool {
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		return false
	}
	return true
}

var e164Re = regexp.MustCompile(`^\+[1-9]\d{7,14}$`)

func isValidE164(phone string) bool {
	phone = strings.TrimSpace(phone)
	return e164Re.MatchString(phone)
}