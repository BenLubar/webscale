package schema // import "github.com/BenLubar/webscale/db/internal/schema"

type change struct {
	description string
	query       string
}

// This array holds the schema change scripts.
//
// Elements of this array must not be altered or removed once they are added.
var all = [...]change{
	{
		description: "initial schema - helpers",
		query: `
create extension citext;
create extension pgcrypto;

-- http://stackoverflow.com/a/20102665/2664560
create operator class _citext_ops default 
for type _citext using gin as 
operator 1 &&(anyarray, anyarray), 
operator 2 @>(anyarray, anyarray), 
operator 3 <@(anyarray, anyarray), 
operator 4 =(anyarray, anyarray), 
function 1 citext_cmp(citext, citext),
function 2 ginarrayextract(anyarray, internal, internal), 
function 3 ginqueryarrayextract(anyarray, internal, smallint, internal, internal, internal, internal), 
function 4 ginarrayconsistent(internal, smallint, anyarray, integer, internal, internal, internal, internal), 
storage citext;

create language plperlu;

create function unicode_nfkd(text) returns text as $$
	use Unicode::Normalize;
	return NFKD($_[0]);
$$ language plperlu immutable;

create function slugify(t text) returns text as $$
	select regexp_replace(trim(lower(regexp_replace(
		translate(unicode_nfkd(t), '-_', '  '),
		'[^\w\s]+', '', 'gi'
	))), '\s+', '-', 'gi')
$$ language sql immutable;

create function make_slug() returns trigger as $$
	begin
		new.slug = slugify(new.name);
		return new;
	end;
$$ language plpgsql;`,
	},
	{
		description: "initial schema - users",
		query: `
create table users (
	id bigserial primary key,
	name varchar(255) not null,
	slug varchar(255) not null constraint users_slug_exists check (slug <> ''),
	password varchar(60),
	email citext,
	join_date timestamp with time zone not null default now(),
	last_seen timestamp with time zone default null,
	address inet[] not null default '{}',
	birthday date default null,
	signature text not null default '',
	bio text not null default '',
	location text not null default '',
	website varchar(2048) not null default '',
	avatar varchar(2048) default null
);

create unique index users_slug on users (slug);
create index users_join_date on users (join_date);
create unique index users_email on users (email);
create index users_address on users (address);

create trigger users_make_slug
before insert or update of name, slug on users
for each row
execute procedure make_slug();`,
	},
	{
		description: "initial schema - sessions",
		query: `
create table sessions (
	id uuid not null primary key default gen_random_uuid(),
	user_id bigint not null references users(id) on delete cascade on update cascade,
	address inet[] not null,
	browser text not null,
	logged_in timestamp with time zone not null default now(),
	last_seen timestamp with time zone not null default now()
);

create index sessions_user on sessions (user_id);
`,
	},
	{
		description: "initial schema - groups",
		query: `
create table groups (
	id bigserial primary key,
	name varchar(255) not null,
	slug varchar(255) not null constraint groups_slug_exists check (slug <> '')
);

create unique index groups_slug on groups (slug);

create trigger groups_make_slug
before insert or update of name, slug on groups
for each row
execute procedure make_slug();

insert into groups (id, name) values
(-1, 'Registered Users'),
(-2, 'Guests'),
(-3, 'Administrators');

create table groups_users (
	group_id bigint references groups(id) on delete cascade on update cascade,
	user_id bigint references users(id) on delete cascade on update cascade,
	primary key (group_id, user_id),
	constraint groups_users_not_guest check (group_id <> -2)
);

create function add_user_to_registered_users() returns trigger as $$
	begin
		insert into groups_users (group_id, user_id) values (-1, new.id);
		return null;
	end;
$$ language plpgsql;

create trigger users_add_user_to_registered_users
after insert on users
for each row
execute procedure add_user_to_registered_users();

create table groups_groups (
	parent_group_id bigint references groups(id) on delete cascade on update cascade,
	group_id bigint references groups(id) on delete cascade on update cascade,
	primary key (parent_group_id, group_id),
	constraint groups_groups_not_special check (parent_group_id > 0 and group_id > 0)
);

create materialized view groups_groups_flat as
	with recursive ggf as (
		select g.id as parent_group_id, g.id as group_id
		from groups as g
		union
		select gg.parent_group_id, ggf.group_id
		from ggf
		inner join groups_groups as gg
		on gg.group_id = ggf.parent_group_id
	) select ggf.parent_group_id, ggf.group_id from ggf;

create function update_groups_groups_flat() returns trigger as $$
	begin
		refresh materialized view groups_groups_flat;
		return null;
	end;
$$ language plpgsql;

create trigger groups_update_groups_groups_flat
after insert or update of id or delete on groups
for each row
execute procedure update_groups_groups_flat();

create trigger groups_truncate_groups_groups_flat
after truncate on groups
for each statement
execute procedure update_groups_groups_flat();

create trigger groups_groups_update_groups_groups_flat
after insert or update or delete on groups_groups
for each row
execute procedure update_groups_groups_flat();

create trigger groups_groups_truncate_groups_groups_flat
after truncate on groups_groups
for each statement
execute procedure update_groups_groups_flat();

create index groups_groups_flat_parent on groups_groups_flat (parent_group_id, group_id);
create index groups_groups_flat_child on groups_groups_flat (group_id, parent_group_id);

create function groups_has_ancestor(ancestor_id bigint, child_id bigint) returns boolean as $$
	select exists(select 1 from groups_groups_flat where parent_group_id = ancestor_id and group_id = child_id)
$$ language sql stable;

alter table groups_groups add constraint groups_groups_no_cycles check(not groups_has_ancestor(group_id, parent_group_id));`,
	},
	{
		description: "initial schema - categories",
		query: `
create table categories (
	id bigserial primary key,
	name varchar(255) not null,
	slug varchar(255) not null constraint categories_slug_exists check (slug <> ''),
	parent_category_id bigint references categories(id) on delete cascade on update cascade default null
);

create unique index categories_slug on categories (slug);

create trigger categories_make_slug
before insert or update of name, slug on categories
for each row
execute procedure make_slug();

create recursive view categories_path (category_id, path, depth) as
	select c.id, array[c.id], 0
	from categories as c
	where c.parent_category_id is null
union
	select c.id, array_append(cp.path, c.id), cp.depth + 1
	from categories as c
	inner join categories_path as cp
	on c.parent_category_id = cp.category_id;

create function categories_has_ancestor(ancestor_id bigint, child_id bigint) returns boolean as $$
	select cp.path @> array[ancestor_id] from categories_path as cp where cp.category_id = child_id
$$ language sql stable;

alter table categories add constraint categories_no_cycles check(not categories_has_ancestor(id, parent_category_id));
`,
	},
	{
		description: "initial schema - topics",
		query: `
create table topics (
	id bigserial primary key,
	name varchar(255) not null,
	slug varchar(255) not null constraint topics_slug_exists check (slug <> ''),
	user_id bigint references users(id) on delete restrict on update cascade,
	category_id bigint not null constraint topics_category_fk references categories(id) on delete restrict on update cascade,
	created_at timestamp with time zone not null default now(),
	bumped_at timestamp with time zone not null default now()
);

create index topics_slug on topics (slug);
create index topics_user on topics (user_id);
create index topics_category on topics (category_id);
create index topics_created on topics (created_at);
create index topics_bumped on topics (bumped_at);

create trigger topics_make_slug
before insert or update of name, slug on topics
for each row
execute procedure make_slug();
`,
	},
	{
		description: "initial schema - posts",
		query: `
create table posts (
	id bigserial primary key,
	topic_id bigint not null references topics(id) on delete cascade on update cascade,
	user_id bigint references users(id) on delete restrict on update cascade,
	parent_post_id bigint references posts(id) on delete set null on update cascade,
	created_at timestamp with time zone not null default now(),
	content text not null,
	tags citext[] not null default '{}'
);

create index posts_user on posts (user_id);
create index posts_topic on posts (topic_id, created_at asc);
create index posts_tags on posts using gin (tags);

create table post_revisions (
	id bigserial primary key,
	post_id bigint references posts(id) on delete cascade on update cascade,
	user_id bigint references users(id) on delete set null on update cascade,
	created_at timestamp with time zone not null default now(),
	content text not null,
	tags citext[] not null default '{}'
);

create index post_revisions_post on post_revisions (post_id);
`,
	},
	{
		description: "initial schema - permissions",
		query: `
create table permissions (
	id bigserial primary key,
	slug varchar(255) not null constraint permissions_slug_slug check (slugify(slug) = slug)
);

create unique index permissions_slug on permissions (slug);

create table permission_sets (
	id bigserial primary key
);

create table permission_sets_users (
	set_id bigint not null references permission_sets(id) on delete cascade on update cascade,
	user_id bigint not null references users(id) on delete cascade on update cascade,
	primary key (set_id, user_id)
);

create table permission_sets_groups (
	set_id bigint not null references permission_sets(id) on delete cascade on update cascade,
	group_id bigint not null references groups(id) on delete cascade on update cascade,
	primary key (set_id, group_id)
);

create table permission_sets_categories (
	set_id bigint not null references permission_sets(id) on delete cascade on update cascade,
	category_id bigint not null references categories(id) on delete cascade on update cascade,
	primary key (set_id, category_id)
);

create table permission_sets_topics (
	set_id bigint not null references permission_sets(id) on delete cascade on update cascade,
	topic_id bigint not null references topics(id) on delete cascade on update cascade,
	primary key (set_id, topic_id)
);

create table permission_sets_posts (
	set_id bigint not null references permission_sets(id) on delete cascade on update cascade,
	post_id bigint not null references posts(id) on delete cascade on update cascade,
	primary key (set_id, post_id)
);

create table group_permissions (
	id bigserial primary key,
	group_id bigint not null references groups(id) on delete cascade on update cascade,
	permission_id bigint not null references permissions(id) on delete cascade on update cascade,
	priority bigint not null,
	set_id bigint references permission_sets(id) on delete restrict on update cascade,

	allow boolean not null,
	sudo boolean not null default false,
	self boolean not null default false,

	constraint guest_permissions_nospecial check (group_id <> -2 or (not sudo and not self))
);

create unique index group_permissions_unique on group_permissions (group_id, permission_id, coalesce(set_id, 0), sudo, self);

create view guest_permissions_flat as
	select gp.id, gp.permission_id, gp.priority, gp.set_id, gp.allow
	from group_permissions as gp
	where gp.group_id = -2;

create view user_permissions_flat as
	select gp.id, gu.user_id, gp.permission_id, gp.priority, gp.set_id, gp.allow, gp.sudo, gp.self
	from group_permissions as gp
	inner join groups_groups_flat as ggf
	on gp.group_id = ggf.parent_group_id
	inner join groups_users as gu
	on gu.group_id = ggf.group_id;

create view permission_sets_users_flat as
	select psu.set_id, psu.user_id
	from permission_sets_users as psu
union all
	select psg.set_id, gu.user_id
	from permission_sets_groups as psg
	inner join groups_groups_flat as ggf
	on psg.group_id = ggf.parent_group_id
	inner join groups_users as gu
	on gu.group_id = ggf.group_id;

create view permission_sets_groups_flat as
	select psg.set_id, ggf.group_id
	from permission_sets_groups as psg
	inner join groups_groups_flat as ggf
	on ggf.parent_group_id = psg.group_id;

create view permission_sets_categories_flat as
	select psc.set_id, cp.category_id
	from permission_sets_categories as psc
	inner join categories_path as cp
	on cp.path @> array[psc.category_id];

create view permission_sets_topics_flat as
	select pst.set_id, pst.topic_id
	from permission_sets_topics as pst
union all
	select pscf.set_id, t.id as topic_id
	from permission_sets_categories_flat as pscf
	inner join topics as t
	on pscf.category_id = t.category_id
union all
	select psu.set_id, t.id as topic_id
	from permission_sets_users as psu
	inner join topics as t
	on psu.user_id = t.user_id;

create view permission_sets_posts_flat as
	select psp.set_id, psp.post_id
	from permission_sets_posts as psp
union all
	select pst.set_id, p.id as post_id
	from permission_sets_topics as pst
	inner join posts as p
	on p.topic_id = pst.topic_id
union all
	select pscf.set_id, p.id as post_id
	from permission_sets_categories_flat as pscf
	inner join topics as t
	on pscf.category_id = t.category_id
	inner join posts as p
	on t.id = p.topic_id
union all
	select psu.set_id, p.id as post_id
	from permission_sets_users as psu
	inner join posts as p
	on psu.user_id = p.user_id;

create view new_user_permissions_flat as
	select gp.permission_id, gp.priority, gp.allow
	from group_permissions as gp
	where gp.group_id in (-1, -2)
	and not gp.sudo
	and (gp.set_id is null or gp.set_id in (select psgf.set_id from permission_sets_groups_flat as psgf where psgf.group_id = -1));

create function can(acting_user_id bigint, permission varchar(255), override boolean) returns boolean as $$
	select case when acting_user_id is null then (
		select gpf.allow
		from guest_permissions_flat as gpf
		inner join permissions as p
		on p.id = gpf.permission_id
		where p.slug = permission
		and gpf.set_id is null
		order by gpf.priority desc, gpf.allow asc
		limit 1
	) else (
		select upf.allow
		from user_permissions_flat as upf
		inner join permissions as p
		on p.id = upf.permission_id
		where p.slug = permission
		and upf.user_id = acting_user_id
		and upf.set_id is null
		and (not upf.sudo or override)
		and not upf.self
		order by upf.priority desc, upf.allow asc
		limit 1
	) end or (override and acting_user_id in (select gu.user_id from groups_users as gu where gu.group_id = -3))
$$ language sql stable;

create function can_user(acting_user_id bigint, permission varchar(255), override boolean, target_user_id bigint) returns boolean as $$
	with sets as (
		select psuf.set_id
		from permission_sets_users_flat as psuf
		where psuf.user_id = target_user_id
	) select (case when permission = 'user-meta' then true else can_user(acting_user_id, 'user-meta', override, target_user_id) end and case when acting_user_id is null then (
		select gpf.allow
		from guest_permissions_flat as gpf
		inner join permissions as p
		on p.id = gpf.permission_id
		where p.slug = permission
		and (gpf.set_id is null or gpf.set_id in (select set_id from sets))
		order by gpf.priority desc, gpf.allow asc
		limit 1
	) else (
		select upf.allow
		from user_permissions_flat as upf
		inner join permissions as p
		on p.id = upf.permission_id
		where p.slug = permission
		and upf.user_id = acting_user_id
		and (upf.set_id is null or upf.set_id in (select set_id from sets))
		and (not upf.sudo or override)
		and (not upf.self or acting_user_id = target_user_id)
		order by upf.priority desc, upf.allow asc
		limit 1
	) end) or (override and acting_user_id in (select gu.user_id from groups_users as gu where gu.group_id = -3))
$$ language sql stable;

create function can_group(acting_user_id bigint, permission varchar(255), override boolean, target_group_id bigint) returns boolean as $$
	with sets as (
		select psgf.set_id
		from permission_sets_groups_flat as psgf
		where psgf.group_id = target_group_id
	) select (case when permission = 'group-meta' then true else can_group(acting_user_id, 'group-meta', override, target_group_id) end and case when acting_user_id is null then (
		select gpf.allow
		from guest_permissions_flat as gpf
		inner join permissions as p
		on p.id = gpf.permission_id
		where p.slug = permission
		and (gpf.set_id is null or gpf.set_id in (select set_id from sets))
		order by gpf.priority desc, gpf.allow asc
		limit 1
	) else (
		select upf.allow
		from user_permissions_flat as upf
		inner join permissions as p
		on p.id = upf.permission_id
		where p.slug = permission
		and upf.user_id = acting_user_id
		and (upf.set_id is null or upf.set_id in (select set_id from sets))
		and (not upf.sudo or override)
		and not upf.self
		order by upf.priority desc, upf.allow asc
		limit 1
	) end) or (override and acting_user_id in (select gu.user_id from groups_users as gu where gu.group_id = -3))
$$ language sql stable;

create function can_category(acting_user_id bigint, permission varchar(255), override boolean, target_category_id bigint) returns boolean as $$
	with sets as (
		select pscf.set_id
		from permission_sets_categories_flat as pscf
		where pscf.category_id = target_category_id
	) select (case when permission = 'category-meta' then true else can_category(acting_user_id, 'category-meta', override, target_category_id) end and case when acting_user_id is null then (
		select gpf.allow
		from guest_permissions_flat as gpf
		inner join permissions as p
		on p.id = gpf.permission_id
		where p.slug = permission
		and (gpf.set_id is null or gpf.set_id in (select set_id from sets))
		order by gpf.priority desc, gpf.allow asc
		limit 1
	) else (
		select upf.allow
		from user_permissions_flat as upf
		inner join permissions as p
		on p.id = upf.permission_id
		where p.slug = permission
		and upf.user_id = acting_user_id
		and (upf.set_id is null or upf.set_id in (select set_id from sets))
		and (not upf.sudo or override)
		and not upf.self
		order by upf.priority desc, upf.allow asc
		limit 1
	) end) or (override and acting_user_id in (select gu.user_id from groups_users as gu where gu.group_id = -3))
$$ language sql stable;

create function can_topic(acting_user_id bigint, permission varchar(255), override boolean, target_topic_id bigint) returns boolean as $$
	with sets as (
		select pstf.set_id
		from permission_sets_topics_flat as pstf
		where pstf.topic_id = target_topic_id
	) select (case when permission = 'topic-meta' then true else can_topic(acting_user_id, 'topic-meta', override, target_topic_id) end and case when acting_user_id is null then (
		select gpf.allow
		from guest_permissions_flat as gpf
		inner join permissions as p
		on p.id = gpf.permission_id
		where p.slug = permission
		and (gpf.set_id is null or gpf.set_id in (select set_id from sets))
		order by gpf.priority desc, gpf.allow asc
		limit 1
	) else (
		select upf.allow
		from user_permissions_flat as upf
		inner join permissions as p
		on p.id = upf.permission_id
		where p.slug = permission
		and upf.user_id = acting_user_id
		and (upf.set_id is null or upf.set_id in (select set_id from sets))
		and (not upf.sudo or override)
		and (not upf.self or (select t.user_id from topics as t where t.id = target_topic_id) = acting_user_id)
		order by upf.priority desc, upf.allow asc
		limit 1
	) end) or (override and acting_user_id in (select gu.user_id from groups_users as gu where gu.group_id = -3))
$$ language sql stable;

create function can_post(acting_user_id bigint, permission varchar(255), override boolean, target_post_id bigint) returns boolean as $$
	with sets as (
		select pspf.set_id
		from permission_sets_posts_flat as pspf
		where pspf.post_id = target_post_id
	) select (case when permission = 'post-meta' then true else can_post(acting_user_id, 'post-meta', override, target_post_id) end and case when acting_user_id is null then (
		select gpf.allow
		from guest_permissions_flat as gpf
		inner join permissions as p
		on p.id = gpf.permission_id
		where p.slug = permission
		and (gpf.set_id is null or gpf.set_id in (select set_id from sets))
		order by gpf.priority desc, gpf.allow asc
		limit 1
	) else (
		select upf.allow
		from user_permissions_flat as upf
		inner join permissions as p
		on p.id = upf.permission_id
		where p.slug = permission
		and upf.user_id = acting_user_id
		and (upf.set_id is null or upf.set_id in (select set_id from sets))
		and (not upf.sudo or override)
		and (not upf.self or (select p.user_id from posts as p where p.id = target_post_id) = acting_user_id)
		order by upf.priority desc, upf.allow asc
		limit 1
	) end) or (override and acting_user_id in (select gu.user_id from groups_users as gu where gu.group_id = -3))
$$ language sql stable;

insert into permissions (slug) values
('global-log-in'),
('global-create-user'),
('global-create-group'),
('global-create-category'),
('global-edit-permissions'),
('user-view-profile'),
('user-view-sessions'),
('user-delete-session'),
('user-meta'),
('user-view-email'),
('user-view-join-date'),
('user-view-last-seen'),
('user-view-ip-address'),
('user-view-birthday'),
('user-view-signature'),
('user-view-bio'),
('user-view-location'),
('user-view-website'),
('user-view-avatar'),
('user-edit-name'),
('user-edit-email'),
('user-edit-birthday'),
('user-edit-signature'),
('user-edit-bio'),
('user-edit-location'),
('user-edit-website'),
('user-edit-avatar'),
('user-ban'),
('group-meta'),
('group-view-members'),
('group-edit-name'),
('group-delete'),
('group-add-member'),
('group-remove-member'),
('group-add-self'),
('group-remove-self'),
('category-meta'),
('category-delete'),
('category-list-topics'),
('category-create-topic'),
('topic-meta'),
('topic-view-author'),
('topic-create-reply'),
('topic-edit-category'),
('topic-edit-name'),
('topic-edit-author'),
('post-meta'),
('post-view-author'),
('post-view-history'),
('post-view-tags'),
('post-edit-author'),
('post-edit-content'),
('post-delete'),
('post-add-tag'),
('post-remove-tag'),
('post-delete-revision');`,
	},
	{
		description: "initial schema - defaults",
		query: `
do $$
declare
	admins bigint := -3;
	guests bigint := -2;
	users bigint := -1;
	mods bigint;
	mods_set bigint;
	staff bigint;
	staff_set bigint;

begin

insert into groups (name) values ('Global Moderators') returning id into strict mods;
insert into permission_sets default values returning id into strict mods_set;
insert into permission_sets_groups (set_id, group_id) values (mods_set, mods);

insert into categories (name) values ('Staff') returning id into strict staff;
insert into permission_sets default values returning id into strict staff_set;
insert into permission_sets_categories (set_id, category_id) values (staff_set, staff);

insert into categories (name) values ('General');

insert into group_permissions (group_id, permission_id, priority, allow, self, sudo, set_id) values
(users, (select id from permissions where slug = 'global-log-in'), 0, true, false, false, null),
(guests, (select id from permissions where slug = 'global-create-user'), 0, true, false, false, null),
(users, (select id from permissions where slug = 'user-view-profile'), 0, true, false, false, null),
(guests, (select id from permissions where slug = 'user-view-profile'), 0, true, false, false, null),
(users, (select id from permissions where slug = 'user-view-sessions'), 0, true, true, false, null),
(users, (select id from permissions where slug = 'user-delete-session'), 0, true, true, false, null),
(users, (select id from permissions where slug = 'user-meta'), 0, true, false, false, null),
(guests, (select id from permissions where slug = 'user-meta'), 0, true, false, false, null),
(users, (select id from permissions where slug = 'user-view-email'), 0, true, true, false, null),
(users, (select id from permissions where slug = 'user-view-join-date'), 0, true, false, false, null),
(guests, (select id from permissions where slug = 'user-view-join-date'), 0, true, false, false, null),
(users, (select id from permissions where slug = 'user-view-last-seen'), 0, true, false, false, null),
(guests, (select id from permissions where slug = 'user-view-last-seen'), 0, true, false, false, null),
(users, (select id from permissions where slug = 'user-view-birthday'), 0, true, false, false, null),
(guests, (select id from permissions where slug = 'user-view-birthday'), 0, true, false, false, null),
(users, (select id from permissions where slug = 'user-view-signature'), 0, true, false, false, null),
(guests, (select id from permissions where slug = 'user-view-signature'), 0, true, false, false, null),
(users, (select id from permissions where slug = 'user-view-bio'), 0, true, false, false, null),
(guests, (select id from permissions where slug = 'user-view-bio'), 0, true, false, false, null),
(users, (select id from permissions where slug = 'user-view-location'), 0, true, false, false, null),
(guests, (select id from permissions where slug = 'user-view-location'), 0, true, false, false, null),
(users, (select id from permissions where slug = 'user-view-website'), 0, true, false, false, null),
(guests, (select id from permissions where slug = 'user-view-website'), 0, true, false, false, null),
(users, (select id from permissions where slug = 'user-view-avatar'), 0, true, false, false, null),
(guests, (select id from permissions where slug = 'user-view-avatar'), 0, true, false, false, null),
(users, (select id from permissions where slug = 'user-edit-email'), 0, true, true, false, null),
(users, (select id from permissions where slug = 'user-edit-birthday'), 0, true, true, false, null),
(users, (select id from permissions where slug = 'user-edit-signature'), 0, true, true, false, null),
(users, (select id from permissions where slug = 'user-edit-bio'), 0, true, true, false, null),
(users, (select id from permissions where slug = 'user-edit-location'), 0, true, true, false, null),
(users, (select id from permissions where slug = 'user-edit-website'), 0, true, true, false, null),
(users, (select id from permissions where slug = 'user-edit-avatar'), 0, true, true, false, null),
(users, (select id from permissions where slug = 'group-meta'), 0, true, false, false, null),
(guests, (select id from permissions where slug = 'group-meta'), 0, true, false, false, null),
(users, (select id from permissions where slug = 'group-view-members'), 0, true, false, false, null),
(guests, (select id from permissions where slug = 'group-view-members'), 0, true, false, false, null),
(mods, (select id from permissions where slug = 'group-remove-self'), 0, true, false, true, mods_set),
(users, (select id from permissions where slug = 'category-meta'), 0, true, false, false, null),
(guests, (select id from permissions where slug = 'category-meta'), 0, true, false, false, null),
(users, (select id from permissions where slug = 'category-meta'), 100, false, false, false, staff_set),
(guests, (select id from permissions where slug = 'category-meta'), 100, false, false, false, staff_set),
(mods, (select id from permissions where slug = 'category-meta'), 200, true, false, false, staff_set),
(admins, (select id from permissions where slug = 'category-meta'), 200, true, false, false, staff_set),
(users, (select id from permissions where slug = 'category-list-topics'), 0, true, false, false, null),
(guests, (select id from permissions where slug = 'category-list-topics'), 0, true, false, false, null),
(users, (select id from permissions where slug = 'category-create-topic'), 0, true, false, false, null),
(users, (select id from permissions where slug = 'topic-meta'), 0, true, false, false, null),
(guests, (select id from permissions where slug = 'topic-meta'), 0, true, false, false, null),
(users, (select id from permissions where slug = 'topic-meta'), 100, false, false, false, staff_set),
(guests, (select id from permissions where slug = 'topic-meta'), 100, false, false, false, staff_set),
(mods, (select id from permissions where slug = 'topic-meta'), 200, true, false, false, staff_set),
(admins, (select id from permissions where slug = 'topic-meta'), 200, true, false, false, staff_set),
(users, (select id from permissions where slug = 'topic-view-author'), 0, true, false, false, null),
(guests, (select id from permissions where slug = 'topic-view-author'), 0, true, false, false, null),
(users, (select id from permissions where slug = 'topic-create-reply'), 0, true, false, false, null),
(mods, (select id from permissions where slug = 'topic-edit-category'), 0, true, false, true, null),
(users, (select id from permissions where slug = 'topic-edit-name'), 0, true, true, false, null),
(mods, (select id from permissions where slug = 'topic-edit-name'), 0, true, false, true, null),
(users, (select id from permissions where slug = 'post-meta'), 0, true, false, false, null),
(guests, (select id from permissions where slug = 'post-meta'), 0, true, false, false, null),
(users, (select id from permissions where slug = 'post-meta'), 100, false, false, false, staff_set),
(guests, (select id from permissions where slug = 'post-meta'), 100, false, false, false, staff_set),
(mods, (select id from permissions where slug = 'post-meta'), 200, true, false, false, staff_set),
(admins, (select id from permissions where slug = 'post-meta'), 200, true, false, false, staff_set),
(users, (select id from permissions where slug = 'post-view-author'), 0, true, false, false, null),
(guests, (select id from permissions where slug = 'post-view-author'), 0, true, false, false, null),
(users, (select id from permissions where slug = 'post-view-history'), 0, true, true, false, null),
(mods, (select id from permissions where slug = 'post-view-history'), 0, true, false, true, null),
(users, (select id from permissions where slug = 'post-view-tags'), 0, true, false, false, null),
(guests, (select id from permissions where slug = 'post-view-tags'), 0, true, false, false, null),
(users, (select id from permissions where slug = 'post-edit-content'), 0, true, true, false, null),
(mods, (select id from permissions where slug = 'post-edit-content'), 0, true, false, true, null),
(mods, (select id from permissions where slug = 'post-delete'), 0, true, false, true, null),
(users, (select id from permissions where slug = 'post-add-tag'), 0, true, true, false, null),
(mods, (select id from permissions where slug = 'post-add-tag'), 0, true, false, true, null),
(users, (select id from permissions where slug = 'post-remove-tag'), 0, true, true, false, null),
(mods, (select id from permissions where slug = 'post-remove-tag'), 0, true, false, true, null);

end;
$$`,
	},
	{
		description: "",
		query: `
create index topics_category_bumped on topics (category_id, bumped_at desc);`,
	},
}
