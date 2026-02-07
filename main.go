package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	godotenv.Load()

	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	if err := db.Exec("USE sat_catalogs").Error; err != nil {
		log.Fatal("Failed to select database:", err)
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Configuración CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	config.AllowHeaders = []string{
		"Origin", "Content-Type", "Accept", "Authorization",
		"Cache-Control", "X-Requested-With", "Pragma",
	}
	r.Use(cors.New(config))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Migration endpoint
	r.POST("/api/migrate", migrateHandler)

	// Setup endpoint
	r.POST("/api/setup", setupHandler)

	// Reset endpoints
	r.POST("/api/reset", resetAllHandler)
	r.POST("/api/reset/:catalog", resetCatalogHandler)

	// Query endpoints
	r.GET("/api/cfdi/:catalog", getCatalog)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on :%s", port)
	r.Run(":" + port)
}

func migrateHandler(c *gin.Context) {
	addressCatalogs := []string{"estados", "municipios", "colonias", "codigos-postales", "localidades"}
	for _, catalog := range addressCatalogs {
		tableName := "cfdi_40_" + strings.Replace(catalog, "-", "_", -1)
		file := "./database/schemas/" + tableName + ".sql"
		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("Error reading %s: %v", file, err)})
			return
		}
		sql := strings.ReplaceAll(string(sqlBytes), "\"", "`")
		sql = strings.ReplaceAll(sql, " text", " TEXT")
		sql = normalizeSchemaSQL(sql)
		log.Printf("Executing schema for %s", filepath.Base(file))
		if err := db.Exec(sql).Error; err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("Error executing %s: %v", file, err)})
			return
		}
		log.Printf("Schema executed successfully for %s", filepath.Base(file))

		// Add indexes for optimization
		if tableName == "cfdi_40_estados" {
			if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_estados_pais ON cfdi_40_estados (pais);").Error; err != nil {
				log.Printf("Failed to create index for estados: %v", err)
			}
		} else if tableName == "cfdi_40_municipios" {
			if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_municipios_estado ON cfdi_40_municipios (estado);").Error; err != nil {
				log.Printf("Failed to create index for municipios: %v", err)
			}
		} else if tableName == "cfdi_40_colonias" {
			if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_colonias_codigo_postal ON cfdi_40_colonias (codigo_postal);").Error; err != nil {
				log.Printf("Failed to create index for colonias: %v", err)
			}
		} else if tableName == "cfdi_40_codigos_postales" {
			if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_codigos_postales_estado_municipio ON cfdi_40_codigos_postales (estado, municipio);").Error; err != nil {
				log.Printf("Failed to create index for codigos_postales: %v", err)
			}
		} else if tableName == "cfdi_40_localidades" {
			if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_localidades_estado ON cfdi_40_localidades (estado);").Error; err != nil {
				log.Printf("Failed to create index for localidades: %v", err)
			}
		}
	}

	c.JSON(200, gin.H{"message": "Migration completed"})
}

func setupHandler(c *gin.Context) {
	addressCatalogs := []string{"estados", "municipios", "colonias", "codigos-postales", "localidades"}
	totalFiles := len(addressCatalogs)
	for i, catalog := range addressCatalogs {
		tableName := "cfdi_40_" + strings.Replace(catalog, "-", "_", -1)
		file := "./database/data/" + tableName + ".sql"
		log.Printf("Setting up data from file %d/%d: %s", i+1, totalFiles, filepath.Base(file))
		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("Error reading %s: %v", file, err)})
			return
		}
		sql := strings.ReplaceAll(string(sqlBytes), "\"", "`")
		sql = strings.ReplaceAll(sql, " text", " TEXT")
		sql = strings.ReplaceAll(sql, "\r\n", "\n")
		statements := splitSQLStatements(sql)
		for _, stmt := range statements {
			clean := strings.TrimSpace(stmt)
			if clean == "" {
				continue
			}
			upper := strings.ToUpper(clean)
			if strings.HasPrefix(upper, "PRAGMA") || strings.HasPrefix(upper, "BEGIN TRANSACTION") || upper == "BEGIN" || strings.HasPrefix(upper, "COMMIT") {
				continue
			}
			if err := db.Exec(clean).Error; err != nil {
				c.JSON(500, gin.H{"error": fmt.Sprintf("Error executing %s: %v", file, err)})
				return
			}
		}
		log.Printf("Data inserted successfully for %s", filepath.Base(file))
	}
	c.JSON(200, gin.H{"message": "Data setup completed"})
}

