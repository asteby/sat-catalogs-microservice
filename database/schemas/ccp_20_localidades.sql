CREATE TABLE IF NOT EXISTS `ccp_20_localidades`(
  `localidad` VARCHAR(255) not null,
  `estado` VARCHAR(255) not null,
  `texto` TEXT not null,
  `vigencia_desde` TEXT not null,
  `vigencia_hasta` TEXT not null,
  PRIMARY KEY(`localidad`, `estado`)
);
