package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/marco-almeida/challenge-02/internal"
	"github.com/rs/zerolog/log"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, ginErr := range c.Errors {
			logLevel := log.Info()

			unwrappedErr := ginErr.Err
			var validationErrors validator.ValidationErrors
			switch {
			// check if error has validator.ValidationErrors
			case errors.As(unwrappedErr, &validationErrors):
				errorResponse := internal.ErrorResponse{
					Error: http.StatusText(http.StatusBadRequest),
				}
				for e := range validationErrors {
					errorResponse.Validations = append(errorResponse.Validations, internal.ValidationError{
						// field should be json representation of field
						Field:   validationErrors[e].Field(),
						Tag:     validationErrors[e].Tag(),
						Message: validationErrorToText(validationErrors[e]),
					})
				}
				c.JSON(http.StatusBadRequest, errorResponse)
			case errors.Is(unwrappedErr, internal.ErrOrderTooHeavy):
				c.JSON(http.StatusBadRequest, gin.H{"error": "vehicle can not handle this order's weight"})
			case errors.Is(unwrappedErr, internal.ErrOrderAlreadyFinished):
				c.JSON(http.StatusBadRequest, gin.H{"error": "order already finished"})
			case errors.Is(unwrappedErr, internal.ErrVehicleNotFound):
				c.JSON(http.StatusBadRequest, gin.H{"error": "vehicle not found"})
			case errors.Is(unwrappedErr, internal.ErrOrderNotFound):
				c.JSON(http.StatusBadRequest, gin.H{"error": "order not found"})
			case errors.Is(unwrappedErr, internal.ErrVehicleAlreadyExists):
				c.JSON(http.StatusBadRequest, gin.H{"error": "vehicle already exists"})
			case errors.Is(unwrappedErr, internal.ErrInvalidParams):
				c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusText(http.StatusBadRequest)})
			case errors.Is(unwrappedErr, internal.ErrForeignKeyConstraintViolation):
				c.JSON(http.StatusConflict, gin.H{"error": http.StatusText(http.StatusConflict)})
			case errors.Is(unwrappedErr, internal.ErrUniqueConstraintViolation):
				c.JSON(http.StatusConflict, gin.H{"error": http.StatusText(http.StatusConflict)})
			case errors.Is(unwrappedErr, internal.ErrNoRows):
				c.JSON(http.StatusNotFound, gin.H{"error": http.StatusText(http.StatusNotFound)})
			default:
				logLevel = log.Error()
				c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
			}
			logLevel.Err(unwrappedErr).Send()
		}
	}
}

func validationErrorToText(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s", e.Field(), e.Param())
	case "min":
		return fmt.Sprintf("%s must be at least than %s", e.Field(), e.Param())
	case "email":
		return "Invalid email format"
	case "len":
		return fmt.Sprintf("%s must be %s characters long", e.Field(), e.Param())
	}
	return fmt.Sprintf("%s is not valid", e.Field())
}
