DROP SCHEMA IF EXISTS dwp_18_dev_00_00_00;
CREATE SCHEMA IF NOT EXISTS dwp_18_dev_00_00_00;
USE dwp_18_dev_00_00_00;

-- tables
-- Table: app_user
CREATE TABLE app_user (
  id      int          NOT NULL AUTO_INCREMENT,
  email   varchar(250) NOT NULL,
  created datetime     NOT NULL DEFAULT NOW(),
  name    varchar(250) NOT NULL,
  updated datetime     NOT NULL DEFAULT NOW(),
  UNIQUE INDEX app_user_ak_1 (email),
  CONSTRAINT app_user_pk PRIMARY KEY (id)
);

-- Table: drawing
CREATE TABLE drawing (
  id        int          NOT NULL AUTO_INCREMENT,
  url       varchar(250) NULL,
  created   datetime     NOT NULL DEFAULT NOW(),
  updated   datetime     NOT NULL DEFAULT NOW(),
  completed datetime     NULL,
  CONSTRAINT drawing_pk PRIMARY KEY (id)
);

-- Table: section
CREATE TABLE section (
  drawing_id  int          NOT NULL,
  type        varchar(10)  NOT NULL,
  app_user_id int          NOT NULL,
  url         varchar(250) NOT NULL,
  created     datetime     NOT NULL,
  updated     datetime     NOT NULL,
  CONSTRAINT section_pk PRIMARY KEY (drawing_id, type)
);

-- foreign keys
-- Reference: image_app_user (table: section)
ALTER TABLE section
  ADD CONSTRAINT image_app_user FOREIGN KEY image_app_user (app_user_id)
REFERENCES app_user (id);

-- Reference: image_drawing (table: section)
ALTER TABLE section
  ADD CONSTRAINT image_drawing FOREIGN KEY image_drawing (drawing_id)
REFERENCES drawing (id);

-- End of file.

