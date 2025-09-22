package repositories

import (
	"context"
	"fmt"

	db "github.com/carlogy/WorkoutBuilder/internal/database"
)

type PlanRepository struct {
	db *db.Queries
}

type PlansReposiotry interface {
	CreatePlan(ctx context.Context, p db.CreatePlanParams) (db.Plan, error)
}

func NewPlanRepository(db *db.Queries) *PlanRepository {
	return &PlanRepository{
		db: db,
	}
}

func (pr *PlanRepository) CreatePlan(ctx context.Context, p db.CreatePlanParams) (db.Plan, error) {

	dbPlan, err := pr.db.CreatePlan(ctx, p)
	if err != nil {
		fmt.Println("error creating plan in db: ", err)
		return db.Plan{}, err
	}

	return dbPlan, nil
}
