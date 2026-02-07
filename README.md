# Microservicio de Catálogos SAT

Un microservicio para consultar catálogos del SAT (autoridad fiscal mexicana) usando Go, Gin, GORM y MariaDB.

## Características

- API REST para catálogos CFDI con búsqueda, paginación y filtros
- Endpoint de migración para crear tablas e índices
- Endpoint de configuración para cargar datos desde archivos SQL
- Soporte Docker con docker-compose
- Consultas optimizadas con índices de base de datos
- Endpoint de verificación de salud

## Prerrequisitos

- Docker y Docker Compose
- Go 1.25+ (para desarrollo local)

## Inicio Rápido

1. Clona el repositorio.

2. Copia `.env.example` a `.env` y ajusta las variables de entorno según sea necesario.

3. Construye y ejecuta los servicios:
   ```bash
   docker-compose up --build
   ```

4. Inicializa la base de datos:
   - POST a `http://localhost:8080/api/migrate` para crear tablas e índices.
   - POST a `http://localhost:8080/api/setup` para cargar los datos de los catálogos.

5. Consulta catálogos:
   - GET `http://localhost:8080/api/cfdi/paises?search=mexico&page=1&limit=10`

## Endpoints de la API

### Verificación de Salud
- `GET /health` - Verificación de salud del servicio.

### Migración y Configuración
- `POST /api/migrate` - Ejecutar migraciones de base de datos (crear tablas e índices para direcciones).
- `POST /api/setup` - Cargar datos desde archivos SQL para direcciones.

### Consulta de Catálogos de Direcciones
- `GET /api/cfdi/{catalog}` - Consultar un catálogo CFDI específico.
  - Catálogos disponibles: `estados`, `municipios`, `colonias`, `codigos-postales`.
  - Parámetros de consulta:
    - `page` (entero, defecto: 1): Página actual.
    - `limit` (entero, defecto: 10): Elementos por página.
    - `search` (cadena): Búsqueda general en el campo "texto".
    - Filtros específicos:
      - **estados**: `pais` (ej. "MEX").
      - **municipios**: `estado` (ej. "JAL").
      - **colonias**: `codigo_postal`, `estado`, `municipio`.
      - **codigos-postales**: `estado`, `municipio`, `codigo_postal`.

#### Ejemplos de Consultas
- Cargar estados: `GET /api/cfdi/estados`
- Municipios de Jalisco: `GET /api/cfdi/municipios?estado=JAL`
- Colonias de Guadalajara: `GET /api/cfdi/colonias?estado=JAL&municipio=039`
- Códigos postales: `GET /api/cfdi/codigos-postales?estado=JAL&municipio=039`

#### Respuesta
```json
{
  "data": [
    {
      "estado": "JAL",
      "texto": "Jalisco"
    }
  ],
  "pagination": {
    "current_page": 1,
    "per_page": 10,
    "total": 32,
    "total_pages": 4
  }
}
```

Otros catálogos disponibles: paises, productos_servicios, formas_pago, etc. (consulta el código para más detalles).

## Variables de Entorno

Copia `.env.example` a `.env` y configura:

- `MYSQL_DATABASE`: Nombre de la base de datos
- `MYSQL_USER`: Usuario de la base de datos
- `MYSQL_PASSWORD`: Contraseña de la base de datos
- `MYSQL_ROOT_PASSWORD`: Contraseña root para MariaDB
- `DB_HOST`: Host para la app (localhost para local, db para Docker)
- `DB_PORT`: Puerto (3306)
- `AUTO_MIGRATE`: Establece en true para ejecutar migraciones al inicio (opcional)

## Desarrollo

Para desarrollo local sin Docker:

1. Instala MariaDB y crea una base de datos.

2. Establece `DB_HOST=localhost` en `.env`.

3. Ejecuta `go mod tidy` y `go run main.go`.

## Estructura del Proyecto

- `main.go`: Configuración del servidor y manejadores
- `models.go`: Modelos GORM para catálogos
- `docker-compose.yml`: Servicios Docker
- `Dockerfile`: Construcción del contenedor de la app
- `database/`: Archivos SQL para esquemas y datos

## Licencia

MIT
