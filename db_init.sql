create table todos (
    ID            integer primary key asc,
    subject       text not null,
    details       text not null,
    priority      integer not null,
    due_to_date   text not null,
    resolved      boolean not null
);

insert into todos (id, subject, details, priority, due_to_date, resolved) values (0, 'Push changes', 'Push changes to branch before leaving the office.', 5, '2024-12-01', 0);
insert into todos (id, subject, details, priority, due_to_date, resolved) values (1, 'Do dishes', 'Cleanup the kitchen.', 2, '2024-12-01', 1);
