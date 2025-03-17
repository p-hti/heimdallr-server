package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/p-hti/heimdallr-server/internal/domain/model"
)

type AuthService interface {
	SaveUser(
		ctx context.Context,
		email string,
		name string,
		passHash string,
	) (
		err error,
	)
	SignIn(
		ctx context.Context,
		cred model.SignInStruct,
	) (
		string,
		string,
		error,
	)
	Sessions
}

type Sessions interface {
	ParseToken(
		ctx context.Context,
		token string,
	) (
		int64,
		error,
	)
	RefreshToken(
		ctx context.Context,
		refreshToken string,
	) (
		string,
		string,
		error,
	)
	TerminateSession(
		ctx context.Context,
		refreshToken string,
	) error

	TerminateAllSessions(
		ctx context.Context,
		uid int64,
	) error
}

func (h *HTTPServer) SignUp(w http.ResponseWriter, r *http.Request) {
	signUp := model.SingUpStruct{}
	if err := json.NewDecoder(r.Body).Decode(&signUp); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.Service.SaveUser(r.Context(), signUp.Email, signUp.Name, signUp.Password)
	if err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func (h *HTTPServer) SignIn(w http.ResponseWriter, r *http.Request) {
	credentials := model.SignInStruct{}
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := h.Service.SignIn(r.Context(), credentials)
	if err != nil {
		http.Error(w, "smth went wrong", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(map[string]string{
		"access_token": accessToken,
	})
	if err != nil {
		http.Error(w, "failed marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Set-Cookie", fmt.Sprintf("refresh-token=%s; HttpOnly", refreshToken))
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (h *HTTPServer) LogOut(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh-token")
	if err != nil {
		http.Error(w, "cookie is not found", http.StatusBadRequest)
		return
	}

	err = h.Service.TerminateSession(r.Context(), cookie.Value)
	if err != nil {
		http.Error(w, "refreshToken is not found", http.StatusInternalServerError)
		return
	}

	// #TODO: add logic for delete (or blacklist) access token
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh-token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
	})

	w.WriteHeader(http.StatusOK)
}

func (h *HTTPServer) RefreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh-token")
	if err != nil {
		http.Error(w, "cookie is not found", http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := h.Service.RefreshToken(r.Context(), cookie.Value)
	if err != nil {
		http.Error(w, "smth went wrong", http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(map[string]string{
		"access_token": accessToken,
	})
	if err != nil {
		http.Error(w, "failed marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Set-Cookie", fmt.Sprintf("refresh-token=%s; HttpOnly", refreshToken))
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (h *HTTPServer) FullLogOut(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("UserID").(int64)

	err := h.Service.TerminateAllSessions(r.Context(), userId)
	if err != nil {
		http.Error(w, "smth went wrong", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
