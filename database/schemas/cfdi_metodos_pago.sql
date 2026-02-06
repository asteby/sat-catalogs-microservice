CREATE TABLE IF NOT EXISTS `cfdi_metodos_pago`(
  `id` VARCHAR(255) not null,
  `texto` TEXT not null,
  `vigencia_desde` TEXT not null,
  `vigencia_hasta` TEXT not null,
  PRIMARY KEY(`id`)
);
