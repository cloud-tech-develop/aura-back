package tenant

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/cloud-tech-develop/aura-back/shared/domain/vo"
	"github.com/cloud-tech-develop/aura-back/shared/errors"
	"github.com/cloud-tech-develop/aura-back/shared/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserID       int64    `json:"user_id"`
	EnterpriseID int64    `json:"enterprise_id"`
	TenantID     int64    `json:"tenant_id"` // Agregado: ID del tenant
	Slug         string   `json:"slug"`
	Email        string   `json:"email"`
	Roles        []string `json:"roles"`
	RoleLevel    int      `json:"role_level"` // Nivel del rol más alto del usuario (0=superadmin, 1=admin, etc.)
	IP           string   `json:"ip"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte(getEnv("JWT_SECRET", "aura-secret-key-change-in-production"))

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(userID int64, e *Enterprise, roles []string, roleLevel int, ip string) (string, error) {
	claims := Claims{
		UserID:       userID,
		EnterpriseID: e.ID,
		TenantID:     e.TenantID, // Include tenant ID
		Slug:         e.Slug,
		Email:        e.Email.String(),
		Roles:        roles,
		RoleLevel:    roleLevel,
		IP:           ip,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// GenerateTokenOffline generates a JWT token for offline mode (simpler)
func GenerateTokenOffline(userID int64, enterpriseID int64, email string) (string, error) {
	claims := Claims{
		UserID:       userID,
		EnterpriseID: enterpriseID,
		Slug:         "offline",
		Email:        email,
		Roles:        []string{"ADMIN"},
		RoleLevel:    100,
		IP:           "127.0.0.1",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func getClientIP(c *gin.Context) string {
	ip := c.GetHeader("X-Forwarded-For")
	if ip != "" {
		ips := net.ParseIP(ip)
		if ips != nil {
			return ips.String()
		}
	}
	ip = c.GetHeader("X-Real-IP")
	if ip != "" {
		return ip
	}
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}
	return ip
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"data": nil, "success": false, "message": errors.ErrAuthHeaderRequired.Error()})
			return
		}

		tokenString := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"data": nil, "success": false, "message": errors.ErrTokenInvalid.Error()})
			return
		}

		clientIP := getClientIP(c)
		if claims.IP != "" && claims.IP != clientIP {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"data": nil, "success": false, "message": errors.ErrIPNotValid.Error()})
			return
		}

		c.Set(string(TenantKey), claims.Slug)
		c.Set("user_id", claims.UserID)
		c.Set("enterprise_id", claims.EnterpriseID)
		c.Set("email", claims.Email)
		c.Set("roles", claims.Roles)
		c.Set("role_level", claims.RoleLevel)
		c.Next()
	}
}

func Login(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email    vo.Email `json:"email" binding:"required"`
			Password string   `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, errors.ErrEmailPasswordRequired.Error())
			return
		}

		email := req.Email.String()
		password := req.Password

		// 1. Authenticate user and get enterprise info
		var user struct {
			ID           int64
			PasswordHash string
			EnterpriseID int64
			TenantID     int64
			Slug         string
			EntStatus    string
			UserActive   bool
		}

		query := `
			SELECT u.id, u.password_hash, e.id, e.tenant_id, e.slug, e.status, u.active 
			FROM public.users u
			JOIN public.enterprises e ON u.enterprise_id = e.id
			WHERE u.email = $1 
			AND u.deleted_at IS NULL 
			AND e.deleted_at IS NULL`

		err := db.QueryRowContext(c.Request.Context(), query, email).Scan(
			&user.ID, &user.PasswordHash, &user.EnterpriseID, &user.TenantID, &user.Slug, &user.EntStatus, &user.UserActive,
		)

		if err == sql.ErrNoRows {
			response.Unauthorized(c, errors.ErrInvalidCredentials.Error())
			return
		}
		if err != nil {
			response.InternalServerError(c, err.Error())
			return
		}

		if user.EntStatus != "ACTIVE" {
			response.Forbidden(c, errors.ErrEnterpriseInactive.Error())
			return
		}

		if !user.UserActive {
			response.Forbidden(c, errors.ErrUserInactive.Error())
			return
		}

		if !CheckPassword(password, user.PasswordHash) {
			response.Unauthorized(c, errors.ErrInvalidCredentials.Error())
			return
		}

		// 2. Fetch roles from public schema and get the highest privilege level (minimum level number)
		roles := []string{}
		roleLevel := 100 // Default high level (no privilege)
		roleRows, err := db.QueryContext(c.Request.Context(), `
			SELECT r.name, r.level 
			FROM public.roles r
			JOIN public.user_roles ur ON ur.role_id = r.id
			WHERE ur.user_id = $1
			AND r.deleted_at IS NULL`, user.ID)
		if err == nil {
			defer roleRows.Close()
			for roleRows.Next() {
				var roleName string
				var level int
				if err := roleRows.Scan(&roleName, &level); err == nil {
					roles = append(roles, roleName)
					if level < roleLevel {
						roleLevel = level
					}
				}
			}
		}

		// 3. Fetch third party info from tenant schema
		var thirdParty struct {
			FirstName *string
			LastName  *string
		}
		tpQuery := fmt.Sprintf(`
			SELECT tp.first_name, tp.last_name 
			FROM %q.third_parties tp
			WHERE tp.user_id = $1 AND tp.deleted_at IS NULL
			LIMIT 1`, user.Slug)
		_ = db.QueryRowContext(c.Request.Context(), tpQuery, user.ID).Scan(
			&thirdParty.FirstName, &thirdParty.LastName,
		)

		// 4. Generate JWT
		ent := Enterprise{
			ID:       user.EnterpriseID,
			TenantID: user.TenantID,
			Slug:     user.Slug,
			Email:    vo.Email(email),
		}

		clientIP := getClientIP(c)
		token, err := GenerateToken(user.ID, &ent, roles, roleLevel, clientIP)
		if err != nil {
			response.InternalServerError(c, errors.ErrTokenGeneration.Error())
			return
		}

		response.OK(c, gin.H{
			"token": token,
			"user": gin.H{
				"id":         user.ID,
				"email":      email,
				"first_name": thirdParty.FirstName,
				"last_name":  thirdParty.LastName,
				"roles":      roles,
			},
			"enterprise": gin.H{
				"id":        user.EnterpriseID,
				"tenant_id": user.TenantID,
				"slug":      user.Slug,
			},
		})
	}
}
