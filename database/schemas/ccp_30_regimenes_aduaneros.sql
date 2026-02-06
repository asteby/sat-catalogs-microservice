CREATE TABLE IF NOT EXISTS `ccp_30_regimenes_aduaneros`(
  `id` VARCHAR(255) not null,
  `texto` TEXT not null,
  `impoexpo` TEXT not null,
  `vigencia_desde` TEXT not null,
  `vigencia_hasta` TEXT not null,
  PRIMARY KEY(`id`)
);
