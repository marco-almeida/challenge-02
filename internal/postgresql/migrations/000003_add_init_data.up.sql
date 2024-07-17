-- init vehicles
INSERT INTO public.vehicle (max_weight_capacity,number_plate) VALUES
	 (40000,'11-ad-92'),
	 (40000,'74-qg-23'),
	 (36000,'10-ae-12'),
	 (32000,'55-bk-30');

-- init orders
INSERT INTO public.order (weight,destination,observations) VALUES
	 (35000.23,'(47.2342,5.45)','to be delivered in gate A'),
	 (1000,'(38.71,-9.14)','Z Lisbon stuff'),
     (10000, '(40.97134, -5.66337)', 'A'),
     (10000, '(39.47590, -0.37696)', 'B'),
     (10000, '(41.15936, -8.62907)', 'C'),
     (10000, '(40.42382, -3.70254)', 'D');

