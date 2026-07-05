SELECT id, user_id, code, original_url, created_at
FROM links
WHERE user_id = $1
ORDER BY created_at DESC;
