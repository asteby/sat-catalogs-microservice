CREATE TABLE IF NOT EXISTS `cfdi_usos_cfdi`(
  `id` VARCHAR(255) not null,
  `texto` TEXT not null,
  `aplica_fisica` int not null,
  `aplica_moral` int not null,
  `vigencia_desde` TEXT not null,
  `vigencia_hasta` TEXT not null,
  PRIMARY KEY(`id`)
);
