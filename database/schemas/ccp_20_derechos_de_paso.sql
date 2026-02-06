CREATE TABLE IF NOT EXISTS `ccp_20_derechos_de_paso`(
  `id` VARCHAR(255) not null,
  `texto` TEXT not null,
  `entre` TEXT not null,
  `hasta` TEXT not null,
  `otorga_recibe` TEXT not null,
  `concesionario` TEXT not null,
  `vigencia_desde` TEXT not null,
  `vigencia_hasta` TEXT not null,
  PRIMARY KEY(`id`)
);
