CREATE TABLE public.attendance (
	student_id int4 NOT NULL,
	"date" date DEFAULT CURRENT_DATE NOT NULL,
	check_in time NULL,
	check_out time NULL,
	status text NOT NULL,
	CONSTRAINT attendance_pkey PRIMARY KEY (student_id, date),
	CONSTRAINT attendance_status_check CHECK ((status = ANY (ARRAY['출석'::text, '결석'::text, '지각'::text])))
);


-- public.attendance foreign keys

ALTER TABLE public.attendance ADD CONSTRAINT attendance_student_id_fkey FOREIGN KEY (student_id) REFERENCES public.students(student_id);