package robokassa

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(password string, data JWT) (string, error) {

	conc := fmt.Sprintf("%s:%s", data.MerchantLogin, password)

	secret := base64.StdEncoding.EncodeToString([]byte(conc))

	// Marshal and Unmarshal (struct -> json -> map[string]interface{})
	// So im think it's good practice for save abstraction
	var dat map[string]interface{}

	pay, err := json.Marshal(data)
	if err != nil {
		return "", err // temp
	}

	json.Unmarshal(pay, &dat)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(dat))

	signToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err // temp
	}

	return signToken, nil
}
