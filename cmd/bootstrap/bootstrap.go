package bootstrap

import (
	"github.com/simondanielsson/apPRoved/cmd/internal/controllers"
	"github.com/simondanielsson/apPRoved/cmd/internal/repositories"
	"github.com/simondanielsson/apPRoved/cmd/internal/services"
	"gorm.io/gorm"
)

func InitRepositories(db *gorm.DB) *repositories.Repositories {
	return &repositories.Repositories{
		ReviewsRepository: repositories.NewReviewsRepository(db),
		UserRepository:    repositories.NewUserRepository(db),
	}
}

func InitServices(repos *repositories.Repositories) *services.Services {
	return &services.Services{
		ReviewsService: services.NewReviewsService(repos.ReviewsRepository),
		UserService:    services.NewUserService(repos.UserRepository),
		AuthService:    services.NewAuthService(repos.UserRepository),
	}
}

func InitControllers(services *services.Services) *controllers.Controllers {
	return &controllers.Controllers{
		ReviewsController: controllers.NewReviewsController(services.ReviewsService),
		UserController:    controllers.NewUserController(services.UserService),
		AuthController:    controllers.NewAuthController(services.AuthService),
	}
}
