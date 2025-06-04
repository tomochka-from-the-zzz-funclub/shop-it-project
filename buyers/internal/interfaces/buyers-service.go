package interfaces

import (
	"context"

	"buyers-service/internal/models/request"
	"buyers-service/internal/models/response"

	"github.com/google/uuid"
)

// @tg http-prefix=api/v1
// @tg http-server log metrics
type BuyersService interface {
	// @tg http-path=/auth/register
	// @tg summary=`Регистрирует нового пользователя и создаёт связанного покупателя`
	// @tg http-method=POST
	// @tg http-success=200
	// @tg 400=`неправильное тело запроса или параметры`
	// @tg 409=`email уже существует`
	Register(ctx context.Context, email, password string, buyer request.BuyerCreate) (userID uuid.UUID, err error)

	// @tg http-path=/auth/login
	// @tg summary=`Аутентифицирует пользователя и возвращает JWT-токен`
	// @tg http-method=POST
	// @tg http-success=200
	// @tg 400=`неправильное тело запроса или параметры`
	// @tg 401=`неверные учетные данные`
	Login(ctx context.Context, email, password string) (token string, err error)

	// @tg http-path=/buyer/
	// @tg http-headers=Authorization
	// @tg summary=`Получает информацию о покупателе, связанном с текущим пользователем`
	// @tg http-method=GET
	// @tg http-success=200
	// @tg 401=`неавторизованный доступ`
	// @tg 404=`покупатель не найден`
	GetBuyer(ctx context.Context, id uuid.UUID) (res response.BuyerInfo, err error)

	// @tg http-path=/buyer/
	// @tg http-headers=Authorization
	// @tg summary=`Удаляет текущего пользователя и связанную информацию о покупателе`
	// @tg http-method=DELETE
	// @tg http-success=200
	// @tg 401=`неавторизованный доступ`
	// @tg 404=`пользователь не найден`
	DeleteUser(ctx context.Context, id uuid.UUID) (err error)
}
