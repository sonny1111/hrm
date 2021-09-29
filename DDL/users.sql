CREATE TABLE users(
    user_id BIGSERIAL PRIMARY KEY,
    first_name VARCHAR(25) NOT NULL,
    last_name VARCHAR(25) NOT NULL,
    middle_name VARCHAR(25),
    username VARCHAR(15) UNIQUE NOT NULL,
    password VARCHAR(25),
    group_id INT NULL,
    role_id INT NULL,
    -- expires VARCHAR,
    -- attemps NUMERIC,
    -- days_b4_expn NUMERIC
);

CREATE TABLE groups(
    group_id BIGSERIAL PRIMARY KEY,
    group_name VARCHAR(12) UNIQUE NOT NULL,
    description VARCHAR(15) NULL,
);

--User can be assigned to one and only one group. But one group can be assigned to many users
ALTER TABLE groups ADD CONSTRAINT usr_grp_grpid_fk FOREIGN KEY(group_id) REFERENCES 
 groups(group_id) ON DELETE CASCADE;

CREATE TABLE roles(
     role_id BIGSERIAL PRIMARY KEY,
    role_name VARCHAR(12) UNIQUE NOT NULL,
    description VARCHAR(15)
);
--User can be assigned one and only one role. But one role can be assigned to many users
ALTER TABLE roles ADD CONSTRAINT usr_rl_rid_fk FOREIGN KEY(role_id) REFERENCES roles(role_id) 
ON DELETE CASCADE;

CREATE TABLE privileges(
    privilege_id BIGSERIAL PRIMARY KEY,
    privilege_name VARCHAR(12) UNIQUE NOT NULL,
    description VARCHAR(15)
);


--More than one privilege can be assigned to a role. Many roles can have the same privilege.
CREATE TABLE role_privilges(
    role_id INT,
    priv_id INT,
);

ALTER TABLE role_privilges ADD CONSTRAINT rl_priv_rid_fk FOREIGN KEY(role_id) REFERENCES roles(id)
ON DELETE CASCADE ON UPDATE SET NULL;
ALTER TABLE role_privilges ADD CONSTRAINT rl_priv_pid_fk FOREIGN KEY(priv_id) REFERENCES privileges(id)
ON DELETE CASCADE ON UPDATE SET NULL;

CREATE TABLE group_roles(
    group_id INT,
    role_id INT,
);

ALTER TABLE group_roles ADD CONSTRAINT grp_rl_grpid_fk FOREIGN KEY(group_id) REFERENCES groups(group_id)
ON DELETE CASCADE;
ALTER TABLE group_roles ADD CONSTRAINT grp_rl_rid_fk FOREIGN KEY(role_id) REFERENCES roles(id)
ON DELETE CASCADE;
