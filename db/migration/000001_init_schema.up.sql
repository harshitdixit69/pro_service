CREATE TABLE `users` (
  `id` INT AUTO_INCREMENT PRIMARY KEY, -- Auto-increment primary key
  `username` VARCHAR(255),
  `first_name` VARCHAR(255),
  `last_name` VARCHAR(255),
  `pre_password` VARCHAR(255),
  `password` VARCHAR(255),
  `encrypted_passwrod` VARCHAR(255),
  `address` VARCHAR(255),
  `phone_no` VARCHAR(255),
  `aadhaar_no` VARCHAR(255),
  `status` INT,
  `is_active` INT,
  `created_date` TIMESTAMP,
  `updated_date` TIMESTAMP,
  `deleted_date` TIMESTAMP
);

CREATE TABLE `services` (
  `id` INT AUTO_INCREMENT PRIMARY KEY, -- Auto-increment primary key
  `service_type` VARCHAR(255),
  `is_active` INT,
  `status` INT,
  `created_date` TIMESTAMP,
  `updated_date` TIMESTAMP,
  `deleted_date` TIMESTAMP
);

CREATE TABLE `service_provider` (
  `id` INT AUTO_INCREMENT PRIMARY KEY, -- Auto-increment primary key
  `service_id` INT,
  `username` VARCHAR(255),
  `first_name` VARCHAR(255),
  `last_name` VARCHAR(255),
  `pre_password` VARCHAR(255),
  `password` VARCHAR(255),
  `encrypted_passwrod` VARCHAR(255),
  `address` VARCHAR(255),
  `phone_no` VARCHAR(255),
  `aadhaar_no` VARCHAR(255),
  `status` INT,
  `is_active` INT,
  `created_date` TIMESTAMP,
  `updated_date` TIMESTAMP,
  `deleted_date` TIMESTAMP
);

CREATE TABLE `availability` (
  `id` INT AUTO_INCREMENT PRIMARY KEY, -- Auto-increment primary key
  `service_id` INT,
  `service_provider_id` INT,
  `status` INT,
  `is_active` INT,
  `created_date` TIMESTAMP,
  `updated_date` TIMESTAMP,
  `deleted_date` TIMESTAMP
);

CREATE TABLE `type_availability` (
  `id` INT AUTO_INCREMENT PRIMARY KEY, -- Auto-increment primary key
  `availability_id` INT,
  `service_provider_id` INT,
  `status` INT,
  `is_active` INT,
  `created_date` TIMESTAMP,
  `updated_date` TIMESTAMP,
  `deleted_date` TIMESTAMP
);

CREATE TABLE `appointermnts` (
  `id` INT AUTO_INCREMENT PRIMARY KEY, -- Auto-increment primary key
  `user_id` INT,
  `type_availability_id` INT,
  `is_accepted` INT,
  `is_rejected` INT,
  `is_pending` INT,
  `status` INT,
  `created_date` TIMESTAMP,
  `updated_date` TIMESTAMP,
  `deleted_date` TIMESTAMP
);

-- Adding foreign keys
ALTER TABLE `service_provider` 
ADD FOREIGN KEY (`service_id`) REFERENCES `services` (`id`);

ALTER TABLE `availability` 
ADD FOREIGN KEY (`service_id`) REFERENCES `services` (`id`),
ADD FOREIGN KEY (`service_provider_id`) REFERENCES `service_provider` (`id`);

ALTER TABLE `type_availability` 
ADD FOREIGN KEY (`availability_id`) REFERENCES `availability` (`id`),
ADD FOREIGN KEY (`service_provider_id`) REFERENCES `service_provider` (`id`);

ALTER TABLE `appointermnts` 
ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
ADD FOREIGN KEY (`type_availability_id`) REFERENCES `type_availability` (`id`);
