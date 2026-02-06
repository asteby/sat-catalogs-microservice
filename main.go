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
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
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

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Migration endpoint
	r.POST("/api/migrate", migrateHandler)

	// Setup endpoint
	r.POST("/api/setup", setupHandler)

	// Query endpoints
	r.GET("/api/cfdi/:catalog", getCatalog)
	// Add more versions or modules as needed

	log.Println("Starting server on :8080")
	r.Run(":8080")
}

func migrateHandler(c *gin.Context) {
	files, err := filepath.Glob("./database/schemas/*.sql")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	for _, file := range files {
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
	}

	c.JSON(200, gin.H{"message": "Migration completed"})
}

func setupHandler(c *gin.Context) {
	files, err := filepath.Glob("./database/data/*.sql")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	totalFiles := len(files)
	for i, file := range files {
		log.Printf("Setting up data from file %d/%d: %s", i+1, totalFiles, filepath.Base(file))
		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("Error reading %s: %v", file, err)})
			return
		}
		sql := strings.ReplaceAll(string(sqlBytes), "\"", "`")
		sql = strings.ReplaceAll(sql, " text", " TEXT")
		sql = strings.ReplaceAll(sql, "\\r\\n", "\\n")
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

	query := "SELECT * FROM " + tableName
	args := []interface{}{}
	conditions := []string{}

	// Specific filters for faster queries
	if catalog == "estados" {
		if pais := c.Query("pais"); pais != "" {
			conditions = append(conditions, "pais = ?")
			args = append(args, pais)
		}
	} else if catalog == "municipios" {
		if estado := c.Query("estado"); estado != "" {
			conditions = append(conditions, "estado = ?")
			args = append(args, estado)
		}
	} else if catalog == "colonias" {
		if cp := c.Query("codigo_postal"); cp != "" {
			conditions = append(conditions, "codigo_postal = ?")
			args = append(args, cp)
		}
		if estado := c.Query("estado"); estado != "" {
			conditions = append(conditions, "estado = ?")
			args = append(args, estado)
		}
	} else if catalog == "codigos_postales" {
		if estado := c.Query("estado"); estado != "" {
			conditions = append(conditions, "d_estado = ?")
			args = append(args, estado)
		}
		if municipio := c.Query("municipio"); municipio != "" {
			conditions = append(conditions, "D_mnpio = ?")
			args = append(args, municipio)
		}
		if cp := c.Query("codigo_postal"); cp != "" {
			conditions = append(conditions, "d_codigo = ?")
			args = append(args, cp)
		}
	}

	// General search on texto if applicable
	if search != "" {
		conditions = append(conditions, "texto LIKE ?")
		args = append(args, "%"+search+"%")
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	log.Printf("Executing query: %s with args: %v", query, args)

	var results []map[string]interface{}
	err := db.Raw(query, args...).Scan(&results).Error
	if err != nil {
		log.Printf("Query failed: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Query succeeded, returning %d results", len(results))
	c.JSON(200, results)
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
