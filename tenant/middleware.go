package tenant

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const TenantKey contextKey = "tenant"

// TokenClaims represents all claims from the JWT token
type TokenClaims struct {
	UserID        int64    `json:"user_id"`
	EnterpriseID  int64    `json:"enterprise_id"`
	TenantID      int64    `json:"tenant_id"`
	Slug         string  `json:"slug"`
	Email        string  `json:"email"`
	Roles        []string `json:"roles"`
	RoleLevel    int      `json:"role_level"`
	IP           string  `json:"ip"`
}

func Middleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. First, try to get tenant from JWT token (set by AuthMiddleware)
		slug, hasSlug := SlugFromContext(c)

		// 2. If not in token, try to resolve from subdomain
		if !hasSlug || slug == "" {
			slug = resolveSubDomain(c)
		}

		// 3. If still no tenant, require authentication
		if slug == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "tenant no encontrado. Inicie sesión para acceder a un tenant"})
			return
		}

		var exists bool
		db.QueryRowContext(c.Request.Context(),
			`SELECT EXISTS(SELECT 1 FROM public.enterprises WHERE slug=$1 AND status='ACTIVE')`, slug,
		).Scan(&exists)

		if !exists {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "enterprise no encontrado o inactivo"})
			return
		}

		c.Set(string(TenantKey), slug)
		c.Next()
	}
}

func resolveSubDomain(c *gin.Context) string {
	host := c.Request.Host
	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}
	parts := strings.Split(host, ".")
	if len(parts) >= 3 {
		subDomain := parts[0]
		if subDomain != "www" && subDomain != "api" {
			return subDomain
		}
	}
	return ""
}

func SubDomainMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		subDomain := resolveSubDomain(c)
		if subDomain == "" {
			c.Next()
			return
		}

		var slug string
		err := db.QueryRowContext(c.Request.Context(),
			`SELECT slug FROM public.enterprises WHERE sub_domain = $1 AND status = 'ACTIVE'`,
			subDomain,
		).Scan(&slug)

		if err == sql.ErrNoRows {
			c.Next()
			return
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error al resolver subdominio"})
			return
		}

		c.Set(string(TenantKey), slug)
		c.Next()
	}
}

func SlugFromContext(c *gin.Context) (string, bool) {
	// 1. First, try to get from context (set by AuthMiddleware)
	slug, ok := c.Get(string(TenantKey))
	if ok && slug.(string) != "" {
		return slug.(string), true
	}

	// 2. Try to get from JWT token directly
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		tokenString := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err == nil && token.Valid && claims.Slug != "" {
			return claims.Slug, true
		}
	}

	return "", false
}

func EnterpriseIDFromContext(c *gin.Context) (int64, bool) {
	// 1. First, try to get from context (set by AuthMiddleware)
	enterpriseID, ok := c.Get("enterprise_id")
	if ok && enterpriseID.(int64) > 0 {
		return enterpriseID.(int64), true
	}

	// 2. Try to get from JWT token directly
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		tokenString := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err == nil && token.Valid && claims.EnterpriseID > 0 {
			return claims.EnterpriseID, true
		}
	}

	return 0, false
}

// ClaimsFromContext returns all token claims from context or JWT token
func ClaimsFromContext(c *gin.Context) (*TokenClaims, bool) {
	// 1. First, try to get from context (set by AuthMiddleware)
	userID, ok := c.Get("user_id")
	if ok && userID.(int64) > 0 {
		enterpriseID, _ := c.Get("enterprise_id")
		email, _ := c.Get("email")
		roles, _ := c.Get("roles")
		roleLevel, _ := c.Get("role_level")
		slug, _ := c.Get(string(TenantKey))

		return &TokenClaims{
			UserID:       userID.(int64),
			EnterpriseID: enterpriseID.(int64),
			Slug:         slug.(string),
			Email:        email.(string),
			Roles:        roles.([]string),
			RoleLevel:    roleLevel.(int),
		}, true
	}

	// 2. Try to get from JWT token directly
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		tokenString := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err == nil && token.Valid {
			return &TokenClaims{
				UserID:       claims.UserID,
				EnterpriseID: claims.EnterpriseID,
				TenantID:     claims.TenantID,
				Slug:         claims.Slug,
				Email:        claims.Email,
				Roles:        claims.Roles,
				RoleLevel:    claims.RoleLevel,
				IP:           claims.IP,
			}, true
		}
	}

	return nil, false
}
