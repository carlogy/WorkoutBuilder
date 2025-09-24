package services

import (
	"context"
	"fmt"
	"time"

	db "github.com/carlogy/WorkoutBuilder/internal/database"
	"github.com/carlogy/WorkoutBuilder/internal/repositories"
	id "github.com/google/uuid"
)

type ExperienceLevel string

const (
	Beginner     ExperienceLevel = "beginner"
	Intermediate ExperienceLevel = "intermediate"
	Advanced     ExperienceLevel = "advanced"
)

type Plan struct {
	Id              id.UUID         `json:"id,omitempty"`
	Name            string          `json:"name"`
	Goal            string          `json:"goal"`
	Days            int             `json:"days"`
	Duration        string          `json:"duration"`
	Description     string          `json:"description"`
	Workouts        []Workout       `json:"workouts"`
	ExperienceLevel ExperienceLevel `json:"experienceLevel"`
	CreatedAt       *time.Time      `json:"createdAt,omitempty"`
	ModifiedAt      *time.Time      `json:"modifiedAt,omitempty"`
}

type PlanService struct {
	planRepo *repositories.PlanRepository
}

type PlanRequestParams struct {
	Name            string          `json:"name"`
	Goal            string          `json:"goal"`
	Days            int             `json:"days"`
	Duration        string          `json:"duration"`
	Description     string          `json:"description"`
	Workouts        []Workout       `json:"workouts"`
	ExperienceLevel ExperienceLevel `json:"experienceLevel"`
}

func NewPlanService(repo repositories.PlanRepository) *PlanService {
	return &PlanService{
		planRepo: &repo,
	}
}

func (ps *PlanService) createDBPlanParams(jep PlanRequestParams) db.CreatePlanParams {

	return db.CreatePlanParams{
		Name:            jep.Name,
		Goal:            jep.Goal,
		Days:            NoneNullIntToNullInt32(&jep.Days),
		Duration:        NoneNullToNullString(&jep.Duration),
		Description:     NoneNullToNullString(&jep.Description),
		ExperienceLevel: NoneNullToNullString((*string)(&jep.ExperienceLevel)),
	}

}

func (ps *PlanService) CreateNewPlan(ctx context.Context, jep PlanRequestParams) (Plan, error) {

	dbPlan, err := ps.planRepo.CreatePlan(ctx, ps.createDBPlanParams(jep))
	if err != nil {
		fmt.Println("Error creating plan from repo: ", err)
		return Plan{}, err
	}

	fmt.Println(dbPlan)

	return Plan{}, nil
}
