WITH authorized_link AS (
    SELECT id
    FROM links
    WHERE id = $1
      AND user_id = $2
),
     per_month_week AS (
    SELECT
        ((EXTRACT(DAY FROM ac.clicked_at)::INT - 1) / 7) + 1 AS week_of_month,
        DATE_TRUNC('month', ac.clicked_at) AS clicked_month,
        COUNT(*)::FLOAT AS clicks
    FROM analytics_clicks ac
             JOIN authorized_link al ON al.id = ac.link_id
    GROUP BY week_of_month, clicked_month
)
SELECT
    weeks.week_of_month,
    COALESCE(AVG(per_month_week.clicks), 0) AS average_clicks
FROM authorized_link
         CROSS JOIN generate_series(1, 5) AS weeks(week_of_month)
         LEFT JOIN per_month_week ON per_month_week.week_of_month = weeks.week_of_month
GROUP BY weeks.week_of_month
ORDER BY weeks.week_of_month;
