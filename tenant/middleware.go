package tenant

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type contextKey string

const TenantKey contextKey = "tenant"

func Middleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		slug := c.GetHeader("X-Tenant")
		if slug == "" {
			slug = resolveSubDomain(c)
		}
		if slug == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "header X-Tenant o subdominio requerido"})
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
	slug, ok := c.Get(string(TenantKey))
	if !ok {
		return "", false
	}
	return slug.(string), true
}
