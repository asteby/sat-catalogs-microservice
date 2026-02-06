CREATE TABLE IF NOT EXISTS `nomina_bancos`(
  `id` VARCHAR(255) not null,
  `texto` TEXT not null,
  `razon_social` TEXT not null,
  `vigencia_desde` TEXT not null,
  `vigencia_hasta` TEXT not null,
  PRIMARY KEY(`id`)
);
