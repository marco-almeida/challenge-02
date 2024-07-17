package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/marco-almeida/challenge-02/internal"
	"github.com/marco-almeida/challenge-02/internal/postgresql/db"
)

// VehicleService defines the methods that the vehicle handler will use
type VehicleService interface {
	Create(ctx context.Context, arg db.CreateVehicleParams) (db.Vehicle, error)
	Get(ctx context.Context, id int64) (db.Vehicle, error)
	Delete(ctx context.Context, id int64) error
	GetAll(ctx context.Context, arg db.GetVehiclesParams) ([]db.Vehicle, error)
}

// AssignedOrderService defines the methods that the vehicle handler will use
type AssignedOrderService interface {
	CreateAssignedOrder(ctx context.Context, arg db.CreateAssignedOrderParams) (db.AssignedOrder, error)
	GetVehicleAssignedOrders(ctx context.Context, arg db.GetVehicleAssignedOrdersParams) ([]db.Order, error)
	GetNextOrder(ctx context.Context, id int64) (db.Order, error)
}

// VehicleHandler is the handler for the vehicle service
type VehicleHandler struct {
	vehicleSvc       VehicleService
	assignedOrderSvc AssignedOrderService
}

// NewVehicleHandler returns a new VehicleHandler
func NewVehicleHandler(vehicleSvc VehicleService, assignedOrderSvc AssignedOrderService) *VehicleHandler {
	return &VehicleHandler{
		vehicleSvc:       vehicleSvc,
		assignedOrderSvc: assignedOrderSvc,
	}
}

func (h *VehicleHandler) RegisterRoutes(r *gin.Engine) {
	groupRoutes := r.Group("/api")

	groupRoutes.POST("/v1/vehicles", h.handleCreateVehicle)
	groupRoutes.GET("/v1/vehicles", h.handleGetAllVehicles)
	groupRoutes.GET("/v1/vehicles/:id", h.handleGetVehicle)
	groupRoutes.DELETE("/v1/vehicles/:id", h.handleDeleteVehicle)
	groupRoutes.POST("/v1/vehicles/:id/orders", h.handleAssignOrderToVehicle)
	groupRoutes.GET("/v1/vehicles/:id/unfinished_orders", h.handleGetUnfinishedOrders)
	groupRoutes.GET("/v1/vehicles/:id/next_order", h.handleGetNextOrder)
}

type createVehicleRequest struct {
	NumberPlate       string  `json:"number_plate" binding:"required"`
	MaxWeightCapacity float32 `json:"max_weight_capacity" binding:"required"`
}

func (h *VehicleHandler) handleCreateVehicle(ctx *gin.Context) {
	var req createVehicleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(fmt.Errorf("%w; %w", internal.ErrInvalidParams, err))
		return
	}

	vehicle, err := h.vehicleSvc.Create(ctx, db.CreateVehicleParams{
		NumberPlate:       req.NumberPlate,
		MaxWeightCapacity: req.MaxWeightCapacity,
	})
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, vehicle)
}

type uriIdRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

func (h *VehicleHandler) handleGetVehicle(ctx *gin.Context) {
	var req uriIdRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.Error(fmt.Errorf("%w; %w", internal.ErrInvalidParams, err))
		return
	}

	vehicle, err := h.vehicleSvc.Get(ctx, req.ID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, vehicle)
}

type getAllVehiclesRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (h *VehicleHandler) handleGetAllVehicles(ctx *gin.Context) {
	var req getAllVehiclesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(fmt.Errorf("%w; %w", internal.ErrInvalidParams, err))
		return
	}

	vehicles, err := h.vehicleSvc.GetAll(ctx, db.GetVehiclesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	})
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, vehicles)
}

func (h *VehicleHandler) handleDeleteVehicle(ctx *gin.Context) {
	var req uriIdRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.Error(fmt.Errorf("%w; %w", internal.ErrInvalidParams, err))
		return
	}

	err := h.vehicleSvc.Delete(ctx, req.ID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

type assignOrderToVehicleRequest struct {
	OrderID int64 `json:"order_id" binding:"required"`
}

func (h *VehicleHandler) handleAssignOrderToVehicle(ctx *gin.Context) {
	var vehicleReq uriIdRequest
	if err := ctx.ShouldBindUri(&vehicleReq); err != nil {
		ctx.Error(fmt.Errorf("%w; %w", internal.ErrInvalidParams, err))
		return
	}

	var orderReq assignOrderToVehicleRequest
	if err := ctx.ShouldBindJSON(&orderReq); err != nil {
		ctx.Error(fmt.Errorf("%w; %w", internal.ErrInvalidParams, err))
		return
	}

	assignedOrder, err := h.assignedOrderSvc.CreateAssignedOrder(ctx, db.CreateAssignedOrderParams{
		VehicleID: vehicleReq.ID,
		OrderID:   orderReq.OrderID,
		AssignedAt: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		},
	})

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, assignedOrder)
}

func (h *VehicleHandler) handleGetUnfinishedOrders(ctx *gin.Context) {
	var vehicleReq uriIdRequest
	if err := ctx.ShouldBindUri(&vehicleReq); err != nil {
		ctx.Error(fmt.Errorf("%w; %w", internal.ErrInvalidParams, err))
		return
	}

	var paginationReq getAllVehiclesRequest
	if err := ctx.ShouldBindQuery(&paginationReq); err != nil {
		ctx.Error(fmt.Errorf("%w; %w", internal.ErrInvalidParams, err))
		return
	}

	orders, err := h.assignedOrderSvc.GetVehicleAssignedOrders(ctx, db.GetVehicleAssignedOrdersParams{
		VehicleID: vehicleReq.ID,
		Limit:     paginationReq.PageSize,
		Offset:    (paginationReq.PageID - 1) * paginationReq.PageSize,
	})

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, orders)
}

func (h *VehicleHandler) handleGetNextOrder(ctx *gin.Context) {
	var req uriIdRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.Error(fmt.Errorf("%w; %w", internal.ErrInvalidParams, err))
		return
	}

	order, err := h.assignedOrderSvc.GetNextOrder(ctx, req.ID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, order)
}
