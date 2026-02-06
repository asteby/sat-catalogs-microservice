CREATE TABLE IF NOT EXISTS `ccp_30_materiales_peligrosos`(
  `id` TEXT not null,
  `texto` TEXT not null,
  `clase_o_div` TEXT not null,
  `peligro_secundario` TEXT not null,
  `nombre_tecnico` TEXT not null,
  `vigencia_desde` TEXT not null,
  `vigencia_hasta` TEXT not null
);
