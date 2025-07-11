CREATE TABLE public.payments (
    payment_id serial4 NOT NULL,
	student_id int4 NOT NULL,
	class_id int4 NOT NULL,
	payment_date date DEFAULT CURRENT_DATE NOT NULL,
	amount int4 NULL,
	enrollment_id int4 NULL,
	CONSTRAINT payment_amount_check CHECK ((amount >= 0)),
	CONSTRAINT payment_unique UNIQUE (student_id, class_id, payment_date),
	CONSTRAINT payments_pkey PRIMARY KEY (payment_id),
	CONSTRAINT payments_enrollment_id_fkey FOREIGN KEY (enrollment_id) REFERENCES public.enrollments(enrollment_id)
);