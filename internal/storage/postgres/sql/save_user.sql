insert into users (email, pass_hash)
values ($1, $2)
returning id;