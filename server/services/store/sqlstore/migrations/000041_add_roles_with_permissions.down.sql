ALTER TABLE {{.prefix}}users DROP COLUMN role_id;

DROP TABLE {{.prefix}}role_permissions;
DROP TABLE {{.prefix}}permissions;
DROP TABLE {{.prefix}}roles;
