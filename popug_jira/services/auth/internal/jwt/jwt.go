package jwt

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ExtendedClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	PublicID string `json:"id"`
	Role     string `json:"role"`
}

func New(jwtExpired time.Duration) *JWTManager {
	return &JWTManager{
		expiredToken: jwtExpired,
		secretKey:    []byte("09876543210"),
	}
}

type JWTManager struct {
	expiredToken time.Duration
	secretKey    []byte
}

func (j *JWTManager) GenerateToken(username, id, role string) (string, error) {
	expirationTime := time.Now().Add(j.expiredToken)

	claims := &ExtendedClaims{
		Username: username,
		PublicID: id,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		log.Errorf("token.SignedString(...): %s", err.Error())

		return "", fmt.Errorf("jwt.NewWithClaims(...): %w", err)
	}

	return tokenString, nil
}

func (j *JWTManager) Authorize(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	authHeader := md["authorization"]
	if len(authHeader) == 0 {
		return "", status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	// The authHeader should be in the format `Bearer <token>`
	splitToken := strings.Split(authHeader[0], " ")
	if len(splitToken) != 2 || strings.ToLower(splitToken[0]) != "bearer" {
		return "", status.Errorf(codes.Unauthenticated, "authorization token format is invalid")
	}

	tokenStr := splitToken[1]

	claims := &ExtendedClaims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", fmt.Errorf("invalid token signature: %w", err)
		}

		return "", fmt.Errorf("invalid token: %w", err)
	}

	if !tkn.Valid {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	return claims.Role, nil
}
