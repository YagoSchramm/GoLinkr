WITH authorized_link AS (
    SELECT id
    FROM links
    WHERE id = $1
      AND user_id = $2
),
     per_date_weekday AS (
    SELECT
        EXTRACT(ISODOW FROM ac.clicked_at)::INT AS day_of_week,
        DATE(ac.clicked_at) AS clicked_date,
        COUNT(*)::FLOAT AS clicks
    FROM analytics_clicks ac
             JOIN authorized_link al ON al.id = ac.link_id
    GROUP BY day_of_week, clicked_date
)
SELECT
    weekdays.day_of_week,
    COALESCE(AVG(per_date_weekday.clicks), 0) AS average_clicks
FROM authorized_link
         CROSS JOIN generate_series(1, 7) AS weekdays(day_of_week)
         LEFT JOIN per_date_weekday ON per_date_weekday.day_of_week = weekdays.day_of_week
GROUP BY weekdays.day_of_week
ORDER BY weekdays.day_of_week;
