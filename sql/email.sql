USE dwp_18_dev_00_00_00;

SELECT
  DISTINCT (email) AS email
FROM section AS s
  INNER JOIN app_user AS a ON s.app_user_id = a.id
WHERE s.drawing_id = ?
