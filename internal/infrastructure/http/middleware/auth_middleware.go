package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/labstack/echo/v4"
)

// CustomClaims contains custom data we want from the token.
type CustomClaims struct {
	Scope string `json:"scope"`
}

// Validate satisfies validator.CustomClaims interface.
func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

// AuthWithScopes ensures the JWT is valid and contains the required scopes.
func AuthWithScopes(requiredScopes ...string) echo.MiddlewareFunc {
	issuerURLString := fmt.Sprintf("https://%s/", os.Getenv("AUTH0_DOMAIN"))

	// URLをパース
	issuerURL, err := url.Parse(issuerURLString)
	if err != nil {
		log.Fatalf("Failed to parse issuer URL: %v", err)
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{os.Getenv("AUTH0_AUDIENCE")},
		validator.WithCustomClaims(func() validator.CustomClaims {
			return &CustomClaims{}
		}),
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		log.Fatalf("Failed to set up the jwt validator: %v", err)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Authorization ヘッダーからトークンを取得
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization header")
			}

			// トークンを抽出（"Bearer " プレフィックスを除去）
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == authHeader {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Authorization header format")
			}

			// Validate JWT
			validatedToken, err := jwtValidator.ValidateToken(c.Request().Context(), token)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Failed to validate token")
			}

			// Extract claims
			claims, ok := validatedToken.(*validator.ValidatedClaims).CustomClaims.(*CustomClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token claims")
			}

			// スコープのチェック
			if !hasRequiredScopes(claims.Scope, requiredScopes) {
				return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("Insufficient scope: required %v", requiredScopes))
			}

			// Store claims in context for downstream use
			c.Set("user", claims)

			// Proceed to the next handler
			return next(c)
		}
	}
}

// hasRequiredScopes checks if the token contains all required scopes.
func hasRequiredScopes(tokenScopes string, requiredScopes []string) bool {
	tokenScopeList := strings.Split(tokenScopes, " ")
	scopeSet := make(map[string]struct{})
	for _, scope := range tokenScopeList {
		scopeSet[scope] = struct{}{}
	}

	for _, requiredScope := range requiredScopes {
		if _, found := scopeSet[requiredScope]; !found {
			return false
		}
	}

	return true
}
