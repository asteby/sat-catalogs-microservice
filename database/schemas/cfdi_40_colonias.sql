CREATE TABLE IF NOT EXISTS `cfdi_40_colonias`(
  `colonia` VARCHAR(255) not null,
  `codigo_postal` VARCHAR(255) not null,
  `texto` TEXT not null,
  PRIMARY KEY(`colonia`, `codigo_postal`)
);
