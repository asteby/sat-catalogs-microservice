CREATE TABLE IF NOT EXISTS `ccp_20_contenedores_maritimos`(
  `id` VARCHAR(255) not null,
  `texto` TEXT not null,
  `vigencia_desde` TEXT not null,
  `vigencia_hasta` TEXT not null,
  PRIMARY KEY(`id`)
);
