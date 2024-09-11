package bootstrap

import (
	"github.com/simondanielsson/apPRoved/cmd/internal/controllers"
	"github.com/simondanielsson/apPRoved/cmd/internal/repositories"
	"github.com/simondanielsson/apPRoved/cmd/internal/services"
)

func InitRepositories() *repositories.Repositories {
	return &repositories.Repositories{
		ReviewsRepository: repositories.NewReviewsRepository(),
		UserRepository:    repositories.NewUserRepository(),
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
		AuthController:    controllers.NewAuthController(services.AuthService, services.UserService),
	}
}
