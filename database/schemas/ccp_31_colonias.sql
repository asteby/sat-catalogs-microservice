CREATE TABLE IF NOT EXISTS `ccp_31_colonias`(
  `colonia` VARCHAR(255) not null,
  `codigo_postal` VARCHAR(255) not null,
  `texto` TEXT not null,
  PRIMARY KEY(`colonia`, `codigo_postal`)
);
