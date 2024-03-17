package router

import (
	"fmt"
	"gossip/internal/utils"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
)

func (router *Router) chatRouter() *chi.Mux {
	chatRouter := chi.NewRouter()

	// new client
	chatRouter.Get("/connect", func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			Username string `query:"username"`
			UserId   string `query:"userId"`
		}
		query, err := utils.GetURLQueryStruct[request](r.URL)
		if err != nil {
			log.Println(err.Error())
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		userId, err := uuid.FromString(query.UserId)
		if err != nil {
			log.Println(err.Error())
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		conn, err := router.ChatService.UpgradeConnection(w, r)
		if err != nil {
			log.Println(err.Error())
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		router.ChatService.NewClient(userId, query.Username, conn)
	})

	// new room
	chatRouter.Post("/rooms", func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			Name string `json:"name"`
		}
		body, err := utils.ReadJSON[request](r)
		if err != nil {
			log.Println(err.Error())
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		router.ChatService.NewRoom(body.Name)

		utils.WriteJSON(w, http.StatusCreated, utils.BaseResponse{
			Error:   false,
			Message: fmt.Sprintf("room created %s", body.Name),
		})
	})

	// destroy room
	chatRouter.Delete(
		"/rooms/{name}",
		func(w http.ResponseWriter, r *http.Request) {
			name := chi.URLParam(r, "name")

			if err := router.ChatService.DestroyRoom(name); err != nil {
				utils.WriteError(w, http.StatusBadRequest, err)
				return
			}

			utils.WriteJSON(w, http.StatusCreated, utils.BaseResponse{
				Error:   false,
				Message: fmt.Sprintf("room destroyed %s", name),
			})
		},
	)

	// client join room
	chatRouter.Post(
		"/rooms/{name}",
		func(w http.ResponseWriter, r *http.Request) {
			roomName := chi.URLParam(r, "name")

			type request struct {
				UserId string `json:"userId"`
			}
			body, err := utils.ReadJSON[request](r)
			if err != nil {
				log.Println(err.Error())
				utils.WriteError(w, http.StatusBadRequest, err)
				return
			}

			userId, err := uuid.FromString(body.UserId)
			if err != nil {
				log.Println(err.Error())
				utils.WriteError(w, http.StatusBadRequest, err)
				return
			}

			if err := router.ChatService.ClientJoinRoom(userId, roomName); err != nil {
				utils.WriteError(w, http.StatusBadRequest, err)
				return
			}

			utils.WriteJSON(w, http.StatusCreated, utils.BaseResponse{
				Error:   false,
				Message: fmt.Sprintf("joined room %s", roomName),
			})
		},
	)

	return chatRouter
}
