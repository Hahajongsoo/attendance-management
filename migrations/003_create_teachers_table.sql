CREATE TABLE public.teachers (
	teacher_id varchar(30) NOT NULL,
	"password" varchar(100) NOT NULL,
	name varchar(30) NOT NULL,
	phone_number varchar(20) NULL,
	CONSTRAINT teacher_pkey PRIMARY KEY (teacher_id)
);