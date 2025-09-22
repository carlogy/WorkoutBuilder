package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/carlogy/WorkoutBuilder/internal/auth"
	"github.com/carlogy/WorkoutBuilder/internal/services"
)

type PlanHandler struct {
	planService *services.PlanService
	authService *services.AuthService
}

func NewPlanHandler(ps services.PlanService, as services.AuthService) *PlanHandler {
	return &PlanHandler{&ps, &as}
}

func (ph PlanHandler) writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Conten-Type", "application/json")

	dat, err := json.Marshal(data)

	if err != nil {
		http.Error(w, "Error serializing response", 500)
		fmt.Printf("Error marshalling db response: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(dat)
	if err != nil {
		log.Printf("Error writing response: %v\n", err)
		return
	}
}

func (ph *PlanHandler) CreatePlanHandler(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, "Invalid Bearer Token", http.StatusUnauthorized)
		fmt.Println("Error gettting auth token: ", err)
		return
	}

	_, err = auth.ValidateJWT(token, ph.authService.Secret)
	if err != nil {
		http.Error(w, "Invalid Token", http.StatusUnauthorized)
		fmt.Println("Error validating JWT token: ", err)
		return
	}

	jpp := services.PlanRequestParams{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&jpp)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusInternalServerError)
		return
	}

	plan, err := ph.planService.CreateNewPlan(r.Context(), jpp)
	if err != nil {
		fmt.Println("Error from service creating plan: ", err)
	}

	ph.writeJSONResponse(w, plan, http.StatusOK)
}
