package controllers

import (
	"coupon/models"
	"coupon/services"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	invalidRequestBody        = "Invalid request body"
	internalServerError       = "Internal server error"
	couponUpdatedSuccessfully = "Coupon updated successfully"
	couponDeletedSuccessfully = "Coupon deleted successfully"
)

func ApplyCoupon(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	couponID := params["id"]

	var cartRequest struct {
		Cart models.Cart `json:"cart"`
	}

	if err := json.NewDecoder(r.Body).Decode(&cartRequest); err != nil {
		handleError(w, invalidRequestBody, http.StatusBadRequest)
		return
	}

	appliedCoupons := make(map[string]bool)

	updatedCart, err := services.ApplyCoupon(cartRequest.Cart, couponID, appliedCoupons)
	if err != nil {
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"updated_cart": updatedCart,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}

func CreateCoupon(w http.ResponseWriter, r *http.Request) {
	var coupon models.Coupon
	if err := json.NewDecoder(r.Body).Decode(&coupon); err != nil {
		handleError(w, invalidRequestBody, http.StatusBadRequest)
		return
	}

	if err := services.CreateCoupon(coupon); err != nil {
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(coupon); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}

func UpdateCoupon(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	couponID := params["id"]

	var updatedCoupon models.Coupon
	if err := json.NewDecoder(r.Body).Decode(&updatedCoupon); err != nil {
		handleError(w, invalidRequestBody, http.StatusBadRequest)
		return
	}

	if err := services.UpdateCoupon(couponID, updatedCoupon); err != nil {
		handleError(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": couponUpdatedSuccessfully,
	}); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}

func GetAllCoupons(w http.ResponseWriter, r *http.Request) {
	coupons := services.GetAllCoupons()
	if err := json.NewEncoder(w).Encode(coupons); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}

func GetCouponByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	coupon, err := services.GetCouponByID(params["id"])
	if err != nil {
		handleError(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(coupon); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}

func DeleteCoupon(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if err := services.DeleteCoupon(params["id"]); err != nil {
		handleError(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": couponDeletedSuccessfully,
	}); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}

func GetApplicableCoupons(w http.ResponseWriter, r *http.Request) {
	var cartRequest struct {
		Cart models.Cart `json:"cart"`
	}

	if err := json.NewDecoder(r.Body).Decode(&cartRequest); err != nil {
		handleError(w, invalidRequestBody, http.StatusBadRequest)
		return
	}

	applicableCoupons := []map[string]interface{}{}
	appliedCoupons := make(map[string]bool)

	for _, coupon := range services.GetAllCoupons() {
		updatedCart, err := services.ApplyCoupon(cartRequest.Cart, coupon.ID, appliedCoupons)
		if err == nil && updatedCart.TotalDiscount > 0 {
			applicableCoupons = append(applicableCoupons, map[string]interface{}{
				"coupon_id": coupon.ID,
				"type":      coupon.Type,
				"discount":  updatedCart.TotalDiscount,
			})
		}
	}

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"applicable_coupons": applicableCoupons,
	}); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}

func handleError(w http.ResponseWriter, message string, statusCode int) {
	log.Printf("Error: %s", message)
	http.Error(w, message, statusCode)
}
