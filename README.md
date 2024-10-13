# Coupon Management API

This project implements a RESTful API for managing various types of discount coupons for an e-commerce platform. The system allows for the application of different coupon types including **Cart-wise**, **Product-wise**, and **Buy X Get Y (BxGy)** offers.

## Features

- Create, Update, retrieve, and delete coupons.
- Apply applicable coupons to a shopping cart and calculate discounts.
- Handle various edge cases related to coupon usage, expiry, limits, and conflicts.

## API Endpoints

- `POST /coupons`: Create a new coupon.
- `GET /coupons`: Retrieve all coupons.
- `PUT /coupons/{id}`: Update a specific coupon by its ID.
- `GET /coupons/{id}`: Retrieve a specific coupon by its ID.
- `DELETE /coupons/{id}`: Delete a specific coupon by its ID.
- `POST /applicable-coupons`: Fetch all applicable coupons for a given cart.
- `POST /apply-coupon/{id}`: Apply a specific coupon to the cart and return the updated cart with discounted prices.

## Coupon Types

### 1. **Cart-wise Coupons**
   Apply a discount to the entire cart if the total amount exceeds a certain threshold.
   - Example: 10% off on orders over $100.
   - Edge Case: If the total cart value is below the threshold, the coupon is not applied.

### 2. **Product-wise Coupons**
   Apply a discount to specific products in the cart.
   - Example: 20% off on Product A.
   - Edge Case: If the product is excluded, the coupon is not applied.

### 3. **Buy X Get Y (BxGy) Coupons**
   Buy X quantity of certain products and get Y quantity of other products for free.
   - Example: Buy 2 items from Product A and get 1 item from Product B for free.
   - Edge Case: If the required quantity of the "Buy" products is not met, the coupon is not applicable.

## Edge Cases Handled

### 1. **Empty Cart Handling**
   - **Scenario**: If the cart is empty, no coupons can be applied.
   - **Handling**: The function checks if the cart has any items. If the cart is empty, it returns an error: `"cart is empty"`.

### 2. **Invalid Coupon Type**
   - **Scenario**: A coupon must have a valid type (i.e., `"cart-wise"`, `"product-wise"`, or `"bxgy"`).
   - **Handling**: If the coupon type is an empty string, the function returns an error: `"invalid coupon type"`.

### 3. **Coupon Expiry Handling**
   - **Scenario**: Coupons can have an expiry date, and attempting to apply an expired coupon should be disallowed.
   - **Handling**: The function checks if the `ExpiryDate` is not `nil` and whether the current date is after the expiry date. If it has expired, it returns an error: `"coupon has expired"`.

### 4. **Usage Limit Checks**
   - **Scenario**: Coupons may have a limit on how many times they can be used.
   - **Handling**: The function checks the current usage count against the maximum allowed. If the limit is exceeded, it returns an error: `"coupon usage limit exceeded"`.

### 5. **Exclusive Coupon Handling**
   - **Scenario**: Some coupons are marked as exclusive and cannot be combined with others.
   - **Handling**: If the coupon is exclusive and there are already applied coupons, the function returns an error: `"this coupon cannot be combined with others"`.

### 6. **Threshold Validation for Cart-wise Coupons**
   - **Scenario**: Cart-wise coupons may require the cart total to exceed a certain threshold for the coupon to be applied.
   - **Handling**: The function checks if the cart total is below the threshold. If it is, it returns an error: `"cart total does not meet the threshold for this coupon"`.

### 7. **Minimum Cart Value Requirement**
   - **Scenario**: Cart-wise coupons may require a minimum cart value for applicability.
   - **Handling**: The function checks if the total amount in the cart is below the minimum required value. If it is, it returns an error: `"cart value is below the minimum required for this coupon"`.

### 8. **Invalid Discount Values**
   - **Scenario**: Coupons may have invalid discount percentages (e.g., negative or zero).
   - **Handling**: The function validates discount values and returns an error if they are invalid. For example, `"invalid discount value in cart-wise coupon"` or `"invalid discount value in product-wise coupon"`.

### 9. **Excluded Products Handling**
   - **Scenario**: Some coupons may exclude specific products from applicability.
   - **Handling**: The function checks if any products in the cart are in the coupon’s `ExcludedProducts` list. If they are, it returns an error: `"coupon cannot be applied to one or more products in your cart"`.

