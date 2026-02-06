CREATE TABLE IF NOT EXISTS `cfdi_40_paises`(
  `id` VARCHAR(255) not null,
  `texto` TEXT not null,
  `patron_codigo_postal` TEXT not null,
  `patron_identidad_tributaria` TEXT not null,
  `validacion_identidad_tributaria` TEXT not null,
  `agrupaciones` TEXT not null,
  PRIMARY KEY(`id`)
);
