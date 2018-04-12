USE dwp_18_dev_00_00_00;

SELECT
  d.id               AS drawing_id,
  d.url              AS drawing_url,
  d.created          AS drawing_created,
  d.updated          AS drawing_updated,
  d.completed        AS drawing_completed,

  MAX(CASE WHEN s.type = 'top'
    THEN s.url
      ELSE NULL END) AS section_top_url,

  MAX(CASE WHEN s.type = 'top'
    THEN a.name
      ELSE NULL END) AS section_top_name,

  MAX(CASE WHEN s.type = 'top'
    THEN a.email
      ELSE NULL END) AS section_top_email,

  MAX(CASE WHEN s.type = 'middle'
    THEN s.url
      ELSE NULL END) AS section_middle_url,

  MAX(CASE WHEN s.type = 'middle'
    THEN a.name
      ELSE NULL END) AS section_middle_name,

  MAX(CASE WHEN s.type = 'middle'
    THEN a.email
      ELSE NULL END) AS section_middle_email,

  MAX(CASE WHEN s.type = 'bottom'
    THEN s.url
      ELSE NULL END) AS section_bottom_url,

  MAX(CASE WHEN s.type = 'bottom'
    THEN a.name
      ELSE NULL END) AS section_bottom_name,

  MAX(CASE WHEN s.type = 'bottom'
    THEN a.email
      ELSE NULL END) AS section_bottom_email

FROM drawing AS d
  LEFT JOIN section s ON d.id = s.drawing_id
  LEFT JOIN app_user a ON s.app_user_id = a.id
WHERE
  d.id <> 0
--   AND d.id = ?
--   AND a.id = ?
GROUP BY
  d.id,
  d.url,
  d.created,
  d.updated,
  d.completed;

INSERT INTO drawing (id, url, created, updated, completed) VALUES (0, NULL, NOW(), NOW(), NULL);

