CREATE TABLE IF NOT EXISTS `cfdi_codigos_postales`(
  `id` VARCHAR(255) not null,
  `estado` TEXT not null,
  `municipio` TEXT not null,
  `localidad` TEXT not null,
  `estimulo_frontera` int not null,
  `vigencia_desde` TEXT not null,
  `vigencia_hasta` TEXT not null,
  `huso_descripcion` TEXT not null,
  `huso_verano_mes_inicio` TEXT not null,
  `huso_verano_dia_inicio` TEXT not null,
  `huso_verano_hora_inicio` TEXT not null,
  `huso_verano_diferencia` TEXT not null,
  `huso_invierno_mes_inicio` TEXT not null,
  `huso_invierno_dia_inicio` TEXT not null,
  `huso_invierno_hora_inicio` TEXT not null,
  `huso_invierno_diferencia` TEXT not null,
  PRIMARY KEY(`id`)
);
