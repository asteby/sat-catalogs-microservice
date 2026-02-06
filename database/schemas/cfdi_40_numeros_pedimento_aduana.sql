CREATE TABLE IF NOT EXISTS `cfdi_40_numeros_pedimento_aduana`(
  `aduana` TEXT not null,
  `patente` TEXT not null,
  `ejercicio` int not null,
  `cantidad` int not null,
  `vigencia_desde` TEXT not null,
  `vigencia_hasta` TEXT not null
);
