CREATE TABLE IF NOT EXISTS `cce_20_fracciones_arancelarias`(
  `fraccion` VARCHAR(255) not null,
  `texto` TEXT not null,
  `vigencia_desde` TEXT not null,
  `vigencia_hasta` TEXT not null,
  `unidad` TEXT not null,
  PRIMARY KEY(`fraccion`)
);
