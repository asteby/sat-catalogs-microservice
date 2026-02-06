CREATE TABLE IF NOT EXISTS `ccp_31_municipios`(
  `municipio` VARCHAR(255) not null,
  `estado` VARCHAR(255) not null,
  `texto` TEXT not null,
  `vigencia_desde` TEXT not null,
  `vigencia_hasta` TEXT not null,
  PRIMARY KEY(`municipio`, `estado`)
);
