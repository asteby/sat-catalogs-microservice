CREATE TABLE IF NOT EXISTS `ccp_31_claves_unidades`(
  `id` VARCHAR(255) not null,
  `texto` TEXT not null,
  `descripcion` TEXT not null,
  `nota` TEXT not null,
  `vigencia_desde` TEXT not null,
  `vigencia_hasta` TEXT not null,
  `simbolo` TEXT not null,
  `bandera` TEXT not null,
  PRIMARY KEY(`id`)
);
