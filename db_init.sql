create table todos (
    ID            integer primary key asc,
    name          text not null,
    content       text not null,
    resolved      boolean not null
);

insert into todos (id, name, content, resolved) values (0, 'Push changes', 'Push changes to branch before leaving the office.', 0);
insert into todos (id, name, content, resolved) values (1, 'Do dishes', 'Cleanup the kitchen.', 1);
