-- name: GetProvidersList :many
SELECT * FROM providers
ORDER BY name;

-- name: GetProvidersDetails :one
SELECT * FROM providers
WHERE name = ? LIMIT 1;

-- name: GetSearchProviders :many
SELECT * FROM providers
WHERE name LIKE CONCAT('%', ?, '%');

-- name: GetVouchersList :many
SELECT * FROM vouchers
JOIN voucher_groups ON vouchers.group_id = voucher_groups.id
ORDER BY vouchers.id;

-- name: StateListUniqueVouchers :many
SELECT DISTINCT state
FROM vouchers
ORDER BY state ASC;

-- name: UpdateVouchersPrintedField :exec
UPDATE vouchers
set printed = 'printed'
WHERE id = ?;

-- name: PrintedListUniqueVouchers :many
SELECT DISTINCT printed
FROM vouchers
ORDER BY state ASC;

-- name: ValidityListUniqueVouchers :many
SELECT DISTINCT validity
FROM vouchers
ORDER BY validity ASC;

-- name: GroupNameListUniqueVouchers :many
SELECT DISTINCT vg.name
FROM voucher_groups vg
INNER JOIN vouchers v ON v.group_id = vg.id
WHERE (:validity IS NULL OR v.validity = :validity)
ORDER BY vg.name ASC;

-- name: ListVouchers :many
SELECT
    v.id,
    v.username,
    v.validity,
    v.state,
    v.printed,
    vg.name AS group_name,
    p.name AS provider_name
FROM vouchers v
JOIN voucher_groups vg ON vg.id = v.group_id
JOIN providers p ON p.id = vg.provider_id
WHERE (:state IS NULL OR v.state = :state)
  AND (:validity IS NULL OR v.validity = :validity)
  AND (:printed IS NULL OR v.printed = :printed)
  AND (:group_name IS NULL OR vg.name = :group_name)
  AND (
    :search IS NULL OR v.id IN (
      SELECT rowid FROM vouchers_fts WHERE vouchers_fts.username MATCH :search || '*'
    )
  )
LIMIT :limit OFFSET :offset;

-- name: CountVouchers :one
SELECT COUNT(*) FROM vouchers v
JOIN voucher_groups vg ON vg.id = v.group_id
JOIN providers p ON p.id = vg.provider_id
WHERE (:state IS NULL OR v.state = :state)
  AND (:validity IS NULL OR v.validity = :validity)
  AND (:printed IS NULL OR v.printed = :printed)
  AND (:group_name IS NULL OR vg.name = :group_name)
  AND (
    :search IS NULL OR v.id IN (
      SELECT rowid FROM vouchers_fts WHERE vouchers_fts.username MATCH :search || '*'
    )
  );

-- name: GetVoucherByID :one
SELECT
    v.id,
    v.username,
    v.validity,
    v.state,
    v.printed,
    vg.name AS group_name,
    p.name AS provider_name
FROM vouchers v
JOIN voucher_groups vg ON vg.id = v.group_id
JOIN providers p ON p.id = vg.provider_id
WHERE v.id = ? LIMIT 1;

-- name: GetExpiredVouchersPerDay :many
SELECT
    strftime('%w', datetime(endtime, 'unixepoch')) as day_of_week,
    validity,
    COUNT(*) as count
FROM vouchers
WHERE state = 'expired'
  AND endtime >= ? AND endtime <= ?
GROUP BY day_of_week, validity
ORDER BY day_of_week, validity;

-- name: GetExpiredVouchersPerWeek :many
SELECT
    strftime('%Y-%W', datetime(endtime, 'unixepoch')) as week,
    validity,
    COUNT(*) as count
FROM vouchers
WHERE state = 'expired'
  AND endtime >= ? AND endtime <= ?
GROUP BY week, validity
ORDER BY week, validity;

-- name: GetExpiredVouchersPerMonth :many
SELECT
    strftime('%Y-%m', datetime(endtime, 'unixepoch')) as month,
    validity,
    COUNT(*) as count
FROM vouchers
WHERE state = 'expired'
  AND endtime >= ? AND endtime <= ?
GROUP BY month, validity
ORDER BY month, validity;