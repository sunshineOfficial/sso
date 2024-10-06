select id        as ID,
       email     as Email,
       pass_hash as PassHash
from users
where email = $1;