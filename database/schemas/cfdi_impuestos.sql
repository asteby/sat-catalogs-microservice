CREATE TABLE IF NOT EXISTS `cfdi_impuestos`(
  `id` VARCHAR(255) not null,
  `texto` TEXT not null,
  `retencion` int not null,
  `traslado` int not null,
  `ambito` TEXT not null,
  `entidad` TEXT not null,
  PRIMARY KEY(`id`)
);
