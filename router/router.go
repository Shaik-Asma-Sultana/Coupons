package router

import (
	"coupon/controllers"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/coupons", controllers.CreateCoupon).Methods("POST")
	router.HandleFunc("/coupons", controllers.GetAllCoupons).Methods("GET")
	router.HandleFunc("/coupons/{id}", controllers.GetCouponByID).Methods("GET")
	router.HandleFunc("/coupons/{id}", controllers.UpdateCoupon).Methods("PUT")
	router.HandleFunc("/coupons/{id}", controllers.DeleteCoupon).Methods("DELETE")
	router.HandleFunc("/applicable-coupons", controllers.GetApplicableCoupons).Methods("POST")
	router.HandleFunc("/apply-coupon/{id}", controllers.ApplyCoupon).Methods("POST")

	return router
}
