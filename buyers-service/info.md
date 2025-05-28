curl-запросы для BuyersService API
1. Регистрация пользователя (Register)
Создаёт нового пользователя и связанного покупателя.
curl -X POST http://127.0.0.1:80/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user1@example.com",
    "password": "password123",
    "buyer": {
      "name": "Анна Иванова",
      "phone": 79991234567,
      "gender": true,
      "birthdate": "1990-05-15T00:00:00Z"
    }
  }'


Ожидаемый ответ (200):{"user_id":"550e8400-e29b-41d4-a716-446655440000"}


Ошибки:
400: Неверное тело (например, пустой email или name):{"error":"invalid email, password, or buyer name"}


409: Email уже существует:{"error":"email already exists"}





2. Аутентификация (Login)
Аутентифицирует пользователя и возвращает JWT-токен.
curl -X POST http://127.0.0.1:80/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user1@example.com",
    "password": "password123"
  }'


Ожидаемый ответ (200):{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}


Ошибки:
400: Неверное тело (например, пустой email):{"error":"invalid email or password"}


401: Неверные учетные данные:{"error":"invalid credentials"}



3. Получение информации о покупателе (GetBuyer)
Возвращает данные покупателя, связанного с пользователем.
curl -X GET http://127.0.0.1:80/api/v1/buyer/ \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."


Ожидаемый ответ (200):{
  "ID": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Анна Иванова",
  "phone": 79991234567,
  "gender": true,
  "birthdate": "1990-05-15T00:00:00Z"
}


Ошибки:
401: Отсутствует или неверный токен:{"error":"invalid token"}


404: Покупатель не найден:{"error":"buyer not found"}





4. Удаление пользователя (DeleteUser)
Удаляет пользователя и связанного покупателя.
curl -X DELETE http://127.0.0.1:80/api/v1/buyer/ \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."


Ожидаемый ответ (200): Пустое тело.
Ошибки:
401: Отсутствует или неверный токен:{"error":"invalid token"}


404: Пользователь не найден:{"error":"user not found"}





Примечания

Замените <token> на токен, полученный из /auth/login.


