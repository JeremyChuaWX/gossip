package router

import (
	"errors"
	"gossip/internal/constants"
	_roomuser "gossip/internal/domains/roomuser"
	_user "gossip/internal/domains/user"
	"gossip/internal/password"
	"gossip/internal/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
)

var userForbiddenError = errors.New("mismatch user ID")

func (router *Router) userRouter() *chi.Mux {
	userRouter := chi.NewMux()
	userRouter.Use(router.Middlewares.AuthMiddleware())

	// find one user
	userRouter.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		urlId := chi.URLParam(r, "id")
		id, err := uuid.FromString(urlId)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		dto := _user.FindOneDTO{
			Id: id,
		}

		user, err := router.UserRepository.FindOne(r.Context(), dto)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		type response struct {
			utils.BaseResponse
			User _user.User `json:"user"`
		}
		utils.WriteJSON(w, http.StatusOK, response{
			BaseResponse: utils.BaseResponse{
				Error:   false,
				Message: "found user",
			},
			User: user,
		})
	})

	// find one user by username
	userRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			Username string `query:"username"`
		}
		req, err := utils.GetURLQueryStruct[request](r.URL)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		dto := _user.FindOneByUsernameDTO{
			Username: req.Username,
		}

		user, err := router.UserRepository.FindOneByUsername(r.Context(), dto)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		type response struct {
			utils.BaseResponse
			User _user.User `json:"user"`
		}
		utils.WriteJSON(w, http.StatusOK, response{
			BaseResponse: utils.BaseResponse{
				Error:   false,
				Message: "found user",
			},
			User: user,
		})
	})

	// update user
	userRouter.Patch("/{id}", func(w http.ResponseWriter, r *http.Request) {
		urlId := chi.URLParam(r, "id")
		id, err := uuid.FromString(urlId)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		authUserId := r.Context().Value(constants.USER_ID_CONTEXT_KEY).(uuid.UUID)
		if authUserId != id {
			utils.WriteError(w, http.StatusForbidden, userForbiddenError)
			return
		}

		type Request struct {
			Username *string `json:"username,omitempty"`
			Password *string `json:"password,omitempty"`
		}
		req, err := utils.ReadJSON[Request](r)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		dto := _user.UpdateDTO{
			Id:           id,
			Username:     req.Username,
			PasswordHash: nil,
		}

		if req.Password != nil {
			passwordHash, err := password.Hash(*req.Password)
			dto.PasswordHash = &passwordHash
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, err)
				return
			}
		}

		user, err := router.UserRepository.Update(r.Context(), dto)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		type response struct {
			utils.BaseResponse
			User _user.User `json:"user"`
		}
		utils.WriteJSON(w, http.StatusOK, response{
			BaseResponse: utils.BaseResponse{
				Error:   false,
				Message: "updated user",
			},
			User: user,
		})
	})

	// delete user
	userRouter.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
		urlId := chi.URLParam(r, "id")
		id, err := uuid.FromString(urlId)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		authUserId := r.Context().Value(constants.USER_ID_CONTEXT_KEY).(uuid.UUID)
		if authUserId != id {
			utils.WriteError(w, http.StatusForbidden, userForbiddenError)
			return
		}

		dto := _user.DeleteDTO{
			Id: id,
		}
		user, err := router.UserRepository.Delete(r.Context(), dto)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		type response struct {
			utils.BaseResponse
			User _user.User `json:"user"`
		}
		utils.WriteJSON(w, http.StatusOK, response{
			BaseResponse: utils.BaseResponse{
				Error:   false,
				Message: "deleted user",
			},
			User: user,
		})
	})

	// join room
	userRouter.Post(
		"/{id}/join-room",
		func(w http.ResponseWriter, r *http.Request) {
			urlId := chi.URLParam(r, "id")
			id, err := uuid.FromString(urlId)
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, err)
				return
			}

			authUserId := r.Context().Value(constants.USER_ID_CONTEXT_KEY).(uuid.UUID)
			if authUserId != id {
				utils.WriteError(w, http.StatusForbidden, userForbiddenError)
				return
			}

			type Request struct {
				RoomId uuid.UUID `json:"roomId"`
			}
			req, err := utils.ReadJSON[Request](r)
			if err != nil {
				utils.WriteError(w, http.StatusBadRequest, err)
				return
			}

			dto := _roomuser.CreateDTO{
				RoomId: req.RoomId,
				UserId: id,
			}
			roomUser, err := router.RoomUserRepository.Create(r.Context(), dto)
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, err)
				return
			}

			type response struct {
				utils.BaseResponse
				RoomUser _roomuser.RoomUser `json:"roomUser"`
			}
			utils.WriteJSON(w, http.StatusOK, response{
				BaseResponse: utils.BaseResponse{
					Error:   false,
					Message: "join room",
				},
				RoomUser: roomUser,
			})
		},
	)

	// leave room
	userRouter.Post(
		"/{id}/leave-room",
		func(w http.ResponseWriter, r *http.Request) {
			urlId := chi.URLParam(r, "id")
			id, err := uuid.FromString(urlId)
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, err)
				return
			}

			authUserId := r.Context().Value(constants.USER_ID_CONTEXT_KEY).(uuid.UUID)
			if authUserId != id {
				utils.WriteError(w, http.StatusForbidden, userForbiddenError)
				return
			}

			type Request struct {
				RoomId uuid.UUID `json:"roomId"`
			}
			req, err := utils.ReadJSON[Request](r)
			if err != nil {
				utils.WriteError(w, http.StatusBadRequest, err)
				return
			}

			dto := _roomuser.DeleteDTO{
				RoomId: req.RoomId,
				UserId: id,
			}
			roomUser, err := router.RoomUserRepository.Delete(r.Context(), dto)
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, err)
				return
			}

			type response struct {
				utils.BaseResponse
				RoomUser _roomuser.RoomUser `json:"roomUser"`
			}
			utils.WriteJSON(w, http.StatusOK, response{
				BaseResponse: utils.BaseResponse{
					Error:   false,
					Message: "leave room",
				},
				RoomUser: roomUser,
			})
		},
	)

	return userRouter
}
