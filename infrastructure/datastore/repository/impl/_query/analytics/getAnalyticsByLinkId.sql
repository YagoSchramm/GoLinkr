SELECT
    a.id,
    l.id,
    a.clicks
FROM links l
         JOIN analytics a ON a.link_id = l.id
WHERE l.id = $1
  AND l.user_id = $2;
