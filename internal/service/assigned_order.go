package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"

	"github.com/marco-almeida/challenge-02/internal"
	"github.com/marco-almeida/challenge-02/internal/postgresql/db"
)

// AssignedOrderRepository defines the methods that any assigned order repository should implement.
type AssignedOrderRepository interface {
	GetVehicleAssignedOrders(ctx context.Context, arg db.GetVehicleAssignedOrdersParams) ([]db.Order, error)
	CreateAssignedOrder(ctx context.Context, arg db.CreateAssignedOrderParams) (db.AssignedOrder, error)
}

type AssignedOrderService struct {
	repo        AssignedOrderRepository
	vehicleRepo VehicleRepository
	orderRepo   OrderRepository
}

// NewAssignedOrderService returns a new AssignedOrderService.
func NewAssignedOrderService(repo AssignedOrderRepository, vehicleRepo VehicleRepository, orderRepo OrderRepository) *AssignedOrderService {
	return &AssignedOrderService{
		repo:        repo,
		vehicleRepo: vehicleRepo,
		orderRepo:   orderRepo,
	}
}

func (s *AssignedOrderService) GetVehicleAssignedOrders(ctx context.Context, arg db.GetVehicleAssignedOrdersParams) ([]db.Order, error) {
	// check if vehicle exists
	_, err := s.vehicleRepo.Get(ctx, arg.VehicleID)

	if err != nil {
		return nil, fmt.Errorf("%w: %s", internal.ErrVehicleNotFound, err.Error())
	}

	// no rows must not be an error
	orders, err := s.repo.GetVehicleAssignedOrders(ctx, arg)
	if err != nil {
		return nil, err
	}

	// sort orders by closest
	sortOrdersByClosest(orders, initialLat, initialLong)

	return orders, nil
}

func (s *AssignedOrderService) CreateAssignedOrder(ctx context.Context, arg db.CreateAssignedOrderParams) (db.AssignedOrder, error) {
	// check if vehicle exists
	vehicle, err := s.vehicleRepo.Get(ctx, arg.VehicleID)
	if err != nil {
		if errors.Is(err, internal.ErrNoRows) {
			return db.AssignedOrder{}, fmt.Errorf("%w: %s", internal.ErrVehicleNotFound, err.Error())
		}
		return db.AssignedOrder{}, err
	}

	// check if order exists
	order, err := s.orderRepo.Get(ctx, arg.OrderID)
	if err != nil {
		if errors.Is(err, internal.ErrNoRows) {
			return db.AssignedOrder{}, fmt.Errorf("%w: %s", internal.ErrOrderNotFound, err.Error())
		}
		return db.AssignedOrder{}, err
	}

	// order mustnt be flagged as finished
	if order.Finished {
		return db.AssignedOrder{}, internal.ErrOrderAlreadyFinished
	}

	// check if vehicle can handle the order
	if vehicle.CurrentWeight+order.Weight > vehicle.MaxWeightCapacity {
		return db.AssignedOrder{}, internal.ErrOrderTooHeavy
	}

	// update vehicle's current weight
	newWeight := vehicle.CurrentWeight + order.Weight
	_, err = s.vehicleRepo.UpdateCurrentWeight(ctx, db.UpdateVehicleCurrentWeightParams{
		CurrentWeight: newWeight,
		ID:            vehicle.ID,
	})

	if err != nil {
		return db.AssignedOrder{}, err
	}

	return s.repo.CreateAssignedOrder(ctx, arg)
}

func (s *AssignedOrderService) GetNextOrder(ctx context.Context, id int64) (db.Order, error) {
	// logic already implemented in GetVehicleAssignedOrders, get the first (closest) order
	orders, err := s.GetVehicleAssignedOrders(ctx, db.GetVehicleAssignedOrdersParams{
		VehicleID: id,
		Limit:     math.MaxInt32,
		Offset:    0,
	})

	fmt.Println(orders)
	if err != nil {
		return db.Order{}, err
	}

	if len(orders) == 0 {
		return db.Order{}, nil
	}

	return orders[0], nil
}

const (
	earthRadiusKm = 6371     // radius of the earth in kilometers.
	initialLat    = 38.71814 // Fintech House in Lisbon
	initialLong   = -9.14552 // Fintech House in Lisbon
)

// sortOrdersByClosest sorts orders by closest to the given coordinates
func sortOrdersByClosest(orders []db.Order, lat float64, long float64) {
	// we start at 38.71814, -9.14552
	sort.Slice(orders, func(i, j int) bool {
		distanceFromI := DistanceBetweenCoordinates(lat, long, orders[i].Destination.P.X, orders[i].Destination.P.Y)
		distanceFromJ := DistanceBetweenCoordinates(lat, long, orders[j].Destination.P.X, orders[j].Destination.P.Y)

		return distanceFromI < distanceFromJ
	})
}

// degreesToRadians converts from degrees to radians.
func degreesToRadians(d float64) float64 {
	return d * math.Pi / 180
}

// DistanceBetweenCoordinates calculates the shortest path between two coordinates on the surface
// of the Earth. This function returns two units of measure, the first is the
// distance in miles, the second is the distance in kilometers.
func DistanceBetweenCoordinates(latitude1, longitude1, latitude2, longitude2 float64) float64 {
	lat1 := degreesToRadians(latitude1)
	lon1 := degreesToRadians(longitude1)
	lat2 := degreesToRadians(latitude2)
	lon2 := degreesToRadians(longitude2)

	diffLat := lat2 - lat1
	diffLon := lon2 - lon1

	a := math.Pow(math.Sin(diffLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*
		math.Pow(math.Sin(diffLon/2), 2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	km := c * earthRadiusKm

	return km
}
