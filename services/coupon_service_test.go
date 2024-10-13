package services

import (
	"coupon/models"
	"testing"
	"time"
)

func TestCreateCoupon(t *testing.T) {
	Coupons = make(map[string]models.Coupon)

	coupon := models.Coupon{
		ID:   "1",
		Type: "cart-wise",
		Details: models.CouponDetails{
			Threshold:    100.0,
			Discount:     10.0,
			MinCartValue: 50.0,
			MaxUses:      5,
			Uses:         0,
			Exclusive:    false,
		},
	}

	err := CreateCoupon(coupon)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = CreateCoupon(coupon)
	if err == nil {
		t.Fatalf("Expected error, got none")
	}
}

func TestUpdateCoupon(t *testing.T) {
	Coupons = make(map[string]models.Coupon)

	coupon := models.Coupon{
		ID:   "1",
		Type: "cart-wise",
		Details: models.CouponDetails{
			Threshold:    100.0,
			Discount:     10.0,
			MinCartValue: 50.0,
			MaxUses:      5,
			Uses:         0,
			Exclusive:    false,
		},
	}

	err := CreateCoupon(coupon)
	if err != nil {
		t.Fatalf("Expected no error when creating the coupon, got %v", err)
	}

	updatedCoupon := models.Coupon{
		Type: "cart-wise",
		Details: models.CouponDetails{
			Threshold: 200.0,
			Discount:  20.0,
		},
	}

	err = UpdateCoupon(coupon.ID, updatedCoupon)
	if err != nil {
		t.Fatalf("Expected no error when updating the coupon, got %v", err)
	}

	updated, err := GetCouponByID(coupon.ID)
	if err != nil {
		t.Fatalf("Expected to retrieve the updated coupon, got error %v", err)
	}

	if updated.Details.Threshold != 200.0 {
		t.Fatalf("Expected updated threshold to be 200.0, got %f", updated.Details.Threshold)
	}

	if updated.Details.Discount != 20.0 {
		t.Fatalf("Expected updated discount to be 20.0, got %f", updated.Details.Discount)
	}

	if updated.Details.MinCartValue != 50.0 {
		t.Fatalf("Expected min cart value to remain 50.0, got %f", updated.Details.MinCartValue)
	}

	if updated.Details.Uses != 0 {
		t.Fatalf("Expected uses to remain 0, got %d", updated.Details.Uses)
	}

	if updated.Details.MaxUses != 5 {
		t.Fatalf("Expected max uses to remain 5, got %d", updated.Details.MaxUses)
	}
}

func TestUpdateCoupon_EdgeCases(t *testing.T) {
	Coupons = make(map[string]models.Coupon)

	coupon := models.Coupon{
		ID:   "1",
		Type: "cart-wise",
		Details: models.CouponDetails{
			Threshold: 100.0,
			Discount:  10.0,
			MaxUses:   5,
			Uses:      0,
		},
	}
	Coupons[coupon.ID] = coupon

	err := UpdateCoupon("999", coupon)
	if err == nil || err.Error() != "coupon not found" {
		t.Fatalf("Expected 'coupon not found', got %v", err)
	}

	invalidCoupon := models.Coupon{
		Details: models.CouponDetails{
			Threshold: -10.0,
		},
	}
	err = UpdateCoupon("1", invalidCoupon)
	if err == nil || err.Error() != "invalid threshold: cannot be negative" {
		t.Fatalf("Expected 'invalid threshold: cannot be negative', got %v", err)
	}

	invalidCoupon = models.Coupon{
		Details: models.CouponDetails{
			Discount: 150.0,
		},
	}
	err = UpdateCoupon("1", invalidCoupon)
	if err == nil || err.Error() != "invalid discount: must be between 0 and 100" {
		t.Fatalf("Expected 'invalid discount: must be between 0 and 100', got %v", err)
	}

	emptyCoupon := models.Coupon{}
	err = UpdateCoupon("1", emptyCoupon)
	if err == nil || err.Error() != "no changes provided" {
		t.Fatalf("Expected 'no changes provided', got %v", err)
	}
}

func TestGetAllCoupons(t *testing.T) {
	Coupons = make(map[string]models.Coupon)

	coupon1 := models.Coupon{
		ID:   "1",
		Type: "cart-wise",
		Details: models.CouponDetails{
			Threshold:    100.0,
			Discount:     10.0,
			MinCartValue: 50.0,
		},
	}
	coupon2 := models.Coupon{
		ID:   "2",
		Type: "product-wise",
		Details: models.CouponDetails{
			ProductID: "A123",
			Discount:  20.0,
		},
	}

	CreateCoupon(coupon1)
	CreateCoupon(coupon2)

	coupons := GetAllCoupons()
	if len(coupons) != 2 {
		t.Fatalf("Expected 2 coupons, got %d", len(coupons))
	}
}

