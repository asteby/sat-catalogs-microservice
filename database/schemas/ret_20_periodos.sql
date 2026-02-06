CREATE TABLE IF NOT EXISTS `ret_20_periodos`(
  `id` VARCHAR(255) not null,
  `texto` TEXT not null,
  `vigencia_desde` TEXT not null,
  `vigencia_hasta` TEXT not null,
  PRIMARY KEY(`id`)
);
