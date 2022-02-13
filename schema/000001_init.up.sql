create table users
(
    id                 serial                                                                                   not null
        constraint users_pkey
            primary key,
    created_at         timestamp    default now()                                                               not null,
    updated_at         timestamp    default now()                                                               not null,
    phone              varchar(255) default ''::character varying                                               not null,
    email              varchar(255) default ''::character varying                                               not null,
    role               text         default 'USER'::text                                                        not null,
    last_name          varchar(255) default ''::character varying                                               not null,
    first_name         varchar(255) default ''::character varying                                               not null,
    middle_name        varchar(255) default ''::character varying                                               not null,
    registered_manager boolean      default false                                                               not null,
    code               varchar(255) default ''::character varying                                               not null,
    image              varchar(255) default ''::character varying not null,
    token              varchar(255) default ''::character varying                                               not null
);

create unique index users_email_uindex
    on users (email);

create unique index users_id_uindex
    on users (id);

create unique index users_phone_uindex
    on users (phone);


create table project
(
    id                      serial                                                                                 not null
        constraint project_pk
            primary key,
    address                 varchar(255) default ''::character varying                                             not null,
    user_id                 integer      default 1                                                                 not null
        constraint project_users_id_fk
            references users,
    start_date              varchar(255) default ''::character varying                                             not null,
    end_date                varchar(255) default ''::character varying                                             not null,
    active                  boolean      default true,
    maker_id                integer      default 0,
    material_costs_over_all integer      default 0                                                                 not null,
    work_costs_over_all     integer      default 0                                                                 not null,
    material_cost_spent     integer      default 0                                                                 not null,
    work_cost_spent         integer      default 0                                                                 not null,
    image                   varchar(255) default ''::character varying not null,
    file                    varchar(255) default ''::character varying                                             not null
);


create unique index project_id_uindex
    on project (id);


create table stages
(
    id         serial                                     not null
        constraint stages_pk
            primary key,
    name       varchar(255) default ''::character varying not null,
    project_id integer
        constraint stages_project_id_fk
            references project,
    phase      boolean      default false                 not null,
    date       varchar(255) default ''::character varying not null
);

create unique index stages_id_uindex
    on stages (id);


create table cards
(
    id                 serial                                                         not null
        constraint card_pk
            primary key,
    title              varchar(255) default ''::character varying                     not null,
    deadline           varchar(255) default ''::character varying                     not null,
    stages_id          integer                                                        not null
        constraint card_stages_id_fk
            references stages,
    rating             varchar(255) default ''::character varying                     not null,
    description        varchar(255) default ''::character varying                     not null,
    status             varchar(255) default ''::character varying not null,
    state              varchar(255) default ''::character varying                     not null,
    left_to_pay        integer      default 0                                         not null,
    spent_on_materials integer      default 0                                         not null,
    images             text[]       default ARRAY []::text[]                          not null
);

create unique index card_id_uindex
    on cards (id);


create table tasks
(
    id       serial                                not null
        constraint task_pk
            primary key,
    title    text    default ''::character varying not null,
    card_id  integer                               not null
        constraint task_card_id_fk
            references cards,
    complete boolean default false,
    images   text[]  default ARRAY []::text[]      not null,
    length   bigint  default 0                     not null
);

create unique index task_id_uindex
    on tasks (id);

create table cheques
(
    id         serial                                                     not null
        constraint cheques_pk
            primary key,
    name       varchar(255)             default ''::character varying     not null,
    cost       integer                  default 0                         not null,
    type       varchar(255)             default 'work'::character varying not null,
    card_id    integer                                                    not null,
    user_id    integer,
    project_id integer,
    length     bigint                   default 0                         not null,
    created_at timestamp with time zone default CURRENT_TIMESTAMP         not null,
    file       text[]                   default ARRAY []::text[]          not null,
    weight     text[]                   default ARRAY []::text[]          not null
);

create unique index cheques_id_uindex
    on cheques (id);


create table files
(
    id         serial                                                 not null
        constraint table_name_pk
            primary key,
    created_at timestamp with time zone default CURRENT_TIMESTAMP,
    url        text                     default ''::text              not null,
    length     integer                                                not null,
    mime       text                     default ''::character varying not null,
    bucket     text                     default ''::character varying not null,
    object     text                     default ''::character varying not null,
    role       varchar(255)             default ''::character varying not null,
    tag        varchar(255)             default ''::character varying not null,
    user_id    integer                  default 0                     not null
);

create table managers
(
    id          serial                                            not null
        constraint managers_pk
            primary key,
    image       varchar(255) default ''::character varying        not null,
    email       varchar(255) default ''::character varying        not null,
    password    varchar(255) default ''::character varying        not null,
    last_name   varchar(255) default ''::character varying        not null,
    first_name  varchar(255) default ''::character varying        not null,
    middle_name varchar(255) default ''::character varying        not null,
    phone       varchar(255) default ''::character varying        not null,
    country     varchar(255) default ''::character varying        not null,
    city        varchar(255) default ''::character varying        not null,
    role        varchar(255) default 'MANAGER'::character varying not null,
    created_at  timestamp    default now(),
    updated_at  timestamp    default now(),
    token       varchar(255) default ''::character varying        not null
);

create unique index managers_email_uindex
    on managers (email);

create table manager_list
(
    id         serial not null
        constraint manager_list_pk
            primary key,
    manager_id integer,
    user_id    integer
);

create unique index manager_list_id_uindex
    on manager_list (id);

create unique index manager_list_pkey
    on manager_list (id);

create table addresses
(
    id          serial                                     not null
        constraint addresses_pk
            primary key,
    title       varchar(255) default ''::character varying not null,
    address     varchar(255) default ''::character varying not null,
    city        varchar(255) default ''::character varying not null,
    country     varchar(255) default ''::character varying not null,
    entrance    varchar(255) default ''::character varying not null,
    description text         default ''::character varying not null,
    user_id     integer                                    not null,
    project_id  integer      default 1                     not null
);

create unique index addresses_id_uindex
    on addresses (id);

create table email_problems
(
    id         serial                             not null
        constraint email_problems_pk
            primary key,
    email_name text default ''::character varying not null,
    email_data text default ''::character varying not null
);

create unique index email_problems_id_uindex
    on email_problems (id);

create table log_sms
(
    id      serial                             not null
        constraint log_sms_pk
            primary key,
    resp    text default ''::character varying not null,
    user_id integer
);

create unique index log_sms_id_uindex
    on log_sms (id);

















