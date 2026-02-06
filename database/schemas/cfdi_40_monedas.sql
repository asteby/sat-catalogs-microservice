CREATE TABLE IF NOT EXISTS `cfdi_40_monedas`(
  `id` VARCHAR(255) not null,
  `texto` TEXT not null,
  `decimales` int not null,
  `porcentaje_variacion` int not null,
  `vigencia_desde` TEXT not null,
  `vigencia_hasta` TEXT not null,
  PRIMARY KEY(`id`)
);
