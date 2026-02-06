CREATE TABLE IF NOT EXISTS `cce_colonias`(
  `colonia` VARCHAR(255) not null,
  `codigo_postal` VARCHAR(255) not null,
  `asentamiento` TEXT not null,
  PRIMARY KEY(`colonia`, `codigo_postal`)
);
