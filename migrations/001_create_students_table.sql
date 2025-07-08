CREATE TABLE public.students (
	student_id int4 NOT NULL,
	"name" varchar(20) NOT NULL,
	grade varchar(10) NOT NULL,
	phone varchar(15) NULL,
	parent_phone varchar(15) NOT NULL,
	CONSTRAINT students_pkey PRIMARY KEY (student_id),
	CONSTRAINT students_student_id_check CHECK (((student_id >= 1000) AND (student_id <= 99999)))
);
