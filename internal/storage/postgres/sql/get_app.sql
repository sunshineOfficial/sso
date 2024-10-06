select id     as ID,
       name   as Name,
       secret as Secret
from apps
where id = $1;