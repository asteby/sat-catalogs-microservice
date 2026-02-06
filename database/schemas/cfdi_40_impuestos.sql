CREATE TABLE IF NOT EXISTS `cfdi_40_impuestos`(
  `id` VARCHAR(255) not null,
  `texto` TEXT not null,
  `retencion` int not null,
  `traslado` int not null,
  `ambito` TEXT not null,
  `vigencia_desde` TEXT not null,
  `vigencia_hasta` TEXT not null,
  PRIMARY KEY(`id`)
);
