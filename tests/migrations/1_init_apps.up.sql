insert into apps (name, secret)
values ('test', 'test-secret')
on conflict do nothing;