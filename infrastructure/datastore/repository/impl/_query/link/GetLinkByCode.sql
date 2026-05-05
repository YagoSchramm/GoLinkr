		SELECT id, code, original_url, created_at
		FROM links
		WHERE code = $1
