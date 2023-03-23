BEGIN;

-- 
-- Member Statuses
--
CREATE TABLE IF NOT EXISTS public.member_status (
	id serial,
	status character varying NOT NULL,
	PRIMARY KEY(id)
);

INSERT INTO public.job_status (id, status) VALUES
(1, 'Pending Approval'),
(2, 'Active'),
(3, 'Inactive')
;

--
-- Member Roles
--
CREATE TABLE IF NOT EXISTS public.member_role (
  id serial,
  created_at timestamp without time zone NOT NULL,
	updated_at timestamp without time zone,
	deleted_at timestamp without time zone,
  color character varying NOT NULL,
  role character varying NOT NULL,
  PRIMARY KEY(id)
);

--
-- Members
--
CREATE TABLE IF NOT EXISTS public.member (
	id character varying NOT NULL,
	created_at timestamp without time zone NOT NULL,
	updated_at timestamp without time zone,
	deleted_at timestamp without time zone,
  avatar_url character varying NOT NULL,
  email character varying NOT NULL,
  external_id character varying NOT NULL,
  first_name character varying NOT NULL,
  last_name character varying NOT NULL,
  password character varying NOT NULL,
  role_id bigint references public.member_role(id),
  status_id bigint references public.member_status(id),
  PRIMARY KEY(id)
);

CREATE INDEX idx_member_email ON public.member (email);

COMMIT;

