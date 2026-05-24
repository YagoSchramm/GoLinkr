SELECT
    l.id AS link_id,
    l.code,
    l.original_url,
    l.created_at,
    a.clicks
FROM links l
         INNER JOIN analytics a ON a.link_id = l.id
WHERE l.id = $1;