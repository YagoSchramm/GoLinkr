UPDATE analytics
SET clicks = clicks + 1
WHERE link_id = $1
RETURNING id, link_id, clicks;