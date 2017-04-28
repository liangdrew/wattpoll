CREATE TABLE polls(
  question VARCHAR(255),
  story_id VARCHAR(255) NOT NULL,
  part_id VARCHAR(255) NOT NULL,
  created TIMESTAMP,
  duration INT NOT NULL,
  PRIMARY KEY (part_id)
);

CREATE TABLE choices(
  id INT NOT NULL AUTO_INCREMENT,
  part_id VARCHAR(255) NOT NULL,
  choice VARCHAR(255),
  choice_index TINYINT NOT NULL,
  votes INT,
  PRIMARY KEY (id),
  FOREIGN KEY (part_id)
	REFERENCES polls(part_id)
	ON DELETE CASCADE
);

CREATE TABLE votes(
  id INT NOT NULL AUTO_INCREMENT,
  part_id VARCHAR(255) NOT NULL,
  username VARCHAR(255),
  choice_index INT NOT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (part_id)
	REFERENCES polls(part_id)
	ON DELETE CASCADE
);
