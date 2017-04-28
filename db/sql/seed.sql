INSERT INTO polls (question, story_id, part_id, created, duration)
VALUES
	('What is your favorite letter?', '107369735', '1', '19731230153000', 1),
	('What index is the first element in an array?', '107369735', '2', '20171230153000', 2);

INSERT INTO choices (part_id, choice, choice_index, votes)
VALUES
	(1, 'A', 1, 10),
	(1, 'B', 2, 30),
	(1, 'C', 3, 20),
	(2, '0', 1, 10),
	(2, '1', 2, 20);

INSERT INTO votes (part_id, username, choice_id)
VALUES
	(1, 'cynthiashu', 2),
	(1, 'asdf', 3),
	(2, 'cynthiashu', 4),
	(2, 'cynthiashu8', 5);
