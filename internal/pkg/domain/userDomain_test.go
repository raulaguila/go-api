package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"github.com/raulaguila/packhub"

	"github.com/raulaguila/go-api/internal/pkg/dto"
)

func TestUserTableName(t *testing.T) {
	user := new(User)
	assert.Equal(t, UserTableName, user.TableName())
}

func TestUserToMap(t *testing.T) {
	password := "hashedPassword"
	token := "token"
	user := User{
		Name:   "Alice",
		Email:  "alice@example.com",
		AuthID: 123,
		Auth: &Auth{
			Status:    true,
			ProfileID: 101,
			Token:     packhub.Pointer(token),
			Password:  packhub.Pointer(password),
		},
	}

	result := user.ToMap()
	expected := &map[string]interface{}{
		"name":    "Alice",
		"mail":    "alice@example.com",
		"auth_id": uint(123),
		"Auth": map[string]interface{}{
			"status":     true,
			"profile_id": uint(101),
			"token":      packhub.Pointer(token),
			"password":   packhub.Pointer(password),
		},
	}

	assert.Equal(t, expected, result)
}

func TestUserBind(t *testing.T) {
	auth := &Auth{}
	user := User{Name: "Old Name", Email: "old@mail.com", Auth: auth}
	input := dto.UserInputDTO{
		Name:      packhub.Pointer("New Name"),
		Email:     packhub.Pointer("new@mail.com"),
		Status:    packhub.Pointer(true),
		ProfileID: packhub.Pointer(uint(101)),
	}

	err := user.Bind(&input)

	assert.NoError(t, err)
	assert.Equal(t, "New Name", user.Name)
	assert.Equal(t, "new@mail.com", user.Email)
	assert.Equal(t, true, auth.Status)
	assert.Equal(t, uint(101), auth.ProfileID)
}

func TestUserValidatePassword(t *testing.T) {
	password := "password123"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	auth := &Auth{Password: packhub.Pointer(string(hash))}
	user := User{Auth: auth}

	t.Run("CorrectPassword", func(t *testing.T) {
		assert.True(t, user.ValidatePassword(password))
	})

	t.Run("WrongPassword", func(t *testing.T) {
		assert.False(t, user.ValidatePassword("wrongpassword"))
	})
}

//
//func TestGenerateToken(t *testing.T) {
//	auth := &Auth{Token: utils.Pointer("token")}
//	user := User{Auth: auth}
//	privateKey := `-----BEGIN RSA PRIVATE KEY-----
//MIIBOgIBAAJBAK5nLGD+TdoG56+GGSZkLq6q3URhOaRUnDJbHt0b9g40kD1F2doY
//fjTaVkw0AY9Vo2v+gZzI+m/OEkl5FdHQOXMCAwEAAQJASdINCMyYDDAgEZ4N7sE/
//EP4DkFxsU9qOyCfiBY1UHeY9SkyEtq/9SGW91GZ/m0cRgQkP7Z4TUTxd/JoblIb/
//60ECIQD/KVgftDZf5YmJQS6lcErV09Cw9ZCDsPVjOZbT2sr12QIhAMV42drn5FRu
//4lh5CU3Rf5YOlqll41UQoF1XG8pVB5zdAiA3+07IP3lu8ZX1SKx9RgYyeXoy/ha3
//zkW+LMPa1rZD0QIhAIybHRRtHsJG/7tkFj3lwTAGanM+myhKhgDGN8Fg0zp1AiEA
//+TF3naExsUAnKWsTde3HUxZTkjvA3GtV95ben534WGI=
//-----END RSA PRIVATE KEY-----`
//
//	encodedKey := base64.StdEncoding.EncodeToString([]byte(privateKey))
//	expire := "10"
//
//	token, err := user.GenerateToken(expire, encodedKey)
//	fmt.Println("b err: ", err)
//
//	assert.NoError(t, err)
//	assert.NotEmpty(t, token)
//
//	decodedKey, _ := base64.StdEncoding.DecodeString(encodedKey)
//	key, _ := jwt.ParseRSAPrivateKeyFromPEM(decodedKey)
//
//	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
//		return key.PublicKey, nil
//	})
//
//	assert.NoError(t, err)
//	assert.NotNil(t, parsedToken)
//	//assert.True(t, parsedToken.Valid)
//
//	claims := parsedToken.Claims.(jwt.MapClaims)
//	assert.Equal(t, "token", claims["token"])
//	assert.Equal(t, "127.0.0.1", claims["ip"])
//	assert.NotEmpty(t, claims["iat"])
//	assert.NotEmpty(t, claims["exp"])
//}
