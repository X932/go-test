package authmiddleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"test-go/internal/response"
	user_service "test-go/internal/services/user"
	"test-go/pkg/constants"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	ID int `json:"id"`
}

const UNAUTHORIZED_MESSAGE = "Unauthorized"
const BEARER_PREFIX = "Bearer"
const HEADER_JWT_SLICE_COUNT = 2

func Authmiddleware(userService user_service.Service, JWT_SECRET string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		headerJwt := ctx.GetHeader("authorization")
		fmt.Println("============== Authmiddleware ==============")
		if headerJwt == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Body{Message: UNAUTHORIZED_MESSAGE})
			return
		}

		splittedHeaderJwt := strings.Fields(headerJwt)

		if len(splittedHeaderJwt) != HEADER_JWT_SLICE_COUNT || splittedHeaderJwt[0] != BEARER_PREFIX {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Body{Message: UNAUTHORIZED_MESSAGE})
			return
		}

		hmacSecret := []byte(JWT_SECRET)
		jwtFromHeader, jwtParsingErr := jwt.Parse(splittedHeaderJwt[1], func(t *jwt.Token) (any, error) {
			return hmacSecret, nil
		})

		if jwtParsingErr != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Body{Message: jwtParsingErr.Error()})
			return
		}

		if jwtFromHeader.Valid == false {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Body{Message: "invalid jwt token"})
			return
		}

		claims, ok := jwtFromHeader.Claims.(jwt.MapClaims)

		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Body{Message: "error parsing claims"})
			return
		}

		subject, subjectErr := claims.GetSubject()

		if subjectErr != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Body{Message: subjectErr.Error()})
			return
		}

		userID, convErr := strconv.Atoi(subject)

		if convErr != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Body{Message: convErr.Error()})
			return
		}

		user, userServiceErr := userService.GetUser(userID)

		if userServiceErr != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Body{Message: userServiceErr.Error()})
			return
		}

		ctx.Set(constants.USER_CTX_KEY, user)

		ctx.Next()
	}
}
