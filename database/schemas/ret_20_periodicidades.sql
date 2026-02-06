CREATE TABLE IF NOT EXISTS `ret_20_periodicidades`(
  `id` VARCHAR(255) not null,
  `texto` TEXT not null,
  `nombre_complemento` TEXT not null,
  `vigencia_desde` TEXT not null,
  `vigencia_hasta` TEXT not null,
  PRIMARY KEY(`id`)
);
