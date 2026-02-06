CREATE TABLE IF NOT EXISTS `cfdi_tipos_comprobantes`(
  `id` VARCHAR(255) not null,
  `texto` TEXT not null,
  `valor_maximo` TEXT not null,
  `vigencia_desde` TEXT not null,
  `vigencia_hasta` TEXT not null,
  PRIMARY KEY(`id`)
);
