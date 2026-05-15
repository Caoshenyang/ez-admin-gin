-- +migrate Down
DELETE FROM casbin_rule WHERE v0 = 'super_admin' AND v1 LIKE '/api/v1/system/notifications%';
