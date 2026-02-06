CREATE TABLE IF NOT EXISTS `cfdi_40_reglas_tasa_cuota`(
  `tipo` TEXT not null,
  `minimo` TEXT not null,
  `valor` TEXT not null,
  `impuesto` TEXT not null,
  `factor` TEXT not null,
  `traslado` int not null,
  `retencion` int not null,
  `vigencia_desde` TEXT not null,
  `vigencia_hasta` TEXT not null
);
