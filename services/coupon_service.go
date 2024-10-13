package services

import (
	"coupon/models"
	"errors"
	"fmt"
	"math"
	"time"
)

var Coupons = make(map[string]models.Coupon)

func CreateCoupon(coupon models.Coupon) error {
	if _, exists := Coupons[coupon.ID]; exists {
		return errors.New("coupon already exists")
	}
	Coupons[coupon.ID] = coupon
	return nil
}

func UpdateCoupon(couponID string, updatedCoupon models.Coupon) error {
	coupon, exists := Coupons[couponID]
	if !exists {
		return errors.New("coupon not found")
	}

	if err := validateCouponDetails(updatedCoupon.Details); err != nil {
		return err
	}

	if isNoChangesProvided(updatedCoupon) {
		return errors.New("no changes provided")
	}

	updateCouponDetails(&coupon, updatedCoupon)
	Coupons[couponID] = coupon

	return nil
}

func validateCouponDetails(details models.CouponDetails) error {
	if details.Threshold < 0 {
		return errors.New("invalid threshold: cannot be negative")
	}
	if details.Discount < 0 || details.Discount > 100 {
		return errors.New("invalid discount: must be between 0 and 100")
	}
	if details.MaxUses < 0 {
		return errors.New("invalid max uses: must be positive")
	}
	if details.Uses > details.MaxUses {
		return errors.New("uses cannot exceed max uses")
	}
	return nil
}

func isNoChangesProvided(updatedCoupon models.Coupon) bool {
	return updatedCoupon.Type == "" &&
		updatedCoupon.Details.Threshold == 0 &&
		updatedCoupon.Details.Discount == 0 &&
		updatedCoupon.Details.MinCartValue == 0 &&
		updatedCoupon.Details.MaxUses == 0 &&
		updatedCoupon.Details.ProductID == "" &&
		updatedCoupon.Details.ExpiryDate == nil &&
		len(updatedCoupon.Details.BuyProducts) == 0 &&
		len(updatedCoupon.Details.GetProducts) == 0 &&
		updatedCoupon.Details.RepetitionLimit == 0
}

func updateCouponDetails(coupon *models.Coupon, updatedCoupon models.Coupon) {
	coupon.Type = updatedCoupon.Type

	if updatedCoupon.Details.Threshold > 0 {
		coupon.Details.Threshold = updatedCoupon.Details.Threshold
	}
	if updatedCoupon.Details.Discount > 0 {
		coupon.Details.Discount = updatedCoupon.Details.Discount
	}
	if updatedCoupon.Details.MinCartValue > 0 {
		coupon.Details.MinCartValue = updatedCoupon.Details.MinCartValue
	}
	if updatedCoupon.Details.MaxUses > 0 {
		coupon.Details.MaxUses = updatedCoupon.Details.MaxUses
	}
	if updatedCoupon.Details.ProductID != "" {
		coupon.Details.ProductID = updatedCoupon.Details.ProductID
	}
	if updatedCoupon.Details.ExpiryDate != nil {
		coupon.Details.ExpiryDate = updatedCoupon.Details.ExpiryDate
	}
	coupon.Details.Exclusive = updatedCoupon.Details.Exclusive

	if len(updatedCoupon.Details.BuyProducts) > 0 {
		coupon.Details.BuyProducts = updatedCoupon.Details.BuyProducts
	}
	if len(updatedCoupon.Details.GetProducts) > 0 {
		coupon.Details.GetProducts = updatedCoupon.Details.GetProducts
	}
	if updatedCoupon.Details.RepetitionLimit > 0 {
		coupon.Details.RepetitionLimit = updatedCoupon.Details.RepetitionLimit
	}
}

func GetAllCoupons() []models.Coupon {
	coupons := make([]models.Coupon, 0, len(Coupons))
	for _, coupon := range Coupons {
		coupons = append(coupons, coupon)
	}
	return coupons
}

func GetCouponByID(id string) (models.Coupon, error) {
	coupon, exists := Coupons[id]
	if !exists {
		return models.Coupon{}, errors.New("coupon not found")
	}
	return coupon, nil
}

func DeleteCoupon(id string) error {
	if _, exists := Coupons[id]; !exists {
		return errors.New("coupon not found")
	}
	delete(Coupons, id)
	return nil
}

func ApplyCoupon(cart models.Cart, couponID string, appliedCoupons map[string]bool) (models.Cart, error) {
	if len(cart.Items) == 0 {
		return cart, errors.New("cart is empty")
	}

	coupon, exists := Coupons[couponID]
	if !exists {
		return cart, errors.New("coupon not found")
	}

	if err := validateCouponApplication(coupon, appliedCoupons); err != nil {
		return cart, err
	}

	discount, totalAmount, err := calculateDiscount(cart, coupon)
	if err != nil {
		return cart, err
	}

	appliedCoupons[couponID] = true
	coupon.Details.Uses++
	Coupons[couponID] = coupon

	discount = math.Round(discount*100) / 100
	cart.TotalPrice = totalAmount
	cart.TotalDiscount += discount
	cart.FinalPrice = totalAmount - discount

	return cart, nil
}

