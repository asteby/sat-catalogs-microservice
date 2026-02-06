CREATE TABLE IF NOT EXISTS `cfdi_claves_unidades`(
  `id` VARCHAR(255) not null,
  `texto` TEXT not null,
  `descripcion` TEXT not null,
  `notas` TEXT not null,
  `vigencia_desde` TEXT not null,
  `vigencia_hasta` TEXT not null,
  `simbolo` TEXT not null,
  PRIMARY KEY(`id`)
);
