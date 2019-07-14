package middleware

import (
	"api/env"
	"api/model"
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

var (
	errInvalidJwtToken = errors.New("authentication: invalid token")
	errExpiredJwtToken = errors.New("authentication: expired token")
)

// Middleware :
type Middleware struct {
	mongo *mongo.Client
}

// New :
func New(db *mongo.Client) *Middleware {
	return &Middleware{mongo: db}
}

// JWTToken :
type JWTToken struct {
	token     string
	issuedAt  time.Time
	expiresAt time.Time
}

// Token :
func (jt JWTToken) Token() string {
	return jt.token
}

// ExpiresAt :
func (jt JWTToken) ExpiresAt() time.Time {
	return jt.expiresAt
}

// GenerateJWTToken :
func GenerateJWTToken(subject string, timeDuration ...time.Duration) (*JWTToken, error) {
	duration := time.Hour * 720
	if len(timeDuration) > 0 {
		duration = timeDuration[0]
	}

	issuedAt := time.Now().UTC()
	expiresAt := issuedAt.Add(duration)

	claims := jwt.Claims{
		Issuer:    os.Getenv("TAIZI_APP_NAME"),
		Subject:   subject,
		Expiry:    jwt.NewNumericDate(expiresAt),
		NotBefore: jwt.NewNumericDate(issuedAt),
		IssuedAt:  jwt.NewNumericDate(issuedAt),
	}

	sign, err := jose.NewSigner(
		jose.SigningKey{
			Algorithm: jose.HS256,
			Key:       []byte(os.Getenv("TAIZI_SECRET_KEY")),
		}, nil,
	)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Signed(sign).Claims(claims).CompactSerialize()
	if err != nil {
		return nil, err
	}

	return &JWTToken{
		token:     token,
		issuedAt:  issuedAt,
		expiresAt: expiresAt,
	}, nil
}

// ValidateJWTToken :
func ValidateJWTToken(token string) (*jwt.Claims, error) {
	return validateJWTToken(token)
}

func validateJWTToken(token string) (*jwt.Claims, error) {
	enc, err := jwt.ParseSigned(token)
	if err != nil {
		return nil, errInvalidJwtToken
	}

	claim := new(jwt.Claims)
	if err := enc.Claims([]byte(os.Getenv("TAIZI_SECRET_KEY")), claim); err != nil {
		return nil, errInvalidJwtToken
	}

	if claim.Expiry.Time().UTC().Unix()-time.Now().UTC().Unix() < 0 {
		return nil, errExpiredJwtToken
	}

	return claim, nil
}

// Authentication :
func (mw Middleware) Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		accessToken := c.Request().Header.Get("Authorization")
		accessToken = strings.TrimSpace(accessToken)
		accessToken = strings.Replace(accessToken, "Bearer ", "", 1)

		if accessToken == "" {
			return c.JSON(http.StatusExpectationFailed, errors.New("invalid access token"))
		}

		claim, err := validateJWTToken(accessToken)
		if err != nil {
			switch err {
			case errExpiredJwtToken:
				return c.JSON(http.StatusForbidden, err)
			case errInvalidJwtToken:
			}

			return c.JSON(http.StatusForbidden, err)
		}

		userID, err := primitive.ObjectIDFromHex(claim.Subject)
		if err != nil {
			return c.JSON(http.StatusForbidden, errors.New("invalid user id"))
		}

		user := model.User{}
		ctx := context.Background()
		coll := mw.mongo.Database(env.DefaultDB).
			Collection(model.CollectionUser)

		if err := coll.FindOne(ctx, bson.D{{"_id", userID}}).Decode(&user); err != nil {
			return c.JSON(http.StatusUnauthorized, "Invalid username or password")
		}

		// set user to the context
		c.Set(model.CollectionUser, user)

		return next(c)
	}
}
