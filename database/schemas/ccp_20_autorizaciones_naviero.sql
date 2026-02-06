CREATE TABLE IF NOT EXISTS `ccp_20_autorizaciones_naviero`(
  `id` VARCHAR(255) not null,
  `vigencia_desde` TEXT not null,
  `vigencia_hasta` TEXT not null,
  PRIMARY KEY(`id`)
);
