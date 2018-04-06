USE dwp_18_dev_00_00_00;
INSERT INTO section (drawing_id, type, app_user_id, url, created, updated)
VALUES (?, ?, ?, ?, NOW(), NOW());

SELECT
  MIN(qty) AS qty,
  type     AS type
FROM (
       SELECT
         COUNT(drawing_id) AS qty,
         'top'             AS type
       FROM section
       WHERE type = 'top'

       UNION

       SELECT
         COUNT(drawing_id) AS qty,
         'middle'          AS type
       FROM section
       WHERE type = 'middle'

       UNION

       SELECT
         COUNT(drawing_id) AS qty,
         'bottom'          AS type
       FROM section
       WHERE type = 'bottom'
     ) AS smry
GROUP BY type
ORDER BY qty ASC
LIMIT 1
