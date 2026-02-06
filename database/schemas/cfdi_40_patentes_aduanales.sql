CREATE TABLE IF NOT EXISTS `cfdi_40_patentes_aduanales`(
  `id` VARCHAR(255) not null,
  `vigencia_desde` TEXT not null,
  `vigencia_hasta` TEXT not null,
  PRIMARY KEY(`id`)
);
