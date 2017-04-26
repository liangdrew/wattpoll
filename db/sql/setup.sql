CREATE TABLE poll(
id INT NOT NULL AUTO_INCREMENT,
question VARCHAR(255),
story_id VARCHAR(255) NOT NULL,
part_id VARCHAR(255) NOT NULL,
create_time TIMESTAMP,
PRIMARY KEY (id)
);

CREATE TABLE choice(
id INT NOT NULL AUTO_INCREMENT,
poll_id INT NOT NULL,
choice_text VARCHAR(255),
choice_index TINYINT NOT NULL,
votes INT,
PRIMARY KEY (id),
FOREIGN KEY (poll_id)
	REFERENCES poll(id)
	ON DELETE CASCADE
);

CREATE TABLE vote(
id INT NOT NULL AUTO_INCREMENT,
poll_id iNT NOT NULL,
username VARCHAR(255),
choice_id INT NOT NULL,
PRIMARY KEY (id),
FOREIGN KEY (poll_id)
	REFERENCES poll(id)
	ON DELETE CASCADE
);