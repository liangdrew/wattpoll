INSERT INTO polls (question, story_id, part_id, created)
VALUES
	('What is your favorite letter?', '107369735', '404210726', '19731230153000'),
	('What index is the first element in an array?', '107369735', '404210777', '19731230153000');

INSERT INTO choices (part_id, choice, choice_index, votes)
VALUES
	(1, 'A', 0, 10),
	(1, 'B', 2, 30),
	(1, 'C', 1, 20),
	(2, '0', 0, 10),
	(2, '1', 1, 20);

INSERT INTO votes (part_id, username, choice_id)
VALUES
	(1, 'cynthiashu', 2),
	(1, '', 3),
	(2, 'cynthiashu', 4),
	(2, 'cynthiashu8', 5);