### 10. **Invalid Product ID for Product-wise Coupons**
   - **Scenario**: Product-wise coupons must specify a valid product ID.
   - **Handling**: If the `ProductID` is empty, the function returns an error: `"invalid product ID in product-wise coupon"`.

### 11. **Invalid BxGy Coupon Details**
   - **Scenario**: BxGy coupons must have valid buy and get product details.
   - **Handling**: The function checks if the required details are present and valid. If not, it returns an error: `"invalid BxGy coupon details"`.

### 12. **Insufficient Buy Products for BxGy Coupons**
   - **Scenario**: If the required buy products are not present in sufficient quantities.
   - **Handling**: If the conditions for applying a BxGy coupon are not met, it returns an error: `"insufficient buy products for applying BxGy coupon"`.

### 13. **Floating Point Precision Handling**
   - **Scenario**: Calculating discounts might lead to floating-point precision issues.
   - **Handling**: Discounts are rounded to two decimal places to avoid discrepancies in final amounts.

### 14. **Unsupported Coupon Types**
   - **Scenario**: If a coupon has an unsupported type.
   - **Handling**: The function returns an error indicating the unsupported coupon type: `"unsupported coupon type: %s"`.

### 1. **Empty Cart Handling**
   - **Scenario**: If the cart is empty, no coupons can be applied.
   - **Handling**: The function checks if the cart has any items. If the cart is empty, it returns an error: `"cart is empty"`.

### 2. **Invalid Coupon Type**
   - **Scenario**: A coupon must have a valid type (i.e., `"cart-wise"`, `"product-wise"`, or `"bxgy"`).
   - **Handling**: If the coupon type is an empty string, the function returns an error: `"invalid coupon type"`.

### 3. **Coupon Expiry Handling**
   - **Scenario**: Coupons can have an expiry date, and attempting to apply an expired coupon should be disallowed.
   - **Handling**: The function checks if the `ExpiryDate` is not `nil` and whether the current date is after the expiry date. If it has expired, it returns an error: `"coupon has expired"`.

### 4. **Usage Limit Checks**
   - **Scenario**: Coupons may have a limit on how many times they can be used.
   - **Handling**: The function checks the current usage count against the maximum allowed. If the limit is exceeded, it returns an error: `"coupon usage limit exceeded"`.

### 5. **Exclusive Coupon Handling**
   - **Scenario**: Some coupons are marked as exclusive and cannot be combined with others.
   - **Handling**: If the coupon is exclusive and there are already applied coupons, the function returns an error: `"this coupon cannot be combined with others"`.

### 6. **Threshold Validation for Cart-wise Coupons**
   - **Scenario**: Cart-wise coupons may require the cart total to exceed a certain threshold for the coupon to be applied.
   - **Handling**: The function checks if the cart total is below the threshold. If it is, it returns an error: `"cart total does not meet the threshold for this coupon"`.

### 7. **Minimum Cart Value Requirement**
   - **Scenario**: Cart-wise coupons may require a minimum cart value for applicability.
   - **Handling**: The function checks if the total amount in the cart is below the minimum required value. If it is, it returns an error: `"cart value is below the minimum required for this coupon"`.

### 8. **Invalid Discount Values**
   - **Scenario**: Coupons may have invalid discount percentages (e.g., negative or zero).
   - **Handling**: The function validates discount values and returns an error if they are invalid. For example, `"invalid discount value in cart-wise coupon"` or `"invalid discount value in product-wise coupon"`.

### 9. **Excluded Products Handling**
   - **Scenario**: Some coupons may exclude specific products from applicability.
   - **Handling**: The function checks if any products in the cart are in the coupon’s `ExcludedProducts` list. If they are, it returns an error: `"coupon cannot be applied to one or more products in your cart"`.

### 10. **Invalid Product ID for Product-wise Coupons**
   - **Scenario**: Product-wise coupons must specify a valid product ID.
   - **Handling**: If the `ProductID` is empty, the function returns an error: `"invalid product ID in product-wise coupon"`.

### 11. **Invalid BxGy Coupon Details**
   - **Scenario**: BxGy coupons must have valid buy and get product details.
   - **Handling**: The function checks if the required details are present and valid. If not, it returns an error: `"invalid BxGy coupon details"`.

### 12. **Insufficient Buy Products for BxGy Coupons**
   - **Scenario**: If the required buy products are not present in sufficient quantities.
   - **Handling**: If the conditions for applying a BxGy coupon are not met, it returns an error: `"insufficient buy products for applying BxGy coupon"`.

