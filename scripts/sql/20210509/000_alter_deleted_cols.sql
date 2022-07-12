update linux_gey_db set deleted = null, deleted_by = null;

alter table linux_gey_db 
alter column deleted type timestamp
using deleted::timestamp without time zone;

alter table linux_gey_db
drop column deleted_date;

alter table linux_gey_db                
add constraint deleted_info_check
CHECK (
    (deleted is null and deleted_by is null) or
    (deleted is not null and deleted_by is not null)
);
