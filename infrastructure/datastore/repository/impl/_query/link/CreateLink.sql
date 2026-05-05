		INSERT INTO links (code, original_url)
		VALUES ($1, $2)
		RETURNING id, created_at
