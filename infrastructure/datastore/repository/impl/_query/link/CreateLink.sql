WITH inserted_link AS (
    INSERT INTO links (code, original_url)
        VALUES ($1, $2)
        RETURNING id, created_at
),
     inserted_analytics AS (
         INSERT INTO analytics (link_id, clicks)
             SELECT id, 0 FROM inserted_link
     )
SELECT id, created_at FROM inserted_link;