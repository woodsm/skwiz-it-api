USE dwp_18_dev_00_00_00;
INSERT INTO section (drawing_id, type, app_user_id, url, created, updated)
VALUES (
  (SELECT id
   FROM drawing
   WHERE completed IS NULL
         AND id NOT IN (SELECT drawing_id
                        FROM section
                        WHERE type = ?)
   LIMIT 1)
  , ?, ?, ?, NOW(), NOW())
