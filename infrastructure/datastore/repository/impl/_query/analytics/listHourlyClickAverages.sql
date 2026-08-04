WITH authorized_link AS (
    SELECT id
    FROM links
    WHERE id = $1
      AND user_id = $2
),
     per_day_hour AS (
    SELECT
        EXTRACT(HOUR FROM ac.clicked_at)::INT AS hour,
        DATE(ac.clicked_at) AS clicked_date,
        COUNT(*)::FLOAT AS clicks
    FROM analytics_clicks ac
             JOIN authorized_link al ON al.id = ac.link_id
    GROUP BY hour, clicked_date
)
SELECT
    hours.hour,
    COALESCE(AVG(per_day_hour.clicks), 0) AS average_clicks
FROM authorized_link
         CROSS JOIN generate_series(0, 23) AS hours(hour)
         LEFT JOIN per_day_hour ON per_day_hour.hour = hours.hour
GROUP BY hours.hour
ORDER BY hours.hour;
