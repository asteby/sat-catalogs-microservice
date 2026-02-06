CREATE TABLE IF NOT EXISTS `cfdi_40_estados`(
  `estado` VARCHAR(255) not null,
  `pais` VARCHAR(255) not null,
  `texto` TEXT not null,
  `vigencia_desde` TEXT not null,
  `vigencia_hasta` TEXT not null,
  PRIMARY KEY(`estado`, `pais`)
);
