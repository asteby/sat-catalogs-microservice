CREATE TABLE IF NOT EXISTS `cfdi_productos_servicios`(
  `id` VARCHAR(255) not null,
  `texto` TEXT not null,
  `iva_trasladado` int not null,
  `ieps_trasladado` int not null,
  `complemento` TEXT not null,
  `vigencia_desde` TEXT not null,
  `vigencia_hasta` TEXT not null,
  `estimulo_frontera` int not null,
  `similares` TEXT not null,
  PRIMARY KEY(`id`)
);
