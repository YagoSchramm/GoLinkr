WITH updated_analytics AS (
    UPDATE analytics
    SET clicks = clicks + 1
    WHERE link_id = $1
    RETURNING id, link_id, clicks
),
     inserted_click AS (
         INSERT INTO analytics_clicks (link_id)
             SELECT link_id FROM updated_analytics
             RETURNING id
     )
SELECT updated_analytics.id, updated_analytics.link_id, updated_analytics.clicks
FROM updated_analytics,
     inserted_click;