func validateCouponApplication(coupon models.Coupon, appliedCoupons map[string]bool) error {
	if coupon.Type == "" {
		return errors.New("invalid coupon type")
	}
	if appliedCoupons[coupon.ID] {
		return errors.New("coupon already applied")
	}
	if coupon.Details.ExpiryDate != nil && time.Now().After(*coupon.Details.ExpiryDate) {
		return errors.New("coupon has expired")
	}
	if coupon.Details.Uses >= coupon.Details.MaxUses {
		return errors.New("coupon usage limit exceeded")
	}
	if coupon.Details.Exclusive && len(appliedCoupons) > 0 {
		return errors.New("this coupon cannot be combined with others")
	}
	return nil
}

func calculateDiscount(cart models.Cart, coupon models.Coupon) (float64, float64, error) {
	totalAmount := 0.0

	for _, item := range cart.Items {
		if item.Quantity > 0 && item.Price > 0 {
			totalAmount += float64(item.Quantity) * item.Price
		}
	}

	switch coupon.Type {
	case "cart-wise":
		return calculateCartWiseDiscount(cart, coupon, totalAmount)
	case "product-wise":
		return calculateProductWiseDiscount(cart, coupon, totalAmount)
	case "bxgy":
		return calculateBxGyDiscount(cart, coupon, totalAmount)
	default:
		return 0, 0, fmt.Errorf("unsupported coupon type: %s", coupon.Type)
	}
}

func calculateCartWiseDiscount(cart models.Cart, coupon models.Coupon, totalAmount float64) (float64, float64, error) {
	if coupon.Details.Threshold <= 0 {
		return 0, 0, errors.New("invalid threshold value in cart-wise coupon")
	}
	if coupon.Details.Discount <= 0 {
		return 0, 0, errors.New("invalid discount value in cart-wise coupon")
	}
	if coupon.Details.MinCartValue > 0 && totalAmount < coupon.Details.MinCartValue {
		return 0, 0, errors.New("cart value is below the minimum required for this coupon")
	}
	if totalAmount >= coupon.Details.Threshold {
		discount := (coupon.Details.Discount / 100) * totalAmount
		return discount, totalAmount, nil
	}
	return 0, 0, errors.New("cart total does not meet the threshold for this coupon")
}

func calculateProductWiseDiscount(cart models.Cart, coupon models.Coupon, totalAmount float64) (float64, float64, error) {
	if coupon.Details.ProductID == "" {
		return 0, 0, errors.New("invalid product ID in product-wise coupon")
	}
	if coupon.Details.Discount <= 0 || coupon.Details.Discount > 100 {
		return 0, 0, errors.New("invalid discount value in product-wise coupon")
	}
	for _, item := range cart.Items {
		if item.ProductID == coupon.Details.ProductID {
			discount := (coupon.Details.Discount / 100) * (float64(item.Quantity) * item.Price)
			item.TotalDiscount += math.Round(discount*100) / 100
			return discount, totalAmount, nil
		}
	}
	return 0, 0, errors.New("coupon cannot be applied to one or more products in your cart")
}

func calculateBxGyDiscount(cart models.Cart, coupon models.Coupon, totalAmount float64) (float64, float64, error) {
	if len(coupon.Details.BuyProducts) == 0 || len(coupon.Details.GetProducts) == 0 || coupon.Details.RepetitionLimit <= 0 {
		return 0, 0, errors.New("invalid BxGy coupon details")
	}

	buyProductMap := make(map[string]int)
	for _, item := range cart.Items {
		buyProductMap[item.ProductID] += item.Quantity
	}

	maxRepetitions := coupon.Details.RepetitionLimit
	for _, buyProduct := range coupon.Details.BuyProducts {
		if buyProductMap[buyProduct.ProductID]/buyProduct.Quantity < maxRepetitions {
			maxRepetitions = buyProductMap[buyProduct.ProductID] / buyProduct.Quantity
		}
	}

	if maxRepetitions > 0 {
		discount := 0.0
		for _, getProduct := range coupon.Details.GetProducts {
			freeQuantity := getProduct.Quantity * maxRepetitions
			for _, item := range cart.Items {
				if item.ProductID == getProduct.ProductID {
					item.Quantity += freeQuantity
					discount += float64(freeQuantity) * item.Price
					item.TotalDiscount += math.Round(discount*100) / 100
					break
				}
			}
		}
		return discount, totalAmount, nil
	}
	return 0, 0, errors.New("insufficient buy products for applying BxGy coupon")
}
