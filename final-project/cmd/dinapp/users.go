package main

import (
	"errors"
	"final-project/pkg/dinapp/model"
	"final-project/pkg/dinapp/validator"
	"log"
	"net/http"
	"time"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// err := app.readJSON(w, r, &input)
	// if err != nil {
	// 	app.badRequestResponse(w, r, err)
	// 	return
	// }
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload badRequestResponse")
		log.Printf("Error registering user1: %v", err)
		return
	}
	user := &model.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		log.Printf("Error registering user2: %v", err)
		return
	}
	v := validator.New()

	if model.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// 	err = app.mailer.Send(user.Email, "user_welcome.tmpl", user)
	// 	if err != nil {
	// 		app.serverErrorResponse(w, r, err)
	// 		return
	// 	}

	//		err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	//		if err != nil {
	//			app.serverErrorResponse(w, r, err)
	//		}
	//	}
	err = app.models.Permissions.AddForUser(int64(user.Id), "menus:read")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Assuming app.models.Token represents your token storage model
	// After the user record has been created in the database, generate a new activation
	// token for the user.

	token, err := app.models.Token.New(int64(user.Id), 3*24*time.Hour, model.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var res struct {
		Token *string     `json:"token"`
		User  *model.User `json:"user"`
	}

	res.Token = &token.Plaintext
	res.User = user

	app.writeJSON(w, http.StatusCreated, envelope{"user": res}, nil)
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TokenPlaintext string `json:"token"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if model.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetForToken(model.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	user.Activated = true

	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.Token.DeleteAllForUser(model.ScopeActivation, int64(user.Id))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
}
