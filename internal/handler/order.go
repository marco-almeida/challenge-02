package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/marco-almeida/challenge-02/internal"
	"github.com/marco-almeida/challenge-02/internal/postgresql/db"
)

// OrderService defines the methods that the order handler will use
type OrderService interface {
	Create(ctx context.Context, arg db.CreateOrderParams) (db.Order, error)
	UpdateObservations(ctx context.Context, arg db.UpdateOrderObservationsParams) (db.Order, error)
	Get(ctx context.Context, id int64) (db.Order, error)
}

// OrderHandler is the handler for the order service
type OrderHandler struct {
	orderSvc OrderService
}

// NewOrderHandler returns a new OrderHandler
func NewOrderHandler(orderSvc OrderService) *OrderHandler {
	return &OrderHandler{
		orderSvc: orderSvc,
	}
}

func (h *OrderHandler) RegisterRoutes(r *gin.Engine) {
	groupRoutes := r.Group("/api")

	groupRoutes.POST("/v1/orders", h.handleCreateOrder)
	groupRoutes.GET("/v1/orders/:id", h.handleGetOrder)
	groupRoutes.PATCH("/v1/orders/:id/observations", h.handleOverrideObservations)
}

type createOrderRequest struct {
	Weight       float32 `json:"weight" binding:"required"`
	DestinationX float64 `json:"destination_x" binding:"required"`
	DestinationY float64 `json:"destination_y" binding:"required"`
	Observations string  `json:"observations"`
}

func (h *OrderHandler) handleCreateOrder(ctx *gin.Context) {
	var req createOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(fmt.Errorf("%w; %w", internal.ErrInvalidParams, err))
		return
	}

	order, err := h.orderSvc.Create(ctx, db.CreateOrderParams{
		Weight: req.Weight,
		Destination: pgtype.Point{
			P:     pgtype.Vec2{X: req.DestinationX, Y: req.DestinationY},
			Valid: true,
		},
		Observations: req.Observations,
		Finished:     false,
	})

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, order)
}

type updateObservationRequest struct {
	Observations string `json:"observations" binding:"required"`
}

// will override the observation field
func (h *OrderHandler) handleOverrideObservations(ctx *gin.Context) {
	var req updateObservationRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(fmt.Errorf("%w; %w", internal.ErrInvalidParams, err))
		return
	}

	orderIdStr := ctx.Param("id")
	if orderIdStr == "" {
		ctx.Error(internal.ErrInvalidParams)
		return
	}

	orderId, err := strconv.ParseInt(orderIdStr, 10, 64)
	if err != nil {
		ctx.Error(internal.ErrInvalidParams)
		return
	}

	order, err := h.orderSvc.UpdateObservations(ctx, db.UpdateOrderObservationsParams{
		ID:           orderId,
		Observations: req.Observations,
	})
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, order)
}

func (h *OrderHandler) handleGetOrder(ctx *gin.Context) {
	orderIdStr := ctx.Param("id")
	if orderIdStr == "" {
		ctx.Error(internal.ErrInvalidParams)
		return
	}

	orderId, err := strconv.ParseInt(orderIdStr, 10, 64)
	if err != nil {
		ctx.Error(internal.ErrInvalidParams)
		return
	}

	order, err := h.orderSvc.Get(ctx, orderId)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, order)
}
