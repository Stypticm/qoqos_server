package handlers

import (
	"api-service/internal/models"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {
	skupkaRouteSetup := func(r chi.Router) {
		h := NewBaseHandler[models.Skupka]()
		r.Get("/", h.List)
		r.Get("/{id}", h.GetByID)
		r.Post("/", h.Create)
		r.Put("/{id}", h.Update)
		r.Patch("/{id}", h.Patch)
		r.Delete("/{id}", h.Delete)
	}

	r.Route("/skupka", skupkaRouteSetup)
	r.Route("/skupkas", skupkaRouteSetup) // alias for frontend compatibility

	r.Route("/masters", func(r chi.Router) {
		h := NewBaseHandler[models.Master]()
		r.Get("/", h.List)
		r.Get("/{id}", h.GetByID)
		r.Group(func(r chi.Router) {
			r.Use(AdminMiddleware)
			r.Post("/", h.Create)
			r.Put("/{id}", h.Update)
			r.Patch("/{id}", h.Patch)
			r.Delete("/{id}", h.Delete)
		})
	})

	r.Route("/points", func(r chi.Router) {
		h := NewBaseHandler[models.Point]()
		r.Get("/", h.List)
		r.Get("/{id}", h.GetByID)
		r.Group(func(r chi.Router) {
			r.Use(AdminMiddleware)
			r.Post("/", h.Create)
			r.Put("/{id}", h.Update)
			r.Patch("/{id}", h.Patch)
			r.Delete("/{id}", h.Delete)
		})
	})

	r.Route("/device-inspections", func(r chi.Router) {
		h := NewBaseHandler[models.DeviceInspection]()
		r.Get("/", h.List)
		r.Get("/{id}", h.GetByID)
		r.Post("/", h.Create)
		r.Put("/{id}", h.Update)
		r.Patch("/{id}", h.Patch)
		r.Delete("/{id}", h.Delete)
	})

	r.Route("/devices", func(r chi.Router) {
		h := NewBaseHandler[models.Device]()
		r.Get("/", h.List)
		r.Get("/{id}", h.GetByID)
		r.Post("/", h.Create)
		r.Put("/{id}", h.Update)
		r.Patch("/{id}", h.Patch)
		r.Delete("/{id}", h.Delete)
	})

	r.Route("/market-prices", func(r chi.Router) {
		h := NewBaseHandler[models.MarketPrice]()
		r.Get("/", h.List)
		r.Get("/{id}", h.GetByID)
		r.Post("/", h.Create)
		r.Put("/{id}", h.Update)
		r.Patch("/{id}", h.Patch)
		r.Delete("/{id}", h.Delete)
	})

	r.Route("/marketplace-lots", func(r chi.Router) {
		h := NewBaseHandler[models.MarketplaceLot]()
		r.Get("/", h.List)
		r.Get("/{id}", h.GetByID)

		r.Group(func(r chi.Router) {
			r.Use(AdminMiddleware)
			r.Post("/", h.Create)
			r.Put("/{id}", h.Update)
			r.Patch("/{id}", h.Patch)
			r.Delete("/{id}", h.Delete)
		})
	})

	r.Route("/cart-items", func(r chi.Router) {
		h := NewBaseHandler[models.CartItem]()
		r.Get("/", h.List)
		r.Get("/{id}", h.GetByID)
		r.Post("/", h.Create)
		r.Put("/{id}", h.Update)
		r.Patch("/{id}", h.Patch)
		r.Delete("/{id}", h.Delete)
	})

	r.Route("/favorite-items", func(r chi.Router) {
		h := NewBaseHandler[models.FavoriteItem]()
		r.Get("/", h.List)
		r.Get("/{id}", h.GetByID)
		r.Post("/", h.Create)
		r.Put("/{id}", h.Update)
		r.Patch("/{id}", h.Patch)
		r.Delete("/{id}", h.Delete)
	})

	r.Route("/orders", func(r chi.Router) {
		h := NewBaseHandler[models.Order]()
		r.Get("/", h.List)
		r.Get("/{id}", h.GetByID)
		r.Post("/", h.Create)
		r.Put("/{id}", h.Update)
		r.Patch("/{id}", h.Patch)
		r.Delete("/{id}", h.Delete)
	})

	r.Route("/order-items", func(r chi.Router) {
		h := NewBaseHandler[models.OrderItem]()
		r.Get("/", h.List)
		r.Get("/{id}", h.GetByID)
		r.Post("/", h.Create)
		r.Put("/{id}", h.Update)
		r.Patch("/{id}", h.Patch)
		r.Delete("/{id}", h.Delete)
	})

	r.Route("/operator-chats", func(r chi.Router) {
		h := NewBaseHandler[models.OperatorChat]()
		r.Get("/", h.List)
		r.Get("/{id}", h.GetByID)
		r.Post("/", h.Create)
		r.Put("/{id}", h.Update)
		r.Patch("/{id}", h.Patch)
		r.Delete("/{id}", h.Delete)
	})

	r.Route("/operator-messages", func(r chi.Router) {
		h := NewBaseHandler[models.OperatorMessage]()
		r.Get("/", h.List)
		r.Get("/{id}", h.GetByID)
		r.Post("/", h.Create)
		r.Put("/{id}", h.Update)
		r.Patch("/{id}", h.Patch)
		r.Delete("/{id}", h.Delete)
	})

	r.Route("/auth-requests", func(r chi.Router) {
		h := NewBaseHandler[models.AuthRequest]()
		r.Get("/", h.List)
		r.Get("/{id}", h.GetByID)
		r.Post("/", h.Create)
		r.Put("/{id}", h.Update)
		r.Patch("/{id}", h.Patch)
		r.Delete("/{id}", h.Delete)
	})

	r.Route("/users", func(r chi.Router) {
		h := NewBaseHandler[models.User]()
		r.Group(func(r chi.Router) {
			r.Use(AdminMiddleware)
			r.Get("/", h.List)
			r.Get("/{id}", h.GetByID)
			r.Post("/", h.Create)
			r.Put("/{id}", h.Update)
			r.Patch("/{id}", h.Patch)
			r.Delete("/{id}", h.Delete)
		})
	})

	r.Route("/quick-leads", func(r chi.Router) {
		h := NewBaseHandler[models.QuickLead]()
		r.Get("/", h.List)
		r.Get("/{id}", h.GetByID)
		r.Post("/", h.Create)
		r.Put("/{id}", h.Update)
		r.Patch("/{id}", h.Patch)
		r.Delete("/{id}", h.Delete)
	})

	r.Route("/blog-posts", func(r chi.Router) {
		h := NewBaseHandler[models.BlogPost]()
		r.Get("/", h.List)
		r.Get("/{id}", h.GetByID)
		r.Group(func(r chi.Router) {
			r.Use(AdminMiddleware)
			r.Post("/", h.Create)
			r.Put("/{id}", h.Update)
			r.Patch("/{id}", h.Patch)
			r.Delete("/{id}", h.Delete)
		})
	})

	r.Route("/trade-in-evaluations", func(r chi.Router) {
		h := NewBaseHandler[models.TradeInEvaluation]()
		r.Get("/", h.List)
		r.Get("/{id}", h.GetByID)
		r.Post("/", h.Create)
		r.Put("/{id}", h.Update)
		r.Patch("/{id}", h.Patch)
		r.Delete("/{id}", h.Delete)
	})

	r.Route("/push-subscriptions", func(r chi.Router) {
		h := NewBaseHandler[models.PushSubscription]()
		r.Get("/", h.List)
		r.Get("/{id}", h.GetByID)
		r.Post("/", h.Create)
		r.Put("/{id}", h.Update)
		r.Patch("/{id}", h.Patch)
		r.Delete("/{id}", h.Delete)
	})

	r.Route("/bot-accesses", func(r chi.Router) {
		h := NewBaseHandler[models.BotAccess]()
		r.Get("/", h.List)
		r.Get("/{id}", h.GetByID)
		r.Post("/", h.Create)
		r.Put("/{id}", h.Update)
		r.Patch("/{id}", h.Patch)
		r.Delete("/{id}", h.Delete)
	})

	r.Route("/repair-requests", func(r chi.Router) {
		h := NewBaseHandler[models.RepairRequest]()
		r.Get("/", h.List)
		r.Get("/{id}", h.GetByID)
		r.Post("/", h.Create)
		r.Put("/{id}", h.Update)
		r.Patch("/{id}", h.Patch)
		r.Delete("/{id}", h.Delete)
	})

	r.Route("/agent-audit-logs", func(r chi.Router) {
		h := NewBaseHandler[models.AgentAuditLog]()
		r.Get("/", h.List)
		r.Get("/{id}", h.GetByID)
		r.Post("/", h.Create)
		r.Put("/{id}", h.Update)
		r.Patch("/{id}", h.Patch)
		r.Delete("/{id}", h.Delete)
	})

	r.Route("/idempotency-keys", func(r chi.Router) {
		h := NewBaseHandler[models.IdempotencyKey]()
		r.Get("/", h.List)
		r.Get("/{id}", h.GetByID)
		r.Post("/", h.Create)
		r.Put("/{id}", h.Update)
		r.Patch("/{id}", h.Patch)
		r.Delete("/{id}", h.Delete)
		r.Get("/distinct/{field}", h.Distinct)
	})

	SetupStockRoutes(r)

	// Auth endpoints
	r.Post("/users/login", HandleLogin)
	r.Post("/users/register", HandleRegister)
}
