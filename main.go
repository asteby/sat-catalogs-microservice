package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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
		if err := db.Exec(string(sqlBytes)).Error; err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("Error executing %s: %v", file, err)})
			return
		}
	}

	// Create indexes for query optimization
	indexSQLs := []string{
		"CREATE INDEX IF NOT EXISTS idx_cfdi_40_paises_texto ON cfdi_40_paises (texto);",
		"CREATE INDEX IF NOT EXISTS idx_cfdi_40_estados_pais ON cfdi_40_estados (pais);",
		"CREATE INDEX IF NOT EXISTS idx_cfdi_40_estados_texto ON cfdi_40_estados (texto);",
		"CREATE INDEX IF NOT EXISTS idx_cfdi_40_municipios_estado ON cfdi_40_municipios (estado);",
		"CREATE INDEX IF NOT EXISTS idx_cfdi_40_municipios_texto ON cfdi_40_municipios (texto);",
		"CREATE INDEX IF NOT EXISTS idx_cfdi_40_colonias_codigo_postal ON cfdi_40_colonias (codigo_postal);",
		"CREATE INDEX IF NOT EXISTS idx_cfdi_40_colonias_estado ON cfdi_40_colonias (estado);",
		"CREATE INDEX IF NOT EXISTS idx_cfdi_40_colonias_texto ON cfdi_40_colonias (texto);",
		"CREATE INDEX IF NOT EXISTS idx_cfdi_40_productos_servicios_texto ON cfdi_40_productos_servicios (texto);",
		"CREATE INDEX IF NOT EXISTS idx_cfdi_40_formas_pago_texto ON cfdi_40_formas_pago (texto);",
		"CREATE INDEX IF NOT EXISTS idx_cfdi_40_monedas_texto ON cfdi_40_monedas (texto);",
		"CREATE INDEX IF NOT EXISTS idx_cfdi_40_usos_cfdi_texto ON cfdi_40_usos_cfdi (texto);",
		"CREATE INDEX IF NOT EXISTS idx_cfdi_40_regimenes_fiscales_texto ON cfdi_40_regimenes_fiscales (texto);",
		"CREATE INDEX IF NOT EXISTS idx_cfdi_40_tipos_comprobantes_texto ON cfdi_40_tipos_comprobantes (texto);",
		"CREATE INDEX IF NOT EXISTS idx_cfdi_40_metodos_pago_texto ON cfdi_40_metodos_pago (texto);",
		"CREATE INDEX IF NOT EXISTS idx_cfdi_40_codigos_postales_estado ON cfdi_40_codigos_postales (d_estado);",
		"CREATE INDEX IF NOT EXISTS idx_cfdi_40_codigos_postales_municipio ON cfdi_40_codigos_postales (D_mnpio);",
		"CREATE INDEX IF NOT EXISTS idx_cfdi_40_codigos_postales_cp ON cfdi_40_codigos_postales (d_codigo);",
	}
	for _, sql := range indexSQLs {
		if err := db.Exec(sql).Error; err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("Error creating index: %v", err)})
			return
		}
	}

	c.JSON(200, gin.H{"message": "Migration and index creation completed"})
}

func setupHandler(c *gin.Context) {
	files, err := filepath.Glob("./database/data/*.sql")
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
		if err := db.Exec(string(sqlBytes)).Error; err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("Error executing %s: %v", file, err)})
			return
		}
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

	var results []map[string]interface{}
	err := db.Raw(query, args...).Scan(&results).Error
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, results)
}
