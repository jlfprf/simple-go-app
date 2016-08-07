-- Simple users table for testing
create table users (id serial not null primary key, name text not null unique, hashedpass text not null);
insert into users values (default, 'jlf', '243261243130244c7150432f35396337634d757a796f4e5a6c6475334f41744c49586f6f6849336545594c516b457a74774c2f536665536e36747065');
select name, hashedpass from users where name = 'jlf' order by name desc;

-- sessions table to save session information
create table sessions (sessionid text primary key not null, name text not null)