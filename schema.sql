create table users (name text not null primary key, hashedpass text not null);
insert into users (name, hashedpass) values ("jlf", "hash");
select name, hashedpass from users where name = "jlf" order by name desc;

create table sessions (cookie text primary key not null, name text not null)