CREATE TABLE IF NOT EXISTS `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `age` int DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_email` (`email`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `webhook_endpoints` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `url` varchar(512) NOT NULL,
  `events` json DEFAULT (JSON_ARRAY()),
  `secret_key` varchar(255) NOT NULL,
  `enabled` tinyint(1) DEFAULT '1',
  `retry_count` int DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_webhook_endpoints_name` (`name`),
  KEY `idx_webhook_endpoints_enabled` (`enabled`),
  KEY `idx_webhook_endpoints_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `integration_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `system_name` varchar(255) NOT NULL,
  `operation` varchar(255) NOT NULL,
  `request_data` json DEFAULT (JSON_OBJECT()),
  `response_data` json DEFAULT (JSON_OBJECT()),
  `status` varchar(50) NOT NULL,
  `execution_time_ms` int DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_integration_logs_system_name` (`system_name`),
  KEY `idx_integration_logs_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `api_usage_stats` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `partner_id` varchar(255) NOT NULL,
  `date` date NOT NULL,
  `endpoint` varchar(255) NOT NULL,
  `request_count` int NOT NULL,
  `success_count` int NOT NULL,
  `error_count` int NOT NULL,
  `avg_response_time_ms` int NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_api_usage_stats_partner_date_endpoint` (`partner_id`,`date`,`endpoint`),
  KEY `idx_api_usage_stats_date` (`date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `rate_limits` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `partner_id` varchar(255) NOT NULL,
  `endpoint` varchar(255) NOT NULL,
  `limit_per_hour` int NOT NULL,
  `limit_per_day` int NOT NULL,
  `current_hour_count` int DEFAULT '0',
  `current_day_count` int DEFAULT '0',
  `reset_time` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_rate_limits_partner_endpoint` (`partner_id`,`endpoint`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `external_products` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `external_id` varchar(255) NOT NULL,
  `partner_id` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `price` decimal(10,2) NOT NULL,
  `stock` int NOT NULL,
  `category` varchar(255) DEFAULT NULL,
  `last_synced` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_external_products_external_partner` (`external_id`,`partner_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `crm_sync_history` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `customer_id` varchar(255) NOT NULL,
  `crm_system` varchar(255) NOT NULL,
  `sync_type` varchar(255) NOT NULL,
  `data_sent` json DEFAULT (JSON_OBJECT()),
  `sync_status` varchar(50) NOT NULL,
  `sync_id` varchar(255) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_crm_sync_history_customer_id` (`customer_id`),
  KEY `idx_crm_sync_history_crm_system` (`crm_system`),
  KEY `idx_crm_sync_history_sync_status` (`sync_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `email_jobs` (
  `id` varchar(255) NOT NULL,
  `job_type` varchar(255) NOT NULL,
  `recipient` varchar(255) NOT NULL,
  `template` varchar(255) NOT NULL,
  `template_data` json DEFAULT (JSON_OBJECT()),
  `provider` varchar(100) NOT NULL,
  `status` varchar(50) NOT NULL,
  `scheduled_at` datetime(3) DEFAULT NULL,
  `sent_at` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_email_jobs_status` (`status`),
  KEY `idx_email_jobs_provider` (`provider`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `shipping_tracking` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `order_id` varchar(255) NOT NULL,
  `tracking_number` varchar(255) NOT NULL,
  `shipping_provider` varchar(100) NOT NULL,
  `status` varchar(50) NOT NULL,
  `estimated_delivery` datetime(3) DEFAULT NULL,
  `delivered_at` datetime(3) DEFAULT NULL,
  `recipient_signature` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_shipping_tracking_number` (`tracking_number`),
  KEY `idx_shipping_tracking_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `analytics_events` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `event_name` varchar(255) NOT NULL,
  `user_id` varchar(255) DEFAULT NULL,
  `properties` json DEFAULT (JSON_OBJECT()),
  `provider` varchar(100) NOT NULL,
  `sent_at` datetime(3) DEFAULT NULL,
  `status` varchar(50) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_analytics_events_event_name` (`event_name`),
  KEY `idx_analytics_events_provider` (`provider`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `integration_health_checks` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `system_name` varchar(255) NOT NULL,
  `endpoint` varchar(512) NOT NULL,
  `status` varchar(50) NOT NULL,
  `response_time_ms` int NOT NULL,
  `last_check` datetime(3) DEFAULT NULL,
  `error_message` text DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_integration_health_checks_system_name` (`system_name`),
  KEY `idx_integration_health_checks_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `webhook_deliveries` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `webhook_endpoint_id` bigint unsigned NOT NULL,
  `event_type` varchar(255) NOT NULL,
  `payload` json DEFAULT (JSON_OBJECT()),
  `response_status` int DEFAULT NULL,
  `response_body` text DEFAULT NULL,
  `delivered_at` datetime(3) DEFAULT NULL,
  `retry_count` int DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_webhook_deliveries_endpoint_id` (`webhook_endpoint_id`),
  KEY `idx_webhook_deliveries_delivered_at` (`delivered_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `connection_test_results` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `system_name` varchar(255) NOT NULL,
  `test_type` varchar(100) NOT NULL,
  `status` varchar(50) NOT NULL,
  `response_time_ms` int DEFAULT NULL,
  `error_details` text DEFAULT NULL,
  `tested_at` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_connection_test_results_system_name` (`system_name`),
  KEY `idx_connection_test_results_test_type` (`test_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `integration_error_stats` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `date` date NOT NULL,
  `system_name` varchar(255) NOT NULL,
  `error_type` varchar(255) NOT NULL,
  `error_count` int NOT NULL,
  `most_common_error` text DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_integration_error_stats` (`date`,`system_name`,`error_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `auths` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `email` varchar(255) NOT NULL,
  `password_hash` varchar(255) NOT NULL,
  `is_active` tinyint(1) DEFAULT '1',
  `last_login_at` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_auths_user_id` (`user_id`),
  UNIQUE KEY `idx_auths_email` (`email`),
  KEY `idx_auths_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `roles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `description` varchar(255) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_roles_name` (`name`),
  KEY `idx_roles_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `permissions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `resource` varchar(255) NOT NULL,
  `action` varchar(255) NOT NULL,
  `scope` varchar(255) NOT NULL DEFAULT '*',
  `description` varchar(255) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_permissions_name` (`name`),
  KEY `idx_permissions_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `scopes` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `resource` varchar(255) NOT NULL,
  `description` varchar(255) DEFAULT NULL,
  `is_active` tinyint(1) DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_scopes_name` (`name`),
  KEY `idx_scopes_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `user_scopes` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `scope_id` bigint unsigned NOT NULL,
  `expires_at` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_scopes_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `user_roles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `role_id` bigint unsigned NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_roles_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `role_permissions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `role_id` bigint unsigned NOT NULL,
  `permission_id` bigint unsigned NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_role_permissions_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `refresh_tokens` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `token` varchar(255) NOT NULL,
  `expires_at` datetime(3) DEFAULT NULL,
  `is_revoked` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_refresh_tokens_token` (`token`),
  KEY `idx_refresh_tokens_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `user_tokens` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `token_type` varchar(50) NOT NULL,
  `token_hash` varchar(255) NOT NULL,
  `expires_at` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_tokens_user_id` (`user_id`),
  KEY `idx_user_tokens_token_type` (`token_type`),
  KEY `idx_user_tokens_expires_at` (`expires_at`),
  KEY `idx_user_tokens_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `membership_tiers` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `level` int NOT NULL,
  `description` varchar(255) DEFAULT NULL,
  `benefits` json DEFAULT NULL,
  `requirements` json DEFAULT NULL,
  `is_active` tinyint(1) DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_membership_tiers_name` (`name`),
  KEY `idx_membership_tiers_level` (`level`),
  KEY `idx_membership_tiers_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `user_memberships` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `tier_id` bigint unsigned NOT NULL,
  `points` int DEFAULT '0',
  `total_spent` decimal(10,2) DEFAULT '0.00',
  `joined_at` datetime(3) DEFAULT NULL,
  `last_activity_at` datetime(3) DEFAULT NULL,
  `expires_at` datetime(3) DEFAULT NULL,
  `is_active` tinyint(1) DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_memberships_user_id` (`user_id`),
  KEY `idx_user_memberships_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `point_transactions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `type` varchar(255) NOT NULL,
  `points` int NOT NULL,
  `description` varchar(255) DEFAULT NULL,
  `reference_type` varchar(255) DEFAULT NULL,
  `reference_id` bigint unsigned DEFAULT NULL,
  `expires_at` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_point_transactions_user_id` (`user_id`),
  KEY `idx_point_transactions_type` (`type`),
  KEY `idx_point_transactions_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `user_profiles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `first_name` varchar(255) DEFAULT NULL,
  `last_name` varchar(255) DEFAULT NULL,
  `phone_number` varchar(255) DEFAULT NULL,
  `date_of_birth` datetime(3) DEFAULT NULL,
  `gender` varchar(255) DEFAULT NULL,
  `address` json DEFAULT NULL,
  `preferences` json DEFAULT NULL,
  `avatar` varchar(255) DEFAULT NULL,
  `bio` varchar(255) DEFAULT NULL,
  `is_verified` tinyint(1) DEFAULT '0',
  `verified_at` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_profiles_user_id` (`user_id`),
  KEY `idx_user_profiles_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `user_activities` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `activity_type` varchar(255) NOT NULL,
  `description` varchar(255) DEFAULT NULL,
  `metadata` json DEFAULT NULL,
  `ip_address` varchar(255) DEFAULT NULL,
  `user_agent` varchar(255) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_activities_user_id` (`user_id`),
  KEY `idx_user_activities_activity_type` (`activity_type`),
  KEY `idx_user_activities_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `notifications` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `type` varchar(255) NOT NULL,
  `title` varchar(255) NOT NULL,
  `message` varchar(255) NOT NULL,
  `data` json DEFAULT (JSON_OBJECT()),
  `is_read` tinyint(1) DEFAULT '0',
  `read_at` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_notifications_user_id` (`user_id`),
  KEY `idx_notifications_type` (`type`),
  KEY `idx_notifications_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `user_preferences` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `category` varchar(255) NOT NULL,
  `key` varchar(255) NOT NULL,
  `value` varchar(255) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_preferences_user_id` (`user_id`),
  KEY `idx_user_preferences_category` (`category`),
  KEY `idx_user_preferences_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `login_attempts` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `email` varchar(255) NOT NULL,
  `ip_address` varchar(255) NOT NULL,
  `user_agent` varchar(255) DEFAULT NULL,
  `success` tinyint(1) DEFAULT NULL,
  `fail_reason` varchar(255) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_login_attempts_email` (`email`),
  KEY `idx_login_attempts_ip_address` (`ip_address`),
  KEY `idx_login_attempts_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `ip_blacklists` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `ip_address` varchar(255) NOT NULL,
  `reason` varchar(255) DEFAULT NULL,
  `expires_at` datetime(3) DEFAULT NULL,
  `is_active` tinyint(1) DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_ip_blacklists_ip_address` (`ip_address`),
  KEY `idx_ip_blacklists_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `user_sessions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `session_id` varchar(255) NOT NULL,
  `ip_address` varchar(255) NOT NULL,
  `user_agent` varchar(255) DEFAULT NULL,
  `expires_at` datetime(3) DEFAULT NULL,
  `is_active` tinyint(1) DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_sessions_session_id` (`session_id`),
  KEY `idx_user_sessions_user_id` (`user_id`),
  KEY `idx_user_sessions_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `security_events` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned DEFAULT NULL,
  `event_type` varchar(255) NOT NULL,
  `description` varchar(255) DEFAULT NULL,
  `ip_address` varchar(255) NOT NULL,
  `user_agent` varchar(255) DEFAULT NULL,
  `severity` varchar(255) NOT NULL,
  `metadata` json DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_security_events_user_id` (`user_id`),
  KEY `idx_security_events_event_type` (`event_type`),
  KEY `idx_security_events_ip_address` (`ip_address`),
  KEY `idx_security_events_severity` (`severity`),
  KEY `idx_security_events_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `rate_limit_rules` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `resource` varchar(255) NOT NULL,
  `max_requests` int NOT NULL,
  `window_size` int NOT NULL,
  `is_active` tinyint(1) DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_rate_limit_rules_name` (`name`),
  KEY `idx_rate_limit_rules_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `rate_limit_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `rule_id` bigint unsigned NOT NULL,
  `ip_address` varchar(255) NOT NULL,
  `user_id` bigint unsigned DEFAULT NULL,
  `requests` int DEFAULT NULL,
  `window_start` datetime(3) DEFAULT NULL,
  `window_end` datetime(3) DEFAULT NULL,
  `blocked` tinyint(1) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_rate_limit_logs_rule_id` (`rule_id`),
  KEY `idx_rate_limit_logs_ip_address` (`ip_address`),
  KEY `idx_rate_limit_logs_user_id` (`user_id`),
  KEY `idx_rate_limit_logs_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `device_fingerprints` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `fingerprint` varchar(255) NOT NULL,
  `device_info` json DEFAULT NULL,
  `is_trusted` tinyint(1) DEFAULT '0',
  `last_seen_at` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_device_fingerprints_user_id` (`user_id`),
  KEY `idx_device_fingerprints_fingerprint` (`fingerprint`),
  KEY `idx_device_fingerprints_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `user_suspensions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `reason` varchar(255) NOT NULL,
  `duration_days` int NOT NULL,
  `suspended_by` bigint unsigned NOT NULL,
  `suspended_at` datetime(3) NOT NULL,
  `expires_at` datetime(3) NOT NULL,
  `status` varchar(255) NOT NULL DEFAULT 'active',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_suspensions_user_id` (`user_id`),
  KEY `idx_user_suspensions_suspended_by` (`suspended_by`),
  KEY `idx_user_suspensions_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `user_approval_queue` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `registration_data` json DEFAULT NULL,
  `risk_assessment` json DEFAULT NULL,
  `assigned_to` bigint unsigned DEFAULT NULL,
  `priority` varchar(255) NOT NULL DEFAULT 'normal',
  `status` varchar(255) NOT NULL DEFAULT 'pending',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_approval_queue_user_id` (`user_id`),
  KEY `idx_user_approval_queue_assigned_to` (`assigned_to`),
  KEY `idx_user_approval_queue_status` (`status`),
  KEY `idx_user_approval_queue_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `system_health_metrics` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `metric_name` varchar(255) NOT NULL,
  `value` varchar(255) NOT NULL,
  `unit` varchar(255) DEFAULT NULL,
  `status` varchar(255) NOT NULL DEFAULT 'healthy',
  `last_updated` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_system_health_metrics_metric_name` (`metric_name`),
  KEY `idx_system_health_metrics_status` (`status`),
  KEY `idx_system_health_metrics_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `fraud_alerts` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `alert_type` varchar(255) NOT NULL,
  `severity` varchar(255) NOT NULL,
  `title` varchar(255) NOT NULL,
  `description` text DEFAULT NULL,
  `status` varchar(255) NOT NULL DEFAULT 'active',
  `triggered_at` datetime(3) NOT NULL,
  `resolved_at` datetime(3) DEFAULT NULL,
  `resolved_by` bigint unsigned DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_fraud_alerts_user_id` (`user_id`),
  KEY `idx_fraud_alerts_alert_type` (`alert_type`),
  KEY `idx_fraud_alerts_severity` (`severity`),
  KEY `idx_fraud_alerts_status` (`status`),
  KEY `idx_fraud_alerts_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `admin_actions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `admin_user_id` bigint unsigned NOT NULL,
  `action_type` varchar(255) NOT NULL,
  `target_type` varchar(255) NOT NULL,
  `target_id` varchar(255) DEFAULT NULL,
  `description` text DEFAULT NULL,
  `result` varchar(255) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_admin_actions_admin_user_id` (`admin_user_id`),
  KEY `idx_admin_actions_action_type` (`action_type`),
  KEY `idx_admin_actions_target_type` (`target_type`),
  KEY `idx_admin_actions_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `system_settings_history` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `setting_key` varchar(255) NOT NULL,
  `old_value` varchar(255) DEFAULT NULL,
  `new_value` varchar(255) NOT NULL,
  `changed_by` bigint unsigned NOT NULL,
  `change_reason` text DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_system_settings_history_setting_key` (`setting_key`),
  KEY `idx_system_settings_history_changed_by` (`changed_by`),
  KEY `idx_system_settings_history_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `admin_settings` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `key_name` varchar(255) NOT NULL,
  `value` varchar(255) NOT NULL,
  `description` text DEFAULT NULL,
  `category` varchar(255) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_admin_settings_key_name` (`key_name`),
  KEY `idx_admin_settings_category` (`category`),
  KEY `idx_admin_settings_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `export_jobs` (
  `id` varchar(255) NOT NULL,
  `job_type` varchar(255) NOT NULL,
  `status` varchar(255) NOT NULL DEFAULT 'pending',
  `parameters` json DEFAULT NULL,
  `created_by` bigint unsigned NOT NULL,
  `progress` int DEFAULT '0',
  `file_path` varchar(255) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_export_jobs_job_type` (`job_type`),
  KEY `idx_export_jobs_status` (`status`),
  KEY `idx_export_jobs_created_by` (`created_by`),
  KEY `idx_export_jobs_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `announcements` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `content` text NOT NULL,
  `type` varchar(255) NOT NULL,
  `priority` varchar(255) NOT NULL DEFAULT 'medium',
  `target_users` varchar(255) NOT NULL DEFAULT 'all',
  `status` varchar(255) NOT NULL DEFAULT 'draft',
  `created_by` bigint unsigned NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_announcements_type` (`type`),
  KEY `idx_announcements_priority` (`priority`),
  KEY `idx_announcements_status` (`status`),
  KEY `idx_announcements_created_by` (`created_by`),
  KEY `idx_announcements_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `audit_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `action` varchar(255) NOT NULL,
  `category` varchar(255) NOT NULL,
  `resource_type` varchar(255) NOT NULL,
  `resource_id` varchar(255) DEFAULT NULL,
  `details` json DEFAULT NULL,
  `ip_address` varchar(255) NOT NULL,
  `user_agent` varchar(255) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_audit_logs_user_id` (`user_id`),
  KEY `idx_audit_logs_action` (`action`),
  KEY `idx_audit_logs_category` (`category`),
  KEY `idx_audit_logs_resource_type` (`resource_type`),
  KEY `idx_audit_logs_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `performance_metrics` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `metric_name` varchar(255) NOT NULL,
  `value` decimal(10,2) NOT NULL,
  `unit` varchar(50) DEFAULT NULL,
  `timestamp` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_performance_metrics_metric_name` (`metric_name`),
  KEY `idx_performance_metrics_timestamp` (`timestamp`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `connection_pool_config` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `pool_name` varchar(255) NOT NULL,
  `max_connections` int NOT NULL,
  `min_connections` int NOT NULL,
  `connection_timeout_ms` int NOT NULL,
  `idle_timeout_ms` int NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_connection_pool_config_pool_name` (`pool_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `cache_config` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `cache_name` varchar(255) NOT NULL,
  `max_size_mb` int NOT NULL,
  `ttl_seconds` int NOT NULL,
  `eviction_policy` varchar(50) NOT NULL,
  `hit_ratio_target` decimal(3,2) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_cache_config_cache_name` (`cache_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `batch_jobs` (
  `id` varchar(255) NOT NULL,
  `job_name` varchar(255) NOT NULL,
  `job_type` varchar(255) NOT NULL,
  `priority` varchar(50) NOT NULL,
  `max_execution_time_minutes` int NOT NULL,
  `retry_count` int NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_batch_jobs_job_type` (`job_type`),
  KEY `idx_batch_jobs_priority` (`priority`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `job_queue` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `job_type` varchar(255) NOT NULL,
  `priority` varchar(50) NOT NULL,
  `payload` json DEFAULT NULL,
  `status` varchar(50) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_job_queue_job_type` (`job_type`),
  KEY `idx_job_queue_status` (`status`),
  KEY `idx_job_queue_priority` (`priority`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `api_response_stats` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `endpoint` varchar(255) NOT NULL,
  `method` varchar(10) NOT NULL,
  `p50_ms` int NOT NULL,
  `p90_ms` int NOT NULL,
  `p95_ms` int NOT NULL,
  `p99_ms` int NOT NULL,
  `request_count` int NOT NULL,
  `date` date NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_api_response_stats_endpoint` (`endpoint`),
  KEY `idx_api_response_stats_date` (`date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `index_usage_stats` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `table_name` varchar(255) NOT NULL,
  `index_name` varchar(255) NOT NULL,
  `usage_count` int NOT NULL,
  `efficiency_ratio` decimal(3,2) NOT NULL,
  `last_used` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_index_usage_stats_table_name` (`table_name`),
  KEY `idx_index_usage_stats_index_name` (`index_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `resource_usage_history` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `timestamp` datetime(3) NOT NULL,
  `cpu_percent` decimal(5,2) NOT NULL,
  `memory_percent` decimal(5,2) NOT NULL,
  `disk_io_mb_per_sec` decimal(8,2) NOT NULL,
  `network_io_mb_per_sec` decimal(8,2) NOT NULL,
  `active_connections` int NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_resource_usage_history_timestamp` (`timestamp`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `load_test_results` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `test_name` varchar(255) NOT NULL,
  `concurrent_users` int NOT NULL,
  `duration_minutes` int NOT NULL,
  `total_requests` int NOT NULL,
  `successful_requests` int NOT NULL,
  `failed_requests` int NOT NULL,
  `avg_response_time_ms` int NOT NULL,
  `max_response_time_ms` int NOT NULL,
  `throughput_rps` decimal(8,2) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_load_test_results_test_name` (`test_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `memory_analysis` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `component` varchar(255) NOT NULL,
  `allocated_mb` int NOT NULL,
  `used_mb` int NOT NULL,
  `free_mb` int NOT NULL,
  `gc_count` int NOT NULL,
  `gc_time_ms` int NOT NULL,
  `timestamp` datetime(3) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_memory_analysis_component` (`component`),
  KEY `idx_memory_analysis_timestamp` (`timestamp`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `concurrent_sessions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `session_count` int NOT NULL,
  `max_allowed` int NOT NULL,
  `last_activity` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_concurrent_sessions_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `query_execution_plans` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `query_hash` varchar(255) NOT NULL,
  `query_text` text NOT NULL,
  `execution_count` int NOT NULL,
  `avg_execution_time_ms` int NOT NULL,
  `index_usage` text DEFAULT NULL,
  `optimization_suggestions` text DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_query_execution_plans_query_hash` (`query_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `system_settings` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `key_name` varchar(255) NOT NULL,
  `value` varchar(255) NOT NULL,
  `description` text DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_system_settings_key_name` (`key_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `partner_api_keys` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `partner_id` varchar(255) NOT NULL,
  `api_key` varchar(255) NOT NULL,
  `partner_name` varchar(255) DEFAULT NULL,
  `permissions` json DEFAULT (JSON_ARRAY()),
  `rate_limit` int DEFAULT NULL,
  `enabled` tinyint(1) DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_partner_api_keys_api_key` (`api_key`),
  KEY `idx_partner_api_keys_partner_id` (`partner_id`),
  KEY `idx_partner_api_keys_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `external_systems` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `type` varchar(100) NOT NULL,
  `endpoint_url` varchar(512) DEFAULT NULL,
  `api_version` varchar(50) DEFAULT NULL,
  `authentication_type` varchar(50) DEFAULT NULL,
  `credentials` json DEFAULT (JSON_OBJECT()),
  `status` varchar(50) NOT NULL DEFAULT 'active',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_external_systems_name` (`name`),
  KEY `idx_external_systems_type` (`type`),
  KEY `idx_external_systems_status` (`status`),
  KEY `idx_external_systems_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `products` (
  `id` varchar(255) NOT NULL,
  `store_id` bigint unsigned NOT NULL,
  `name` varchar(255) NOT NULL,
  `price` decimal(10,2) NOT NULL,
  `category` varchar(255) DEFAULT NULL,
  `stock` int NOT NULL DEFAULT '0',
  `description` text DEFAULT NULL,
  `exclusive_tier` varchar(50) DEFAULT NULL,
  `limited_quantity` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_products_store_id` (`store_id`),
  KEY `idx_products_category` (`category`),
  KEY `idx_products_exclusive_tier` (`exclusive_tier`),
  KEY `idx_products_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `vip_events` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `description` text DEFAULT NULL,
  `tier_requirement` varchar(50) NOT NULL,
  `capacity` int NOT NULL,
  `current_attendees` int DEFAULT '0',
  `event_date` datetime(3) NOT NULL,
  `location` varchar(255) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_vip_events_tier_requirement` (`tier_requirement`),
  KEY `idx_vip_events_event_date` (`event_date`),
  KEY `idx_vip_events_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `concierge_staff` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `specialization` varchar(255) DEFAULT NULL,
  `tier_level` varchar(50) NOT NULL,
  `available` tinyint(1) DEFAULT '1',
  `rating` decimal(3,2) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_concierge_staff_tier_level` (`tier_level`),
  KEY `idx_concierge_staff_available` (`available`),
  KEY `idx_concierge_staff_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `designers` (
  `id` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `specialty` varchar(255) DEFAULT NULL,
  `tier_exclusive` varchar(50) DEFAULT NULL,
  `hourly_rate` decimal(10,2) DEFAULT NULL,
  `available` tinyint(1) DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_designers_tier_exclusive` (`tier_exclusive`),
  KEY `idx_designers_available` (`available`),
  KEY `idx_designers_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `private_sales` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `description` text DEFAULT NULL,
  `tier_requirement` varchar(50) NOT NULL,
  `start_date` datetime(3) NOT NULL,
  `end_date` datetime(3) NOT NULL,
  `discount_rate` decimal(3,2) NOT NULL,
  `invitation_only` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_private_sales_tier_requirement` (`tier_requirement`),
  KEY `idx_private_sales_start_date` (`start_date`),
  KEY `idx_private_sales_end_date` (`end_date`),
  KEY `idx_private_sales_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `tier_benefits` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `tier` varchar(50) NOT NULL,
  `benefit_type` varchar(100) NOT NULL,
  `benefit_name` varchar(255) NOT NULL,
  `description` text DEFAULT NULL,
  `value` decimal(10,2) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_tier_benefits_tier` (`tier`),
  KEY `idx_tier_benefits_benefit_type` (`benefit_type`),
  KEY `idx_tier_benefits_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `investment_products` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `description` text DEFAULT NULL,
  `minimum_investment` decimal(15,2) NOT NULL,
  `expected_return` decimal(5,4) NOT NULL,
  `risk_level` varchar(50) NOT NULL,
  `tier_requirement` varchar(50) NOT NULL,
  `available` tinyint(1) DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_investment_products_tier_requirement` (`tier_requirement`),
  KEY `idx_investment_products_risk_level` (`risk_level`),
  KEY `idx_investment_products_available` (`available`),
  KEY `idx_investment_products_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `gift_card_types` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `brand` varchar(255) NOT NULL,
  `denomination` decimal(10,2) NOT NULL,
  `tier_requirement` varchar(50) NOT NULL,
  `exchange_rate_bonus` decimal(3,2) DEFAULT '1.00',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_gift_card_types_tier_requirement` (`tier_requirement`),
  KEY `idx_gift_card_types_brand` (`brand`),
  KEY `idx_gift_card_types_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `annual_user_statistics` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `year` int NOT NULL,
  `total_spending` decimal(15,2) DEFAULT '0.00',
  `purchases_count` int DEFAULT '0',
  `points_earned` int DEFAULT '0',
  `points_used` int DEFAULT '0',
  `tier_benefits_used` int DEFAULT '0',
  `exclusive_events_attended` int DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_annual_user_statistics_user_year` (`user_id`,`year`),
  KEY `idx_annual_user_statistics_year` (`year`),
  KEY `idx_annual_user_statistics_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `purchases` (
  `id` varchar(255) NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `store_id` bigint unsigned NOT NULL,
  `total_amount` decimal(10,2) NOT NULL,
  `points_earned` int DEFAULT '0',
  `tier_discount_applied` decimal(3,2) DEFAULT '0.00',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_purchases_user_id` (`user_id`),
  KEY `idx_purchases_store_id` (`store_id`),
  KEY `idx_purchases_created_at` (`created_at`),
  KEY `idx_purchases_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `concierge_requests` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `staff_id` bigint unsigned NOT NULL,
  `service_type` varchar(255) NOT NULL,
  `request_details` text DEFAULT NULL,
  `status` varchar(50) NOT NULL DEFAULT 'pending',
  `priority` varchar(50) NOT NULL DEFAULT 'medium',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_concierge_requests_user_id` (`user_id`),
  KEY `idx_concierge_requests_staff_id` (`staff_id`),
  KEY `idx_concierge_requests_status` (`status`),
  KEY `idx_concierge_requests_priority` (`priority`),
  KEY `idx_concierge_requests_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `event_registrations` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `event_id` bigint unsigned NOT NULL,
  `attendees_count` int DEFAULT '1',
  `special_requests` text DEFAULT NULL,
  `status` varchar(50) NOT NULL DEFAULT 'pending',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_event_registrations_user_id` (`user_id`),
  KEY `idx_event_registrations_event_id` (`event_id`),
  KEY `idx_event_registrations_status` (`status`),
  KEY `idx_event_registrations_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `premium_deliveries` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `order_id` varchar(255) NOT NULL,
  `delivery_type` varchar(100) NOT NULL,
  `time_slot` varchar(100) DEFAULT NULL,
  `special_instructions` text DEFAULT NULL,
  `status` varchar(50) NOT NULL DEFAULT 'pending',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_premium_deliveries_user_id` (`user_id`),
  KEY `idx_premium_deliveries_order_id` (`order_id`),
  KEY `idx_premium_deliveries_status` (`status`),
  KEY `idx_premium_deliveries_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `survey_templates` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `type` varchar(100) NOT NULL,
  `questions` json DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_survey_templates_type` (`type`),
  KEY `idx_survey_templates_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 外部キー制約の追加
ALTER TABLE `user_memberships` ADD CONSTRAINT `fk_user_memberships_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `user_memberships` ADD CONSTRAINT `fk_user_memberships_tier_id` FOREIGN KEY (`tier_id`) REFERENCES `membership_tiers` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE `user_profiles` ADD CONSTRAINT `fk_user_profiles_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE `user_activities` ADD CONSTRAINT `fk_user_activities_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE `notifications` ADD CONSTRAINT `fk_notifications_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE `user_preferences` ADD CONSTRAINT `fk_user_preferences_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE `user_scopes` ADD CONSTRAINT `fk_user_scopes_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `user_scopes` ADD CONSTRAINT `fk_user_scopes_scope_id` FOREIGN KEY (`scope_id`) REFERENCES `scopes` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE `user_roles` ADD CONSTRAINT `fk_user_roles_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `user_roles` ADD CONSTRAINT `fk_user_roles_role_id` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE `role_permissions` ADD CONSTRAINT `fk_role_permissions_role_id` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `role_permissions` ADD CONSTRAINT `fk_role_permissions_permission_id` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE `refresh_tokens` ADD CONSTRAINT `fk_refresh_tokens_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE `user_tokens` ADD CONSTRAINT `fk_user_tokens_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE `point_transactions` ADD CONSTRAINT `fk_point_transactions_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE `user_sessions` ADD CONSTRAINT `fk_user_sessions_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE `security_events` ADD CONSTRAINT `fk_security_events_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;

ALTER TABLE `rate_limit_logs` ADD CONSTRAINT `fk_rate_limit_logs_rule_id` FOREIGN KEY (`rule_id`) REFERENCES `rate_limit_rules` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `rate_limit_logs` ADD CONSTRAINT `fk_rate_limit_logs_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;

ALTER TABLE `device_fingerprints` ADD CONSTRAINT `fk_device_fingerprints_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE `user_suspensions` ADD CONSTRAINT `fk_user_suspensions_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `user_suspensions` ADD CONSTRAINT `fk_user_suspensions_suspended_by` FOREIGN KEY (`suspended_by`) REFERENCES `users` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE `user_approval_queue` ADD CONSTRAINT `fk_user_approval_queue_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `user_approval_queue` ADD CONSTRAINT `fk_user_approval_queue_assigned_to` FOREIGN KEY (`assigned_to`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;

ALTER TABLE `fraud_alerts` ADD CONSTRAINT `fk_fraud_alerts_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fraud_alerts` ADD CONSTRAINT `fk_fraud_alerts_resolved_by` FOREIGN KEY (`resolved_by`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;

ALTER TABLE `admin_actions` ADD CONSTRAINT `fk_admin_actions_admin_user_id` FOREIGN KEY (`admin_user_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE `system_settings_history` ADD CONSTRAINT `fk_system_settings_history_changed_by` FOREIGN KEY (`changed_by`) REFERENCES `users` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE `export_jobs` ADD CONSTRAINT `fk_export_jobs_created_by` FOREIGN KEY (`created_by`) REFERENCES `users` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE `announcements` ADD CONSTRAINT `fk_announcements_created_by` FOREIGN KEY (`created_by`) REFERENCES `users` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE `audit_logs` ADD CONSTRAINT `fk_audit_logs_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE `concurrent_sessions` ADD CONSTRAINT `fk_concurrent_sessions_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- プレミアム機能の外部キー制約
ALTER TABLE `concierge_requests` ADD CONSTRAINT `fk_concierge_requests_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `concierge_requests` ADD CONSTRAINT `fk_concierge_requests_staff_id` FOREIGN KEY (`staff_id`) REFERENCES `concierge_staff` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE `event_registrations` ADD CONSTRAINT `fk_event_registrations_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `event_registrations` ADD CONSTRAINT `fk_event_registrations_event_id` FOREIGN KEY (`event_id`) REFERENCES `vip_events` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE `premium_deliveries` ADD CONSTRAINT `fk_premium_deliveries_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE `purchases` ADD CONSTRAINT `fk_purchases_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE `annual_user_statistics` ADD CONSTRAINT `fk_annual_user_statistics_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE `webhook_deliveries` ADD CONSTRAINT `fk_webhook_deliveries_endpoint_id` FOREIGN KEY (`webhook_endpoint_id`) REFERENCES `webhook_endpoints` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- 追加のインデックス
CREATE INDEX `idx_users_name` ON `users` (`name`);
CREATE INDEX `idx_users_age` ON `users` (`age`);
CREATE INDEX `idx_users_created_at` ON `users` (`created_at`);

CREATE INDEX `idx_auths_is_active` ON `auths` (`is_active`);
CREATE INDEX `idx_auths_last_login_at` ON `auths` (`last_login_at`);

CREATE INDEX `idx_user_memberships_points` ON `user_memberships` (`points`);
CREATE INDEX `idx_user_memberships_total_spent` ON `user_memberships` (`total_spent`);
CREATE INDEX `idx_user_memberships_joined_at` ON `user_memberships` (`joined_at`);
CREATE INDEX `idx_user_memberships_last_activity_at` ON `user_memberships` (`last_activity_at`);
CREATE INDEX `idx_user_memberships_expires_at` ON `user_memberships` (`expires_at`);
CREATE INDEX `idx_user_memberships_is_active` ON `user_memberships` (`is_active`);

CREATE INDEX `idx_membership_tiers_is_active` ON `membership_tiers` (`is_active`);

CREATE INDEX `idx_point_transactions_created_at` ON `point_transactions` (`created_at`);
CREATE INDEX `idx_point_transactions_expires_at` ON `point_transactions` (`expires_at`);
CREATE INDEX `idx_point_transactions_reference_type` ON `point_transactions` (`reference_type`);

CREATE INDEX `idx_user_activities_created_at` ON `user_activities` (`created_at`);

CREATE INDEX `idx_notifications_is_read` ON `notifications` (`is_read`);
CREATE INDEX `idx_notifications_read_at` ON `notifications` (`read_at`);
CREATE INDEX `idx_notifications_created_at` ON `notifications` (`created_at`);

CREATE INDEX `idx_login_attempts_success` ON `login_attempts` (`success`);
CREATE INDEX `idx_login_attempts_created_at` ON `login_attempts` (`created_at`);

CREATE INDEX `idx_ip_blacklists_is_active` ON `ip_blacklists` (`is_active`);
CREATE INDEX `idx_ip_blacklists_expires_at` ON `ip_blacklists` (`expires_at`);

CREATE INDEX `idx_user_sessions_expires_at` ON `user_sessions` (`expires_at`);
CREATE INDEX `idx_user_sessions_is_active` ON `user_sessions` (`is_active`);

CREATE INDEX `idx_security_events_created_at` ON `security_events` (`created_at`);

CREATE INDEX `idx_rate_limit_rules_resource` ON `rate_limit_rules` (`resource`);
CREATE INDEX `idx_rate_limit_rules_is_active` ON `rate_limit_rules` (`is_active`);

CREATE INDEX `idx_rate_limit_logs_window_start` ON `rate_limit_logs` (`window_start`);
CREATE INDEX `idx_rate_limit_logs_window_end` ON `rate_limit_logs` (`window_end`);
CREATE INDEX `idx_rate_limit_logs_blocked` ON `rate_limit_logs` (`blocked`);

CREATE INDEX `idx_device_fingerprints_is_trusted` ON `device_fingerprints` (`is_trusted`);
CREATE INDEX `idx_device_fingerprints_last_seen_at` ON `device_fingerprints` (`last_seen_at`);

CREATE INDEX `idx_user_suspensions_suspended_at` ON `user_suspensions` (`suspended_at`);
CREATE INDEX `idx_user_suspensions_expires_at` ON `user_suspensions` (`expires_at`);
CREATE INDEX `idx_user_suspensions_status` ON `user_suspensions` (`status`);

CREATE INDEX `idx_fraud_alerts_triggered_at` ON `fraud_alerts` (`triggered_at`);
CREATE INDEX `idx_fraud_alerts_resolved_at` ON `fraud_alerts` (`resolved_at`);

CREATE INDEX `idx_admin_actions_created_at` ON `admin_actions` (`created_at`);

CREATE INDEX `idx_system_settings_history_created_at` ON `system_settings_history` (`created_at`);

CREATE INDEX `idx_export_jobs_created_at` ON `export_jobs` (`created_at`);
CREATE INDEX `idx_export_jobs_updated_at` ON `export_jobs` (`updated_at`);

CREATE INDEX `idx_announcements_created_at` ON `announcements` (`created_at`);
CREATE INDEX `idx_announcements_updated_at` ON `announcements` (`updated_at`);

CREATE INDEX `idx_audit_logs_created_at` ON `audit_logs` (`created_at`);

-- プレミアム機能のインデックス
CREATE INDEX `idx_products_price` ON `products` (`price`);
CREATE INDEX `idx_products_stock` ON `products` (`stock`);
CREATE INDEX `idx_products_limited_quantity` ON `products` (`limited_quantity`);
CREATE INDEX `idx_products_created_at` ON `products` (`created_at`);

CREATE INDEX `idx_vip_events_capacity` ON `vip_events` (`capacity`);
CREATE INDEX `idx_vip_events_current_attendees` ON `vip_events` (`current_attendees`);

CREATE INDEX `idx_concierge_staff_rating` ON `concierge_staff` (`rating`);

CREATE INDEX `idx_designers_hourly_rate` ON `designers` (`hourly_rate`);

CREATE INDEX `idx_private_sales_invitation_only` ON `private_sales` (`invitation_only`);
CREATE INDEX `idx_private_sales_discount_rate` ON `private_sales` (`discount_rate`);

CREATE INDEX `idx_investment_products_minimum_investment` ON `investment_products` (`minimum_investment`);
CREATE INDEX `idx_investment_products_expected_return` ON `investment_products` (`expected_return`);

CREATE INDEX `idx_gift_card_types_denomination` ON `gift_card_types` (`denomination`);
CREATE INDEX `idx_gift_card_types_exchange_rate_bonus` ON `gift_card_types` (`exchange_rate_bonus`);

CREATE INDEX `idx_purchases_total_amount` ON `purchases` (`total_amount`);
CREATE INDEX `idx_purchases_points_earned` ON `purchases` (`points_earned`);
CREATE INDEX `idx_purchases_tier_discount_applied` ON `purchases` (`tier_discount_applied`);

CREATE INDEX `idx_concierge_requests_created_at` ON `concierge_requests` (`created_at`);

CREATE INDEX `idx_event_registrations_attendees_count` ON `event_registrations` (`attendees_count`);

CREATE INDEX `idx_premium_deliveries_delivery_type` ON `premium_deliveries` (`delivery_type`);

-- 複合インデックス
CREATE INDEX `idx_users_email_deleted_at` ON `users` (`email`, `deleted_at`);
CREATE INDEX `idx_auths_email_is_active` ON `auths` (`email`, `is_active`);
CREATE INDEX `idx_user_memberships_user_tier_active` ON `user_memberships` (`user_id`, `tier_id`, `is_active`);
CREATE INDEX `idx_point_transactions_user_type_created` ON `point_transactions` (`user_id`, `type`, `created_at`);
CREATE INDEX `idx_notifications_user_read_created` ON `notifications` (`user_id`, `is_read`, `created_at`);
CREATE INDEX `idx_login_attempts_email_success_created` ON `login_attempts` (`email`, `success`, `created_at`);
CREATE INDEX `idx_security_events_user_severity_created` ON `security_events` (`user_id`, `severity`, `created_at`);
CREATE INDEX `idx_fraud_alerts_user_status_triggered` ON `fraud_alerts` (`user_id`, `status`, `triggered_at`);
CREATE INDEX `idx_products_category_tier_price` ON `products` (`category`, `exclusive_tier`, `price`);
CREATE INDEX `idx_vip_events_tier_date` ON `vip_events` (`tier_requirement`, `event_date`);
CREATE INDEX `idx_concierge_requests_user_status_created` ON `concierge_requests` (`user_id`, `status`, `created_at`);
