CREATE TABLE IF NOT EXISTS `ccp_31_configuraciones_autotransporte`(
  `id` VARCHAR(255) not null,
  `texto` TEXT not null,
  `numero_de_ejes` int not null,
  `numero_de_llantas` int not null,
  `remolque` TEXT not null,
  `vigencia_desde` TEXT not null,
  `vigencia_hasta` TEXT not null,
  PRIMARY KEY(`id`)
);
