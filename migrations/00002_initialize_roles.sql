-- +migrate Up
INSERT INTO roles
(id, `type`, created_at, updated_at)
VALUES(1, 'USER', '2019-04-01 00:00:00', '2019-04-01 00:00:00');
INSERT INTO roles
(id, `type`, created_at, updated_at)
VALUES(2, 'ORGANIZATION_MEMBER', '2019-04-01 00:00:00', '2019-04-01 00:00:00');
INSERT INTO roles
(id, `type`, created_at, updated_at)
VALUES(3, 'ORGANIZATION_ADMIN', '2019-04-01 00:00:00', '2019-04-01 00:00:00');
INSERT INTO roles
(id, `type`, created_at, updated_at)
VALUES(4, 'SUPER_ADMIN', '2019-04-01 00:00:00', '2019-04-01 00:00:00');

-- +migrate Down
TRUNCATE roles;
