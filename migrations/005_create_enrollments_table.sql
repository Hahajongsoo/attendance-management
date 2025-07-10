CREATE TABLE public.enrollments (
	student_id int4 NOT NULL,
	class_id int4 NOT NULL,
	enrolled_date date DEFAULT CURRENT_DATE NOT NULL,
	enrollment_id serial4 NOT NULL,
	CONSTRAINT enrollment_unique UNIQUE (student_id, class_id),
	CONSTRAINT enrollments_pkey PRIMARY KEY (enrollment_id),
	CONSTRAINT enrollment_lecture_id_fkey FOREIGN KEY (class_id) REFERENCES public.classes(class_id),
	CONSTRAINT enrollment_student_id_fkey FOREIGN KEY (student_id) REFERENCES public.students(student_id)
);