CREATE TABLE IF NOT EXISTS `ccp_30_estaciones`(
  `id` VARCHAR(255) not null,
  `texto` TEXT not null,
  `clave_transporte` TEXT not null,
  `nacionalidad` TEXT not null,
  `designador_iata` TEXT not null,
  `linea_ferrea` TEXT not null,
  `vigencia_desde` TEXT not null,
  `vigencia_hasta` TEXT not null,
  PRIMARY KEY(`id`)
);
