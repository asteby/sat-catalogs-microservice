CREATE TABLE IF NOT EXISTS `ccp_31_tipos_estacion`(
  `id` VARCHAR(255) not null,
  `texto` TEXT not null,
  `claves_transportes` TEXT not null,
  `vigencia_desde` TEXT not null,
  `vigencia_hasta` TEXT not null,
  PRIMARY KEY(`id`)
);