func TestGetCouponByID(t *testing.T) {
	Coupons = make(map[string]models.Coupon)

	coupon := models.Coupon{
		ID:   "1",
		Type: "cart-wise",
		Details: models.CouponDetails{
			Threshold: 100.0,
			Discount:  10.0,
		},
	}

	CreateCoupon(coupon)

	retrievedCoupon, err := GetCouponByID("1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if retrievedCoupon.ID != coupon.ID {
		t.Fatalf("Expected coupon ID %s, got %s", coupon.ID, retrievedCoupon.ID)
	}

	_, err = GetCouponByID("2")
	if err == nil {
		t.Fatalf("Expected error for non-existing coupon, got none")
	}
}

func TestDeleteCoupon(t *testing.T) {
	Coupons = make(map[string]models.Coupon)

	coupon := models.Coupon{
		ID:   "1",
		Type: "cart-wise",
		Details: models.CouponDetails{
			Threshold: 100.0,
			Discount:  10.0,
		},
	}

	CreateCoupon(coupon)

	err := DeleteCoupon("1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = DeleteCoupon("1")
	if err == nil {
		t.Fatalf("Expected error for non-existing coupon, got none")
	}
}

func TestApplyCoupon_EmptyCart(t *testing.T) {
	cart := models.Cart{Items: []models.CartItem{}}
	_, err := ApplyCoupon(cart, "1", make(map[string]bool))
	if err == nil || err.Error() != "cart is empty" {
		t.Fatalf("Expected 'cart is empty' error, got %v", err)
	}
}

func TestApplyCoupon_InvalidCoupon(t *testing.T) {
	cart := models.Cart{
		Items: []models.CartItem{
			{ProductID: "A123", Quantity: 1, Price: 100.0},
		},
	}

	_, err := ApplyCoupon(cart, "non_existing_coupon", make(map[string]bool))
	if err == nil || err.Error() != "coupon not found" {
		t.Fatalf("Expected 'coupon not found' error, got %v", err)
	}
}

func TestApplyCoupon_ExpiredCoupon(t *testing.T) {
	Coupons = make(map[string]models.Coupon)

	expiryDate := time.Now().Add(-24 * time.Hour)
	coupon := models.Coupon{
		ID:   "1",
		Type: "cart-wise",
		Details: models.CouponDetails{
			Threshold:  100.0,
			Discount:   10.0,
			ExpiryDate: &expiryDate,
		},
	}
	CreateCoupon(coupon)

	cart := models.Cart{
		Items: []models.CartItem{
			{ProductID: "A123", Quantity: 1, Price: 100.0},
		},
	}

	_, err := ApplyCoupon(cart, "1", make(map[string]bool))
	if err == nil || err.Error() != "coupon has expired" {
		t.Fatalf("Expected 'coupon has expired' error, got %v", err)
	}
}

func TestApplyCoupon_UsageLimitExceeded(t *testing.T) {
	Coupons = make(map[string]models.Coupon)

	coupon := models.Coupon{
		ID:   "1",
		Type: "cart-wise",
		Details: models.CouponDetails{
			Threshold: 100.0,
			Discount:  10.0,
			MaxUses:   1,
			Uses:      1,
		},
	}
	CreateCoupon(coupon)

	cart := models.Cart{
		Items: []models.CartItem{
			{ProductID: "A123", Quantity: 1, Price: 100.0},
		},
	}

	_, err := ApplyCoupon(cart, "1", make(map[string]bool))
	if err == nil || err.Error() != "coupon usage limit exceeded" {
		t.Fatalf("Expected 'coupon usage limit exceeded' error, got %v", err)
	}
}

func TestApplyCoupon_ExclusiveCoupon(t *testing.T) {
	Coupons = make(map[string]models.Coupon)

	coupon1 := models.Coupon{
		ID:   "1",
		Type: "cart-wise",
		Details: models.CouponDetails{
			Threshold: 100.0,
			Discount:  10.0,
			Exclusive: true,
			MaxUses:   1,
			Uses:      0,
		},
	}

	coupon2 := models.Coupon{
		ID:   "2",
		Type: "cart-wise",
		Details: models.CouponDetails{
			Threshold: 100.0,
			Discount:  20.0,
			Exclusive: false,
		},
	}

	CreateCoupon(coupon1)
	CreateCoupon(coupon2)

	cart := models.Cart{
		Items: []models.CartItem{
			{ProductID: "A123", Quantity: 1, Price: 100.0},
		},
	}

	appliedCoupons := make(map[string]bool)

	_, err := ApplyCoupon(cart, "1", appliedCoupons)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if Coupons["1"].Details.Uses != 1 {
		t.Fatalf("Expected coupon1 usage count to be 1, got %d", Coupons["1"].Details.Uses)
	}
}
