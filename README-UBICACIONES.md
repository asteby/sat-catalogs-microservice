# API de Ubicaciones SAT

## Introducción

Esta API proporciona acceso a los catálogos de direcciones del SAT (estados, municipios, colonias, códigos postales) con búsqueda, filtros, paginación y consultas en cascada para integraciones web.

## Tablas y Campos

Basado en los esquemas SQL:

### Estados
- **estado** (VARCHAR, PK): Código del estado (ej. "JAL")
- **pais** (VARCHAR): País (ej. "MEX")
- **texto** (TEXT): Nombre del estado
- **vigencia_desde** (TEXT): Fecha de inicio de vigencia
- **vigencia_hasta** (TEXT): Fecha de fin de vigencia

### Municipios
- **municipio** (VARCHAR, PK): Código del municipio
- **estado** (VARCHAR, PK): Código del estado
- **texto** (TEXT): Nombre del municipio
- **vigencia_desde** (TEXT): Fecha de inicio
- **vigencia_hasta** (TEXT): Fecha de fin

### Colonias
- **colonia** (VARCHAR, PK): Nombre de la colonia
- **codigo_postal** (VARCHAR, PK): Código postal
- **texto** (TEXT): Descripción

### Códigos Postales
- **id** (VARCHAR, PK): Código postal
- **estado** (TEXT): Código del estado
- **municipio** (TEXT): Código del municipio
- **localidad** (TEXT): Localidad
- **estimulo_frontera** (INT): Indicador de estímulo fronterizo
- **vigencia_desde** (TEXT): Fecha de inicio
- **vigencia_hasta** (TEXT): Fecha de fin
- **huso_descripcion** (TEXT): Descripción del huso horario
- **huso_verano_mes_inicio** (TEXT): Mes inicio verano
- **huso_verano_dia_inicio** (TEXT): Día inicio verano
- **huso_verano_hora_inicio** (TEXT): Hora inicio verano
- **huso_verano_diferencia** (TEXT): Diferencia verano
- **huso_invierno_mes_inicio** (TEXT): Mes inicio invierno
- **huso_invierno_dia_inicio** (TEXT): Día inicio invierno
- **huso_invierno_hora_inicio** (TEXT): Hora inicio invierno
- **huso_invierno_diferencia** (TEXT): Diferencia invierno

## Endpoints

### GET /api/cfdi/{catalog}

Consulta un catálogo específico. Catálogos soportados: `estados`, `municipios`, `colonias`, `codigos-postales`.

#### Parámetros de Consulta
- `page` (entero, default: 1): Número de página para paginación.
- `limit` (entero, default: 10): Número de elementos por página (máx. recomendado: 100).
- `search` (cadena): Búsqueda general en el campo "texto" (coincidencia parcial, insensible a mayúsculas).
- Filtros específicos por catálogo (combinados con AND):
  - **estados**: `pais` (ej. "MEX" para México).
  - **municipios**: `estado` (ej. "JAL" para Jalisco).
  - **colonias**: `codigo_postal` (ej. "44100"), `estado` (ej. "JAL"), `municipio` (ej. "039" para Guadalajara).
  - **codigos-postales**: `estado`, `municipio`, `codigo_postal`.

#### Respuesta
```json
{
  "data": [
    {
      "estado": "JAL",
      "pais": "MEX",
      "texto": "Jalisco",
      "vigencia_desde": "2017-01-01",
      "vigencia_hasta": ""
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

## Ejemplos de Consultas

### Carga en Cascada (para Dropdowns Web)
1. **Cargar Estados**: `GET /api/cfdi/estados`
   - Respuesta: Lista de estados con códigos y textos.

2. **Cargar Municipios de un Estado**: `GET /api/cfdi/municipios?estado=JAL`
   - Filtra municipios de Jalisco.

3. **Cargar Colonias de un Municipio**: `GET /api/cfdi/colonias?estado=JAL&municipio=039`
   - Usa join con códigos postales para filtrar colonias de Guadalajara, Jalisco.

4. **Cargar Códigos Postales**: `GET /api/cfdi/codigos-postales?estado=JAL&municipio=039`
   - Códigos postales del mismo municipio.

### Búsqueda y Filtros
- **Buscar Estados por País**: `GET /api/cfdi/estados?pais=MEX&search=jal`
- **Buscar Colonias con Texto**: `GET /api/cfdi/colonias?search=centro&estado=JAL`
- **Filtrar por Código Postal**: `GET /api/cfdi/colonias?codigo_postal=44100`

### Paginación
- **Página Específica**: `GET /api/cfdi/estados?page=2&limit=5`
  - Página 2, 5 elementos por página.

## Notas Técnicas
- **Índices**: Las tablas tienen índices en campos de filtro para rendimiento (pais, estado, codigo_postal, etc.).
- **Joins**: Colonias usa join con códigos_postales para filtros avanzados.
- **Campos de Respuesta**: Incluyen todos los campos de la tabla; "texto" es el nombre descriptivo.
- **Errores**: 500 con `{"error": "mensaje"}` si falla la consulta.
- **Vigencia**: Campos vigencia_desde/hasta indican validez; vacío significa actual.
- **Huso Horario**: Solo en códigos_postales, para zonas horarias mexicanas.

## Uso en Aplicaciones Web
- Implementa dropdowns en cascada: Estado → Municipio → Colonia.
- Usa search para autocompletado.
- Pagina para listas grandes.
- Campos clave: estado, municipio, colonia, codigo_postal para integraciones.
