WITH inserted_link AS (
    INSERT INTO links (user_id, code, original_url)
        VALUES ($1, $2, $3)
        RETURNING id, user_id, created_at
),
     inserted_analytics AS (
         INSERT INTO analytics (link_id, clicks)
             SELECT id, 0 FROM inserted_link
     )
SELECT id, user_id, created_at FROM inserted_link;
