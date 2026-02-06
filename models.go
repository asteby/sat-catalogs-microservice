package main

// Pais model for cfdi_40_paises
type Pais struct {
	ID                            string `gorm:"primaryKey"`
	Texto                         string
	PatronCodigoPostal            string
	PatronIdentidadTributaria     string
	ValidacionIdentidadTributaria string
	Agrupaciones                  string
}

func (Pais) TableName() string {
	return "cfdi_40_paises"
}

// ProductoServicio model for cfdi_40_productos_servicios
type ProductoServicio struct {
	ID               string `gorm:"primaryKey"`
	Texto            string
	IvaTrasladado    int
	IepsTrasladado   int
	Complemento      string
	VigenciaDesde    string
	VigenciaHasta    string
	EstimuloFrontera int
	Similares        string
}

func (ProductoServicio) TableName() string {
	return "cfdi_40_productos_servicios"
}

// Estado model for cfdi_40_estados
type Estado struct {
	Estado        string `gorm:"primaryKey"`
	Pais          string `gorm:"primaryKey"`
	Texto         string
	VigenciaDesde string
	VigenciaHasta string
}

func (Estado) TableName() string {
	return "cfdi_40_estados"
}

// Colonia model for cfdi_40_colonias
type Colonia struct {
	Colonia      string `gorm:"primaryKey"`
	CodigoPostal string `gorm:"primaryKey"`
	Texto        string
}

func (Colonia) TableName() string {
	return "cfdi_40_colonias"
}

// Municipio model for cfdi_40_municipios
type Municipio struct {
	Municipio     string `gorm:"primaryKey"`
	Estado        string `gorm:"primaryKey"`
	Texto         string
	VigenciaDesde string
	VigenciaHasta string
}

func (Municipio) TableName() string {
	return "cfdi_40_municipios"
}

// FormasPago model for cfdi_40_formas_pago
type FormasPago struct {
	ID                              string `gorm:"primaryKey"`
	Texto                           string
	EsBancarizado                   int
	RequiereNumeroOperacion         int
	PermiteBancoOrdenanteRFC        int
	PermiteCuentaOrdenante          int
	PatronCuentaOrdenante           string
	PermiteBancoBeneficiarioRFC     int
	PermiteCuentaBeneficiario       int
	PatronCuentaBeneficiario        string
	PermiteTipoCadenaPago           int
	RequiereBancoOrdenanteNombreExt int
	VigenciaDesde                   string
	VigenciaHasta                   string
}

func (FormasPago) TableName() string {
	return "cfdi_40_formas_pago"
}

// Monedas model for cfdi_40_monedas
type Monedas struct {
	ID                  string `gorm:"primaryKey"`
	Texto               string
	Decimales           int
	PorcentajeVariacion int
	VigenciaDesde       string
	VigenciaHasta       string
}

func (Monedas) TableName() string {
	return "cfdi_40_monedas"
}

// UsosCFDI model for cfdi_40_usos_cfdi
type UsosCFDI struct {
	ID                          string `gorm:"primaryKey"`
	Texto                       string
	AplicaFisica                int
	AplicaMoral                 int
	VigenciaDesde               string
	VigenciaHasta               string
	RegimenesFiscalesReceptores string
}

func (UsosCFDI) TableName() string {
	return "cfdi_40_usos_cfdi"
}

// RegimenesFiscales model for cfdi_40_regimenes_fiscales
type RegimenesFiscales struct {
	ID            string `gorm:"primaryKey"`
	Texto         string
	AplicaFisica  int
	AplicaMoral   int
	VigenciaDesde string
	VigenciaHasta string
}

func (RegimenesFiscales) TableName() string {
	return "cfdi_40_regimenes_fiscales"
}

// TiposComprobantes model for cfdi_40_tipos_comprobantes
type TiposComprobantes struct {
	ID            string `gorm:"primaryKey"`
	Texto         string
	ValorMaximo   string
	VigenciaDesde string
	VigenciaHasta string
}

func (TiposComprobantes) TableName() string {
	return "cfdi_40_tipos_comprobantes"
}

// MetodosPago model for cfdi_40_metodos_pago
type MetodosPago struct {
	ID            string `gorm:"primaryKey"`
	Texto         string
	VigenciaDesde string
	VigenciaHasta string
}

func (MetodosPago) TableName() string {
	return "cfdi_40_metodos_pago"
}
