package vehicles

import (
	"fmt"
	"testing"
	"time"

	"github.com/devmcclu/the-block/backend/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:?_foreign_keys=on"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&database.Vehicle{}, &database.DamageNote{}, &database.VehicleImage{}, &database.Bid{})
	require.NoError(t, err)

	return db
}

func newService(db *gorm.DB) RealVehiclesService {
	return RealVehiclesService{
		DB:                      db,
		MaxAuctionDurationHours: 720,
		MinBidIncrement:         100,
	}
}

func activeAuctionStart() string {
	return time.Now().UTC().Add(-1 * time.Hour).Format("2006-01-02T15:04:05")
}

func expiredAuctionStart() string {
	return time.Now().UTC().Add(-721 * time.Hour).Format("2006-01-02T15:04:05")
}

func seedVehicle(t *testing.T, db *gorm.DB, opts ...func(*database.Vehicle)) database.Vehicle {
	t.Helper()
	v := database.Vehicle{
		ExternalID:   "test-vehicle-1",
		VIN:          "1HGBH41JXMN109186",
		Year:         2021,
		Make:         "Honda",
		VehicleModel: "Civic",
		Trim:         "EX",
		AuctionStart: activeAuctionStart(),
		StartingBid:  1000,
		CurrentBid:   0,
		BidCount:     0,
	}
	for _, opt := range opts {
		opt(&v)
	}
	require.NoError(t, db.Create(&v).Error)
	return v
}

func TestPlaceBid_FirstBidAtStartingBid(t *testing.T) {
	db := setupTestDB(t)
	svc := newService(db)

	seedVehicle(t, db)

	updated, err := svc.UpdateVehicle("test-vehicle-1", database.VehicleUpdate{BidAmount: new(1000)})
	require.NoError(t, err)
	assert.Equal(t, 1000, updated.CurrentBid)
	assert.Equal(t, 1, updated.BidCount)
}

func TestPlaceBid_FirstBidBelowStartingBid(t *testing.T) {
	db := setupTestDB(t)
	svc := newService(db)

	seedVehicle(t, db)

	_, err := svc.UpdateVehicle("test-vehicle-1", database.VehicleUpdate{BidAmount: new(999)})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrBidTooLow)
}

func TestPlaceBid_SecondBidMeetsIncrement(t *testing.T) {
	db := setupTestDB(t)
	svc := newService(db)

	seedVehicle(t, db)

	_, err := svc.UpdateVehicle("test-vehicle-1", database.VehicleUpdate{BidAmount: new(1000)})
	require.NoError(t, err)

	// Minimum second bid is 1000 + 100 = 1100
	updated, err := svc.UpdateVehicle("test-vehicle-1", database.VehicleUpdate{BidAmount: new(1100)})
	require.NoError(t, err)
	assert.Equal(t, 1100, updated.CurrentBid)
	assert.Equal(t, 2, updated.BidCount)
}

func TestPlaceBid_SecondBidBelowIncrement(t *testing.T) {
	db := setupTestDB(t)
	svc := newService(db)

	seedVehicle(t, db)

	_, err := svc.UpdateVehicle("test-vehicle-1", database.VehicleUpdate{BidAmount: new(1000)})
	require.NoError(t, err)

	// 1099 < 1000 + 100 = 1100
	_, err = svc.UpdateVehicle("test-vehicle-1", database.VehicleUpdate{BidAmount: new(1099)})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrBidTooLow)
}

func TestPlaceBid_AuctionEnded(t *testing.T) {
	db := setupTestDB(t)
	svc := newService(db)

	seedVehicle(t, db, func(v *database.Vehicle) {
		v.AuctionStart = expiredAuctionStart()
	})

	_, err := svc.UpdateVehicle("test-vehicle-1", database.VehicleUpdate{BidAmount: new(1000)})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrAuctionEnded)
}

func TestPlaceBid_NilBidAmount(t *testing.T) {
	db := setupTestDB(t)
	svc := newService(db)

	seedVehicle(t, db)

	_, err := svc.UpdateVehicle("test-vehicle-1", database.VehicleUpdate{BidAmount: nil})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrBidTooLow)
}

func TestPlaceBid_ZeroBidAmount(t *testing.T) {
	db := setupTestDB(t)
	svc := newService(db)

	seedVehicle(t, db)

	_, err := svc.UpdateVehicle("test-vehicle-1", database.VehicleUpdate{BidAmount: new(0)})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrBidTooLow)
}