### 13. **Floating Point Precision Handling**
   - **Scenario**: Calculating discounts might lead to floating-point precision issues.
   - **Handling**: Discounts are rounded to two decimal places to avoid discrepancies in final amounts.

### 14. **Unsupported Coupon Types**
   - **Scenario**: If a coupon has an unsupported type.
   - **Handling**: The function returns an error indicating the unsupported coupon type: `"unsupported coupon type: %s"`.

### 15. **Non-Existent Coupon**
   - **Scenario**: The user tries to update a coupon that does not exist.
   - **Handling**: If the coupon ID is not found, the function returns an error: `"coupon not found"`.

### 16. **Invalid Coupon Type**
   - **Scenario**: The user provides an invalid coupon type (e.g., an empty or unsupported type).
   - **Handling**: The function returns an error: `"invalid coupon type"`.

### 17. **Invalid Threshold, Discount, or MaxUses Values**
   - **Scenario**: The user tries to update the coupon with negative or zero values for `Threshold`, `Discount`, or `MaxUses`.
   - **Handling**: The function validates these fields and returns errors like:
     - `"invalid threshold value: cannot be negative"`
     - `"invalid discount value: must be between 0 and 100"`
     - `"invalid max uses: cannot be negative"`

### 18. **Uses Exceeding MaxUses**
   - **Scenario**: The user tries to set `MaxUses` lower than the current number of uses.
   - **Handling**: The function checks that the `MaxUses` is not lower than the current `Uses` and returns an error: `"uses cannot exceed max uses"`.

### 19. **Invalid Product ID for Product-wise Coupons**
   - **Scenario**: The user tries to update a product-wise coupon with an invalid or empty `ProductID`.
   - **Handling**: The function checks for a valid `ProductID` and returns an error: `"invalid product ID for product-wise coupon"`.

### 20. **Invalid Expiry Date**
   - **Scenario**: The user tries to update the coupon with an expiry date in the past.
   - **Handling**: The function checks the `ExpiryDate` and returns an error if the date is in the past: `"invalid expiry date: cannot be in the past"`.

### 21. **Repetition Limit for BxGy Coupons**
   - **Scenario**: The user tries to update a BxGy coupon with an invalid repetition limit (e.g., zero or negative).
   - **Handling**: The function validates the repetition limit for BxGy coupons and returns an error if invalid: `"invalid repetition limit for BxGy coupon"`.

## Example API Payloads

### Create a Cart-wise Coupon:

```json
{
  "id": "1",
  "type": "cart-wise",
  "details": {
    "threshold": 100.0,
    "discount": 10.0,
    "min_cart_value": 50.0,
    "expiry_date": "2024-12-31T23:59:59Z",
    "max_uses": 5,
    "uses": 0,
    "exclusive": false
  }
}
```

### Create a Product-wise Coupon:

```json
{
  "id": "2",
  "type": "product-wise",
  "details": {
    "product_id": "A123",
    "discount": 20.0,
    "excluded_products": ["B123", "C456"],
    "expiry_date": "2024-12-31T23:59:59Z",
    "max_uses": 10,
    "uses": 0,
    "exclusive": true
  }
}
```

### Create a Buy X Get Y (BxGy) Coupon:

```json
{
  "id": "3",
  "type": "bxgy",
  "details": {
    "buy_products": [
      { "product_id": "A123", "quantity": 2 }
    ],
    "get_products": [
      { "product_id": "B456", "quantity": 1 }
    ],
    "repetition_limit": 3,
    "expiry_date": "2024-12-31T23:59:59Z",
    "max_uses": 15,
    "uses": 0,
    "exclusive": false
  }
}
```

### Apply a Specific Coupon to a Cart:

```json
{
  "cart": {
    "items": [
      { "product_id": "A123", "quantity": 2, "price": 100.0 },
      { "product_id": "B456", "quantity": 1, "price": 50.0 }
    ]
  }
}
```

### Get Applicable Coupons for a Cart:

```json
{
  "cart": {
    "items": [
      { "product_id": "A123", "quantity": 2, "price": 100.0 },
      { "product_id": "B456", "quantity": 1, "price": 50.0 }
    ]
  }
}
```

### Update a Specific Coupon:

```json
{
  "type": "cart-wise",
  "details": {
    "discount": 15.0,
    "threshold": 200.0,
    "expiry_date": "2024-12-31T23:59:59Z",
    "max_uses": 5,
    "exclusive": true
  }
}
```
