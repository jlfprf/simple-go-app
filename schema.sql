-- Simple users table for testing
create table users (id serial not null primary key, name text not null unique, hashedpass text not null);
insert into users values (default, 'jlf', '$2a$10$Bxol2j3PpKNJtSP5rYQuF.Txy.GWBL1KgaNYAVrdHZYIy9wd.Kogi');
select name, hashedpass from users where name = 'jlf' order by name desc;

-- sessions table to save session information
create table sessions (sessionid text primary key not null, name text not null)