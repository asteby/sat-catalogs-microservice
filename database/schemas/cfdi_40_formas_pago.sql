CREATE TABLE IF NOT EXISTS `cfdi_40_formas_pago`(
  `id` VARCHAR(255) not null,
  `texto` TEXT not null,
  `es_bancarizado` int not null,
  `requiere_numero_operacion` int not null,
  `permite_banco_ordenante_rfc` int not null,
  `permite_cuenta_ordenante` int not null,
  `patron_cuenta_ordenante` TEXT not null,
  `permite_banco_beneficiario_rfc` int not null,
  `permite_cuenta_beneficiario` int not null,
  `patron_cuenta_beneficiario` TEXT not null,
  `permite_tipo_cadena_pago` int not null,
  `requiere_banco_ordenante_nombre_ext` int not null,
  `vigencia_desde` TEXT not null,
  `vigencia_hasta` TEXT not null,
  PRIMARY KEY(`id`)
);
