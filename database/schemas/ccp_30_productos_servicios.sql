CREATE TABLE IF NOT EXISTS `ccp_30_productos_servicios`(
  `id` VARCHAR(255) not null,
  `texto` TEXT not null,
  `similares` TEXT not null,
  `material_peligroso` TEXT not null,
  `vigencia_desde` TEXT not null,
  `vigencia_hasta` TEXT not null,
  PRIMARY KEY(`id`)
);
