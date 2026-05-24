SELECT
    COALESCE(a.clicks, 0) AS clicks
FROM links l
         LEFT JOIN analytics a ON a.link_id = l.id
WHERE l.id = $1;