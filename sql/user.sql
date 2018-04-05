USE dwp_18_dev_00_00_00;
INSERT INTO app_user (email, created, name, updated) VALUES (?, NOW(), ?, NOW())
ON DUPLICATE KEY UPDATE updated = VALUES(updated)
