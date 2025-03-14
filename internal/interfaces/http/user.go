package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Financial-Partner/server/internal/contextutil"
	"github.com/Financial-Partner/server/internal/entities"
	"github.com/Financial-Partner/server/internal/interfaces/http/dto"
	httperror "github.com/Financial-Partner/server/internal/interfaces/http/error"
	responde "github.com/Financial-Partner/server/internal/interfaces/http/respond"
)

//go:generate mockgen -source=user.go -destination=mocks/user_mock.go -package=mocks

type UserService interface {
	GetUser(ctx context.Context, email string) (*entities.User, error)
	GetOrCreateUser(ctx context.Context, email, name string) (*entities.User, error)
	UpdateUserName(ctx context.Context, email, name string) (*entities.User, error)
}

// CreateUser CreateUser
// @Summary CreateUser
// @Description Create a new user when first login
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer {token}" default "Bearer "
// @Param request body dto.CreateUserRequest true "Create user request"
// @Success 200 {object} dto.CreateUserResponse "Create user successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request format"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /users [post]
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.WithError(err).Warnf("Invalid request format")
		responde.WithError(w, r, h.log, err, httperror.ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	ctxEmail := r.Context().Value(contextutil.UserEmailKey)
	if ctxEmail == nil || ctxEmail.(string) != req.Email {
		h.log.Warnf("Email mismatch or not found in context")
		responde.WithError(w, r, h.log, nil, httperror.ErrEmailMismatch, http.StatusUnauthorized)
		return
	}

	userEntity, err := h.userService.GetOrCreateUser(r.Context(), req.Email, req.Name)
	if err != nil {
		h.log.WithError(err).Errorf("Failed to create user")
		responde.WithError(w, r, h.log, err, httperror.ErrFailedToCreateUser, http.StatusInternalServerError)
		return
	}

	response := dto.CreateUserResponse{
		ID:        userEntity.ID.Hex(),
		Email:     userEntity.Email,
		Name:      userEntity.Name,
		Diamonds:  userEntity.Wallet.Diamonds,
		Savings:   userEntity.Wallet.Savings,
		CreatedAt: userEntity.CreatedAt.Format(time.RFC3339),
	}

	responde.WithJSON(w, r, response, http.StatusOK)
}

// UpdateUser UpdateUser
// @Summary UpdateUser
// @Description Update the current user's nickname
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer {token}" default "Bearer "
// @Param request body dto.UpdateUserRequest true "Update user request"
// @Success 200 {object} dto.UpdateUserResponse "Update user successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request format"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /users/me [put]
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.WithError(err).Warnf("Invalid request format")
		responde.WithError(w, r, h.log, err, httperror.ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	email := r.Context().Value(contextutil.UserEmailKey)
	if email == nil {
		h.log.Errorf("User email not found in context")
		responde.WithError(w, r, h.log, nil, httperror.ErrEmailNotFound, http.StatusInternalServerError)
		return
	}

	updatedUser, err := h.userService.UpdateUserName(r.Context(), email.(string), req.Name)
	if err != nil {
		h.log.WithError(err).Errorf("Failed to update user")
		responde.WithError(w, r, h.log, err, httperror.ErrFailedToUpdateUser, http.StatusInternalServerError)
		return
	}

	response := dto.UpdateUserResponse{
		ID:        updatedUser.ID.Hex(),
		Email:     updatedUser.Email,
		Name:      updatedUser.Name,
		Diamonds:  updatedUser.Wallet.Diamonds,
		Savings:   updatedUser.Wallet.Savings,
		UpdatedAt: updatedUser.UpdatedAt.Format(time.RFC3339),
	}

	responde.WithJSON(w, r, response, http.StatusOK)
}

// GetUser GetUser
// @Summary GetUser
// @Description Get the detailed information of the current user, with the option to return specific data fields
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer {token}" default "Bearer "
// @Param scope query []string false "Fields to include (profile, wallet, character). If not specified, returns all" collectionFormat(multi)
// @Success 200 {object} dto.GetUserResponse "Successfully retrieved user information"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /users/me [get]
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value(contextutil.UserEmailKey)
	if email == nil {
		h.log.Errorf("User email not found in context")
		responde.WithError(w, r, h.log, nil, httperror.ErrEmailNotFound, http.StatusInternalServerError)
		return
	}

	scopes := r.URL.Query()["scope"]
	logger := h.log.WithFields(map[string]any{
		"email":  email,
		"scopes": scopes,
	})

	userEntity, err := h.userService.GetUser(r.Context(), email.(string))
	if err != nil {
		logger.WithError(err).Errorf("Failed to get user")
		responde.WithError(w, r, h.log, err, httperror.ErrFailedToGetUser, http.StatusInternalServerError)
		return
	}

	response := buildUserResponse(userEntity, scopes)

	responde.WithJSON(w, r, response, http.StatusOK)
}

func buildUserResponse(user *entities.User, scopes []string) dto.GetUserResponse {
	response := dto.GetUserResponse{
		ID:        user.ID.Hex(),
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}

	if len(scopes) == 0 {
		scopes = []string{"profile", "wallet", "character"}
	}

	for _, scope := range scopes {
		switch scope {
		case "profile":
			response.Email = &user.Email
			response.Name = &user.Name
		case "wallet":
			response.Wallet = &dto.WalletResponse{
				Diamonds: user.Wallet.Diamonds,
				Savings:  user.Wallet.Savings,
			}
		case "character":
			response.Character = &dto.CharacterResponse{
				ID:       user.Character.ID,
				Name:     user.Character.Name,
				ImageURL: user.Character.ImageURL,
			}
		}
	}

	return response
}