func TestPlaceBid_VehicleNotFound(t *testing.T) {
	db := setupTestDB(t)
	svc := newService(db)

	_, err := svc.UpdateVehicle("nonexistent", database.VehicleUpdate{BidAmount: new(1000)})
	require.Error(t, err)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestPlaceBid_CreatesBidRecord(t *testing.T) {
	db := setupTestDB(t)
	svc := newService(db)

	seedVehicle(t, db)

	_, err := svc.UpdateVehicle("test-vehicle-1", database.VehicleUpdate{BidAmount: new(1500)})
	require.NoError(t, err)

	bids, err := svc.GetAllBids()
	require.NoError(t, err)
	require.Len(t, bids, 1)
	assert.Equal(t, "test-vehicle-1", bids[0].VehicleExID)
	assert.Equal(t, 1500, bids[0].BidAmount)
	assert.False(t, bids[0].IsBuyNow)
	assert.Equal(t, "2021 Honda Civic EX", bids[0].VehicleName)
}

func TestPlaceBid_BuyNowFlagged(t *testing.T) {
	db := setupTestDB(t)
	svc := newService(db)

	buyNow := 5000
	seedVehicle(t, db, func(v *database.Vehicle) {
		v.BuyNowPrice = &buyNow
	})

	_, err := svc.UpdateVehicle("test-vehicle-1", database.VehicleUpdate{BidAmount: new(5000)})
	require.NoError(t, err)

	bids, err := svc.GetAllBids()
	require.NoError(t, err)
	require.Len(t, bids, 1)
	assert.True(t, bids[0].IsBuyNow)
}

func TestPlaceBid_RegularBidNotFlaggedBuyNow(t *testing.T) {
	db := setupTestDB(t)
	svc := newService(db)

	buyNow := 5000
	seedVehicle(t, db, func(v *database.Vehicle) {
		v.BuyNowPrice = &buyNow
	})

	_, err := svc.UpdateVehicle("test-vehicle-1", database.VehicleUpdate{BidAmount: new(1000)})
	require.NoError(t, err)

	bids, err := svc.GetAllBids()
	require.NoError(t, err)
	require.Len(t, bids, 1)
	assert.False(t, bids[0].IsBuyNow)
}

func TestGetAllBids_OrderedNewestFirst(t *testing.T) {
	db := setupTestDB(t)
	svc := newService(db)

	seedVehicle(t, db)

	_, err := svc.UpdateVehicle("test-vehicle-1", database.VehicleUpdate{BidAmount: new(1000)})
	require.NoError(t, err)
	_, err = svc.UpdateVehicle("test-vehicle-1", database.VehicleUpdate{BidAmount: new(1100)})
	require.NoError(t, err)
	_, err = svc.UpdateVehicle("test-vehicle-1", database.VehicleUpdate{BidAmount: new(1200)})
	require.NoError(t, err)

	bids, err := svc.GetAllBids()
	require.NoError(t, err)
	require.Len(t, bids, 3)
	assert.Equal(t, 1200, bids[0].BidAmount)
	assert.Equal(t, 1100, bids[1].BidAmount)
	assert.Equal(t, 1000, bids[2].BidAmount)
}

func TestPlaceBid_MultipleBidsIncrementCount(t *testing.T) {
	db := setupTestDB(t)
	svc := newService(db)

	seedVehicle(t, db)

	for i := 0; i < 5; i++ {
		amount := 1000 + (i * 100)
		_, err := svc.UpdateVehicle("test-vehicle-1", database.VehicleUpdate{BidAmount: new(amount)})
		require.NoError(t, err, fmt.Sprintf("bid %d (amount %d) should succeed", i+1, amount))
	}

	vehicle, err := svc.GetVehicle("test-vehicle-1")
	require.NoError(t, err)
	assert.Equal(t, 5, vehicle.BidCount)
	assert.Equal(t, 1400, vehicle.CurrentBid)
}

func TestBuyNow_Success(t *testing.T) {
	db := setupTestDB(t)
	svc := newService(db)

	buyNow := 5000
	seedVehicle(t, db, func(v *database.Vehicle) {
		v.BuyNowPrice = &buyNow
	})

	vehicle, err := svc.BuyNow("test-vehicle-1")
	require.NoError(t, err)

	// Buy Now should not change current_bid or bid_count
	assert.Equal(t, 0, vehicle.CurrentBid)
	assert.Equal(t, 0, vehicle.BidCount)

	// But it should create a bid record
	bids, err := svc.GetAllBids()
	require.NoError(t, err)
	require.Len(t, bids, 1)
	assert.True(t, bids[0].IsBuyNow)
	assert.Equal(t, 5000, bids[0].BidAmount)
}

func TestBuyNow_NoBuyNowPrice(t *testing.T) {
	db := setupTestDB(t)
	svc := newService(db)

	seedVehicle(t, db) // no BuyNowPrice set

	_, err := svc.BuyNow("test-vehicle-1")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNoBuyNow)
}

func TestBuyNow_AuctionEnded(t *testing.T) {
	db := setupTestDB(t)
	svc := newService(db)

	buyNow := 5000
	seedVehicle(t, db, func(v *database.Vehicle) {
		v.BuyNowPrice = &buyNow
		v.AuctionStart = expiredAuctionStart()
	})

	_, err := svc.BuyNow("test-vehicle-1")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrAuctionEnded)
}

func TestBuyNow_DoesNotAffectBids(t *testing.T) {
	db := setupTestDB(t)
	svc := newService(db)

	buyNow := 5000
	seedVehicle(t, db, func(v *database.Vehicle) {
		v.BuyNowPrice = &buyNow
	})

	// Place a regular bid first
	_, err := svc.UpdateVehicle("test-vehicle-1", database.VehicleUpdate{BidAmount: new(1000)})
	require.NoError(t, err)

	// Buy Now should not change the bid state
	vehicle, err := svc.BuyNow("test-vehicle-1")
	require.NoError(t, err)
	assert.Equal(t, 1000, vehicle.CurrentBid)
	assert.Equal(t, 1, vehicle.BidCount)
}
