ALTER TABLE type_availability ADD `date` DATE NULL;
ALTER TABLE type_availability ADD `time` DATETIME NULL;
ALTER TABLE type_availability ADD temp_block_schedule INT UNSIGNED NULL;
ALTER TABLE type_availability ADD schedule_confirmed INT UNSIGNED NULL;
ALTER TABLE type_availability ADD schedule_cancelled INT UNSIGNED NULL;
