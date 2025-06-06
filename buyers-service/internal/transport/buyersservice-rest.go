// GENERATED BY 'T'ransport 'G'enerator. DO NOT EDIT.
package transport

import (
	"context"
	"github.com/gofiber/fiber/v2"
)

func (http *httpBuyersService) register(ctx context.Context, request requestBuyersServiceRegister) (response responseBuyersServiceRegister, err error) {

	response.UserID, err = http.svc.Register(ctx, request.Email, request.Password, request.Buyer)
	if err != nil {
		if http.errorHandler != nil {
			err = http.errorHandler(err)
		}
	}
	return
}
func (http *httpBuyersService) serveRegister(ctx *fiber.Ctx) (err error) {

	var request requestBuyersServiceRegister
	ctx.Response().SetStatusCode(200)
	if err = ctx.BodyParser(&request); err != nil {
		ctx.Response().SetStatusCode(fiber.StatusBadRequest)
		_, err = ctx.WriteString("request body could not be decoded: " + err.Error())
		return
	}

	var response responseBuyersServiceRegister
	if response, err = http.register(ctx.UserContext(), request); err == nil {
		var iResponse interface{} = response
		if redirect, ok := iResponse.(withRedirect); ok {
			return ctx.Redirect(redirect.RedirectTo())
		}

		return sendResponse(ctx, response)
	}
	if errCoder, ok := err.(withErrorCode); ok {
		ctx.Status(errCoder.Code())
	} else {
		ctx.Status(fiber.StatusInternalServerError)
	}
	return sendResponse(ctx, err)
}
func (http *httpBuyersService) login(ctx context.Context, request requestBuyersServiceLogin) (response responseBuyersServiceLogin, err error) {

	response.Token, err = http.svc.Login(ctx, request.Email, request.Password)
	if err != nil {
		if http.errorHandler != nil {
			err = http.errorHandler(err)
		}
	}
	return
}
func (http *httpBuyersService) serveLogin(ctx *fiber.Ctx) (err error) {

	var request requestBuyersServiceLogin
	ctx.Response().SetStatusCode(200)
	if err = ctx.BodyParser(&request); err != nil {
		ctx.Response().SetStatusCode(fiber.StatusBadRequest)
		_, err = ctx.WriteString("request body could not be decoded: " + err.Error())
		return
	}

	var response responseBuyersServiceLogin
	if response, err = http.login(ctx.UserContext(), request); err == nil {
		var iResponse interface{} = response
		if redirect, ok := iResponse.(withRedirect); ok {
			return ctx.Redirect(redirect.RedirectTo())
		}

		return sendResponse(ctx, response)
	}
	if errCoder, ok := err.(withErrorCode); ok {
		ctx.Status(errCoder.Code())
	} else {
		ctx.Status(fiber.StatusInternalServerError)
	}
	return sendResponse(ctx, err)
}
func (http *httpBuyersService) getBuyer(ctx context.Context, request requestBuyersServiceGetBuyer) (response responseBuyersServiceGetBuyer, err error) {

	response.Res, err = http.svc.GetBuyer(ctx, request.Id)
	if err != nil {
		if http.errorHandler != nil {
			err = http.errorHandler(err)
		}
	}
	return
}
func (http *httpBuyersService) serveGetBuyer(ctx *fiber.Ctx) (err error) {

	var request requestBuyersServiceGetBuyer
	ctx.Response().SetStatusCode(200)
	if err = ctx.BodyParser(&request); err != nil {
		ctx.Response().SetStatusCode(fiber.StatusBadRequest)
		_, err = ctx.WriteString("request body could not be decoded: " + err.Error())
		return
	}

	var response responseBuyersServiceGetBuyer
	if response, err = http.getBuyer(ctx.UserContext(), request); err == nil {
		var iResponse interface{} = response
		if redirect, ok := iResponse.(withRedirect); ok {
			return ctx.Redirect(redirect.RedirectTo())
		}

		return sendResponse(ctx, response)
	}
	if errCoder, ok := err.(withErrorCode); ok {
		ctx.Status(errCoder.Code())
	} else {
		ctx.Status(fiber.StatusInternalServerError)
	}
	return sendResponse(ctx, err)
}
func (http *httpBuyersService) deleteUser(ctx context.Context, request requestBuyersServiceDeleteUser) (response responseBuyersServiceDeleteUser, err error) {

	err = http.svc.DeleteUser(ctx, request.Id)
	if err != nil {
		if http.errorHandler != nil {
			err = http.errorHandler(err)
		}
	}
	return
}
func (http *httpBuyersService) serveDeleteUser(ctx *fiber.Ctx) (err error) {

	var request requestBuyersServiceDeleteUser
	ctx.Response().SetStatusCode(200)
	if err = ctx.BodyParser(&request); err != nil {
		ctx.Response().SetStatusCode(fiber.StatusBadRequest)
		_, err = ctx.WriteString("request body could not be decoded: " + err.Error())
		return
	}

	var response responseBuyersServiceDeleteUser
	if response, err = http.deleteUser(ctx.UserContext(), request); err == nil {
		var iResponse interface{} = response
		if redirect, ok := iResponse.(withRedirect); ok {
			return ctx.Redirect(redirect.RedirectTo())
		}

		return sendResponse(ctx, response)
	}
	if errCoder, ok := err.(withErrorCode); ok {
		ctx.Status(errCoder.Code())
	} else {
		ctx.Status(fiber.StatusInternalServerError)
	}
	return sendResponse(ctx, err)
}
