package controllers

import (
	"net/http"
	"tigerhallKittens/app/lib/web"
	"tigerhallKittens/app/models"
	"tigerhallKittens/app/services"
	"tigerhallKittens/app/utils"
)

type UserController interface {
	CreateUser(request *web.Request, w http.ResponseWriter) (*web.JSONResponse, web.ErrorInterface)
}

type userController struct {
	userService services.UserService
}

func NewUserController() UserController {
	return &userController{
		userService: services.NewUserService(),
	}
}

// CreateUser creates a new user with the given details.
func (uc *userController) CreateUser(request *web.Request, w http.ResponseWriter) (*web.JSONResponse, web.ErrorInterface) {
	var createUserRequest = new(models.CreateUserRequest)

	err := request.ParseAndValidateBody(createUserRequest)
	if err != nil {
		return nil, web.ErrBadRequest(err.Error())
	}

	response, serviceErr := uc.userService.CreateUser(request.Context(), createUserRequest)
	if serviceErr != nil {
		return nil, web.NewError(web.InternalServerError, serviceErr.Error(), "", http.StatusInternalServerError)
	}

	var jsonResponse web.JSONResponse
	jsonResponse, err = utils.StructToMap(response)
	if err != nil {
		return nil, web.NewError(web.InternalServerError, serviceErr.Error(), "", http.StatusInternalServerError)
	}

	return &jsonResponse, nil
}
