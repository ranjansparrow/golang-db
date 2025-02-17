drop table if exists album;

create table album (
                       id INT primary key AUTO_INCREMENT not null,
                       title varchar(128) not null,
                       artist varchar(255) not null,
                       price decimal(5,2) not null
);

insert into album (title, artist, price)
values
    ('Another Brick in the Wall', 'Pink Floyd', 9.99),
    ('The Dark Side of the Moon', 'Pink Floyd', 19.99),
    ('Back in Black', 'AC/DC', 19.99),
    ('The Bodyguard', 'Whitney Houston', 15.99),
    ('Bat Out of Hell', 'Meat Loaf', 15.99),
    ('Their Greatest Hits (1971-1975)', 'Eagles', 19.99),
    ('Saturday Night Fever', 'Bee Gees', 15.99),
    ('Rumours', 'Fleetwood Mac', 19.99),
    ('Grease', 'Various Artists', 15.99),
    ('The Joshua Tree', 'U2', 19.99),
    ('Thriller', 'Michael Jackson', 15.99);