func getCatalog(c *gin.Context) {
	catalog := c.Param("catalog")
	tableName := "cfdi_40_" + strings.Replace(catalog, "-", "_", -1)
	search := c.Query("search")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	offset := (page - 1) * limit

	log.Printf("Querying catalog %s, table %s, search %s", catalog, tableName, search)

	catalogKey := strings.Replace(catalog, "-", "_", -1)

	query := "SELECT * FROM " + tableName
	args := []interface{}{}
	conditions := []string{}
	joinClause := ""

	// Specific filters for faster queries
	if catalogKey == "estados" {
		if pais := c.Query("pais"); pais != "" {
			conditions = append(conditions, "pais = ?")
			args = append(args, pais)
		}
	} else if catalogKey == "municipios" {
		if estado := c.Query("estado"); estado != "" {
			conditions = append(conditions, "estado = ?")
			args = append(args, estado)
		}
	} else if catalogKey == "colonias" {
		joinClause = ""
		estadoColonias := c.Query("estado")
		municipioColonias := c.Query("municipio")
		if estadoColonias != "" || municipioColonias != "" {
			joinClause = " INNER JOIN cfdi_40_codigos_postales cp ON " + tableName + ".codigo_postal = cp.id"
		}
		if cp := c.Query("codigo_postal"); cp != "" {
			conditions = append(conditions, "codigo_postal = ?")
			args = append(args, cp)
		}
		if estadoColonias != "" {
			conditions = append(conditions, "cp.estado = ?")
			args = append(args, estadoColonias)
		}
		if municipioColonias != "" {
			conditions = append(conditions, "cp.municipio = ?")
			args = append(args, municipioColonias)
		}
	} else if catalogKey == "codigos_postales" {
		if estado := c.Query("estado"); estado != "" {
			conditions = append(conditions, "estado = ?")
			args = append(args, estado)
		}
		if municipio := c.Query("municipio"); municipio != "" {
			conditions = append(conditions, "municipio = ?")
			args = append(args, municipio)
		}
		if cp := c.Query("codigo_postal"); cp != "" {
			conditions = append(conditions, "id = ?")
			args = append(args, cp)
		}
		if localidad := c.Query("localidad"); localidad != "" {
			conditions = append(conditions, "localidad = ?")
			args = append(args, localidad)
		}
	} else if catalogKey == "localidades" {
		if estado := c.Query("estado"); estado != "" {
			conditions = append(conditions, "estado = ?")
			args = append(args, estado)
		}
	}

	// General search on texto if applicable
	if search != "" {
		conditions = append(conditions, "texto LIKE ?")
		args = append(args, "%"+search+"%")
	}

	query += joinClause

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	log.Printf("Executing query: %s with args: %v", query, args)

	// Count query for pagination
	countQuery := "SELECT COUNT(*) FROM " + tableName + joinClause
	conditionsArgs := args[:len(args)-2]
	if len(conditions) > 0 {
		countQuery += " WHERE " + strings.Join(conditions, " AND ")
	}
	var total int64
	err := db.Raw(countQuery, conditionsArgs...).Scan(&total).Error
	if err != nil {
		log.Printf("Count query failed: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	totalPages := (total + int64(limit) - 1) / int64(limit)

	var results []map[string]interface{}
	err = db.Raw(query, args...).Scan(&results).Error
	if err != nil {
		log.Printf("Query failed: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Query succeeded, returning %d results", len(results))
	c.JSON(200, gin.H{
		"data": results,
		"pagination": gin.H{
			"current_page": page,
			"per_page":     limit,
			"total":        total,
			"total_pages":  totalPages,
		},
	})
}

func splitSQLStatements(sql string) []string {
	var statements []string
	var current strings.Builder
	inSingleQuote := false
	sqlBytes := []byte(sql)
	for i := 0; i < len(sqlBytes); i++ {
		ch := sqlBytes[i]
		if ch == '\'' {
			if inSingleQuote {
				if i+1 < len(sqlBytes) && sqlBytes[i+1] == '\'' {
					current.WriteByte(ch)
					i++
					current.WriteByte(sqlBytes[i])
					continue
				}
				inSingleQuote = false
			} else {
				inSingleQuote = true
			}
			current.WriteByte(ch)
			continue
		}
		if ch == ';' && !inSingleQuote {
			stmt := strings.TrimSpace(current.String())
			if stmt != "" {
				statements = append(statements, stmt)
			}
			current.Reset()
			continue
		}
		current.WriteByte(ch)
	}
	if tail := strings.TrimSpace(current.String()); tail != "" {
		statements = append(statements, tail)
	}
	return statements
}

func normalizeSchemaSQL(sql string) string {
	rePK := regexp.MustCompile("(?i)PRIMARY\\s+KEY\\s*\\(([^)]+)\\)")
	matches := rePK.FindAllStringSubmatch(sql, -1)
	if len(matches) == 0 {
		return sql
	}
	primaryCols := make(map[string]struct{})
	for _, m := range matches {
		cols := strings.Split(m[1], ",")
		for _, col := range cols {
			name := strings.TrimSpace(col)
			name = strings.Trim(name, "`\"")
			if name == "" {
				continue
			}
			primaryCols[strings.ToLower(name)] = struct{}{}
		}
	}
	updated := sql
	for col := range primaryCols {
		reCol := regexp.MustCompile("(?i)(`" + regexp.QuoteMeta(col) + "`\\s+)TEXT")
		updated = reCol.ReplaceAllString(updated, "$1VARCHAR(255)")
	}
	return updated
}

func resetAllHandler(c *gin.Context) {
	catalogs := []string{"estados", "municipios", "colonias", "codigos-postales", "localidades"}
	for _, catalog := range catalogs {
		tableName := "cfdi_40_" + strings.Replace(catalog, "-", "_", -1)
		log.Printf("Truncating table %s", tableName)
		if err := db.Exec("TRUNCATE TABLE " + tableName).Error; err != nil {
			if err := db.Exec("DELETE FROM " + tableName).Error; err != nil {
				log.Printf("Error truncating %s: %v", tableName, err)
				c.JSON(500, gin.H{"error": fmt.Sprintf("Error resetting %s: %v", tableName, err)})
				return
			}
		}
	}
	c.JSON(200, gin.H{"message": "All tables reset successfully"})
}

func resetCatalogHandler(c *gin.Context) {
	catalog := c.Param("catalog")
	tableName := "cfdi_40_" + strings.Replace(catalog, "-", "_", -1)
	
	var count int64
	if err := db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_name = ?", tableName).Scan(&count).Error; err != nil {
		count = 0
	}
	if count == 0 {
		c.JSON(400, gin.H{"error": fmt.Sprintf("Table %s does not exist", tableName)})
		return
	}
	
	log.Printf("Truncating table %s", tableName)
	if err := db.Exec("TRUNCATE TABLE " + tableName).Error; err != nil {
		if err := db.Exec("DELETE FROM " + tableName).Error; err != nil {
			log.Printf("Error truncating %s: %v", tableName, err)
			c.JSON(500, gin.H{"error": fmt.Sprintf("Error resetting %s: %v", tableName, err)})
			return
		}
	}
	c.JSON(200, gin.H{"message": fmt.Sprintf("Table %s reset successfully", tableName)})
}
