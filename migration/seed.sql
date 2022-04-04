-- CLASS
INSERT INTO class(name, start_date, end_date,capacity) VALUES  ('pilates','2021-12-01','2021-12-20', 20);
INSERT INTO class(name, start_date, end_date,capacity) VALUES  ('crossfit','2021-12-01','2021-12-30', 40);
INSERT INTO class(name, start_date, end_date,capacity) VALUES  ('jiu-jitsu','2022-04-04','2022-04-20', 30);

-- MEMBER
INSERT INTO member(name) VALUES  ('Alice');
INSERT INTO member(name) VALUES  ('Bob');
INSERT INTO member(name) VALUES  ('Charlie');

-- BOOKING
INSERT INTO booking(class_id, member_id, date) VALUES (1,1,'2021-12-01');
INSERT INTO booking(class_id, member_id, date) VALUES (2,1,'2021-12-02');
INSERT INTO booking(class_id, member_id, date) VALUES (1,3,'2021-12-01');
INSERT INTO booking(class_id, member_id, date) VALUES (1,3,'2021-12-02');
INSERT INTO booking(class_id, member_id, date) VALUES (1,3,'2021-12-03');
