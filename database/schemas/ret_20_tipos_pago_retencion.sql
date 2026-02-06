CREATE TABLE IF NOT EXISTS `ret_20_tipos_pago_retencion`(
  `id` VARCHAR(255) not null,
  `texto` TEXT not null,
  `tipo_impuesto` TEXT not null,
  `vigencia_desde` TEXT not null,
  `vigencia_hasta` TEXT not null,
  PRIMARY KEY(`id`)
);
