CREATE TABLE public.classes (
	class_id int4 NOT NULL,
	class_name varchar(50) NOT NULL,
	days varchar(20) NOT NULL,
	start_time time NOT NULL,
	end_time time NOT NULL,
	price int4 NULL,
	teacher_id varchar(30) NULL,
	CONSTRAINT class_check CHECK ((end_time > start_time)),
	CONSTRAINT class_pkey PRIMARY KEY (class_id),
	CONSTRAINT class_price_check CHECK ((price >= 0)),
	CONSTRAINT fk_lecture_teacher FOREIGN KEY (teacher_id) REFERENCES public.teachers(teacher_id)
);