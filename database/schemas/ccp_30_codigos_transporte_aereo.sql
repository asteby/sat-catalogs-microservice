CREATE TABLE IF NOT EXISTS `ccp_30_codigos_transporte_aereo`(
  `id` VARCHAR(255) not null,
  `nacionalidad` TEXT not null,
  `texto` TEXT not null,
  `designador_oaci` TEXT not null,
  `vigencia_desde` TEXT not null,
  `vigencia_hasta` TEXT not null,
  PRIMARY KEY(`id`)
);
