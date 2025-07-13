package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/kitesi/relaytalk/api/routes/auth"
	"github.com/kitesi/relaytalk/api/routes/channels"
	"github.com/kitesi/relaytalk/api/routes/messages"
	"github.com/kitesi/relaytalk/api/routes/servers"
	"github.com/kitesi/relaytalk/db"
)

func RegisterRoutes(store *db.Queries, r chi.Router) {
	// Public
	r.Post("/register", auth.Register(store))
	r.Post("/login", auth.Login(store))

	// Protected
	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware(store))

		// Server routes
		r.Post("/servers", servers.CreateServer(store))
		// r.Delete("/servers/{serverId}", servers.DeleteServerHandler(store))

		// Channel routes
		r.Post("/servers/{serverId}/channels", channels.CreateChannel(store))
		// r.Get("/servers/{serverId}/channels", channels.ListChannelsHandler(store))

		r.Post("/channels/{channelId}/messages", messages.SendMessage(store))
	})
}
