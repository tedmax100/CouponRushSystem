-- User table
CREATE TABLE public.user (
    id BIGINT  NOT NULL PRIMARY KEY,
    name VARCHAR(30) NOT NULL
);

-- Coupon activaty
CREATE TYPE coupon_active_state AS ENUM ('NOT_OPEN', 'OPENING', 'CLOSED');

CREATE TABLE coupon_active (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    date Date NOT NULL,
    begin_time TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    end_time TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    purchase_begin_time TIMESTAMP WITHOUT TIME ZONE,
    purchase_end_time TIMESTAMP WITHOUT TIME ZONE,
    state coupon_active_state NOT NULL
);

-- history about user reserve coupon
CREATE TABLE user_reserved_coupon_active_history (
    user_id BIGINT NOT NULL,
    active_id BIGINT NOT NULL,
    serial_num BIGINT NOT NULL,
    reserved_at BIGINT,
    PRIMARY KEY (user_id, active_id),
    FOREIGN KEY (user_id) REFERENCES public.user (id),
    FOREIGN KEY (active_id) REFERENCES public.coupon_active (id)
);

CREATE TYPE coupon_state AS ENUM ('UNRESERVED', 'RESERVED', 'USED');

-- Coupon table
CREATE TABLE coupon (
    id BIGINT PRIMARY KEY,
    active_id BIGINT,
    coupon_code VARCHAR(50) NOT NULL,
    state coupon_state NOT NULL,
    buyer BIGINT,
    buy_time BIGINT,
    created_at BIGINT NOT NULL,
    FOREIGN KEY (active_id) REFERENCES public.coupon_active (id),
    FOREIGN KEY (buyer) REFERENCES public.user (id)
);

INSERT INTO public."user" (id, name) VALUES (1, 'user1');
INSERT INTO public."user" (id, name) VALUES (2, 'user2');
INSERT INTO public."user" (id, name) VALUES (3, 'user3');
INSERT INTO public."user" (id, name) VALUES (4, 'user4');
INSERT INTO public."user" (id, name) VALUES (5, 'user5');
INSERT INTO public."user" (id, name) VALUES (6, 'user6');
INSERT INTO public."user" (id, name) VALUES (7, 'user7');
INSERT INTO public."user" (id, name) VALUES (8, 'user8');
INSERT INTO public."user" (id, name) VALUES (9, 'user9');
INSERT INTO public."user" (id, name) VALUES (10, 'user10');
INSERT INTO public."user" (id, name) VALUES (11, 'user11');
INSERT INTO public."user" (id, name) VALUES (12, 'user12');
INSERT INTO public."user" (id, name) VALUES (13, 'user13');
INSERT INTO public."user" (id, name) VALUES (14, 'user14');
INSERT INTO public."user" (id, name) VALUES (15, 'user15');
INSERT INTO public."user" (id, name) VALUES (16, 'user16');
INSERT INTO public."user" (id, name) VALUES (17, 'user17');
INSERT INTO public."user" (id, name) VALUES (18, 'user18');
INSERT INTO public."user" (id, name) VALUES (19, 'user19');
INSERT INTO public."user" (id, name) VALUES (20, 'user20');
INSERT INTO public."user" (id, name) VALUES (21, 'user21');
INSERT INTO public."user" (id, name) VALUES (22, 'user22');
INSERT INTO public."user" (id, name) VALUES (23, 'user23');
INSERT INTO public."user" (id, name) VALUES (24, 'user24');
INSERT INTO public."user" (id, name) VALUES (25, 'user25');
INSERT INTO public."user" (id, name) VALUES (26, 'user26');
INSERT INTO public."user" (id, name) VALUES (27, 'user27');
INSERT INTO public."user" (id, name) VALUES (28, 'user28');
INSERT INTO public."user" (id, name) VALUES (29, 'user29');
INSERT INTO public."user" (id, name) VALUES (30, 'user30');
INSERT INTO public."user" (id, name) VALUES (31, 'user31');
INSERT INTO public."user" (id, name) VALUES (32, 'user32');
INSERT INTO public."user" (id, name) VALUES (33, 'user33');
INSERT INTO public."user" (id, name) VALUES (34, 'user34');
INSERT INTO public."user" (id, name) VALUES (35, 'user35');
INSERT INTO public."user" (id, name) VALUES (36, 'user36');
INSERT INTO public."user" (id, name) VALUES (37, 'user37');
INSERT INTO public."user" (id, name) VALUES (38, 'user38');
INSERT INTO public."user" (id, name) VALUES (39, 'user39');
INSERT INTO public."user" (id, name) VALUES (40, 'user40');
INSERT INTO public."user" (id, name) VALUES (41, 'user41');
INSERT INTO public."user" (id, name) VALUES (42, 'user42');
INSERT INTO public."user" (id, name) VALUES (43, 'user43');
INSERT INTO public."user" (id, name) VALUES (44, 'user44');
INSERT INTO public."user" (id, name) VALUES (45, 'user45');
INSERT INTO public."user" (id, name) VALUES (46, 'user46');
INSERT INTO public."user" (id, name) VALUES (47, 'user47');
INSERT INTO public."user" (id, name) VALUES (48, 'user48');
INSERT INTO public."user" (id, name) VALUES (49, 'user49');
INSERT INTO public."user" (id, name) VALUES (50, 'user50');
INSERT INTO public."user" (id, name) VALUES (51, 'user51');
INSERT INTO public."user" (id, name) VALUES (52, 'user52');
INSERT INTO public."user" (id, name) VALUES (53, 'user53');
INSERT INTO public."user" (id, name) VALUES (54, 'user54');
INSERT INTO public."user" (id, name) VALUES (55, 'user55');
INSERT INTO public."user" (id, name) VALUES (56, 'user56');
INSERT INTO public."user" (id, name) VALUES (57, 'user57');
INSERT INTO public."user" (id, name) VALUES (58, 'user58');
INSERT INTO public."user" (id, name) VALUES (59, 'user59');
INSERT INTO public."user" (id, name) VALUES (60, 'user60');
INSERT INTO public."user" (id, name) VALUES (61, 'user61');
INSERT INTO public."user" (id, name) VALUES (62, 'user62');
INSERT INTO public."user" (id, name) VALUES (63, 'user63');
INSERT INTO public."user" (id, name) VALUES (64, 'user64');
INSERT INTO public."user" (id, name) VALUES (65, 'user65');
INSERT INTO public."user" (id, name) VALUES (66, 'user66');
INSERT INTO public."user" (id, name) VALUES (67, 'user67');
INSERT INTO public."user" (id, name) VALUES (68, 'user68');
INSERT INTO public."user" (id, name) VALUES (69, 'user69');
INSERT INTO public."user" (id, name) VALUES (70, 'user70');
INSERT INTO public."user" (id, name) VALUES (71, 'user71');
INSERT INTO public."user" (id, name) VALUES (72, 'user72');
INSERT INTO public."user" (id, name) VALUES (73, 'user73');
INSERT INTO public."user" (id, name) VALUES (74, 'user74');
INSERT INTO public."user" (id, name) VALUES (75, 'user75');
INSERT INTO public."user" (id, name) VALUES (76, 'user76');
INSERT INTO public."user" (id, name) VALUES (77, 'user77');
INSERT INTO public."user" (id, name) VALUES (78, 'user78');
INSERT INTO public."user" (id, name) VALUES (79, 'user79');
INSERT INTO public."user" (id, name) VALUES (80, 'user80');
INSERT INTO public."user" (id, name) VALUES (81, 'user81');
INSERT INTO public."user" (id, name) VALUES (82, 'user82');
INSERT INTO public."user" (id, name) VALUES (83, 'user83');
INSERT INTO public."user" (id, name) VALUES (84, 'user84');
INSERT INTO public."user" (id, name) VALUES (85, 'user85');
INSERT INTO public."user" (id, name) VALUES (86, 'user86');
INSERT INTO public."user" (id, name) VALUES (87, 'user87');
INSERT INTO public."user" (id, name) VALUES (88, 'user88');
INSERT INTO public."user" (id, name) VALUES (89, 'user89');
INSERT INTO public."user" (id, name) VALUES (90, 'user90');
INSERT INTO public."user" (id, name) VALUES (91, 'user91');
INSERT INTO public."user" (id, name) VALUES (92, 'user92');
INSERT INTO public."user" (id, name) VALUES (93, 'user93');
INSERT INTO public."user" (id, name) VALUES (94, 'user94');
INSERT INTO public."user" (id, name) VALUES (95, 'user95');
INSERT INTO public."user" (id, name) VALUES (96, 'user96');
INSERT INTO public."user" (id, name) VALUES (97, 'user97');
INSERT INTO public."user" (id, name) VALUES (98, 'user98');
INSERT INTO public."user" (id, name) VALUES (99, 'user99');
INSERT INTO public."user" (id, name) VALUES (100, 'user100');
INSERT INTO public."user" (id, name) VALUES (101, 'user101');
INSERT INTO public."user" (id, name) VALUES (102, 'user102');
INSERT INTO public."user" (id, name) VALUES (103, 'user103');
INSERT INTO public."user" (id, name) VALUES (104, 'user104');
INSERT INTO public."user" (id, name) VALUES (105, 'user105');
INSERT INTO public."user" (id, name) VALUES (106, 'user106');
INSERT INTO public."user" (id, name) VALUES (107, 'user107');
INSERT INTO public."user" (id, name) VALUES (108, 'user108');
INSERT INTO public."user" (id, name) VALUES (109, 'user109');
INSERT INTO public."user" (id, name) VALUES (110, 'user110');
INSERT INTO public."user" (id, name) VALUES (111, 'user111');
INSERT INTO public."user" (id, name) VALUES (112, 'user112');
INSERT INTO public."user" (id, name) VALUES (113, 'user113');
INSERT INTO public."user" (id, name) VALUES (114, 'user114');
INSERT INTO public."user" (id, name) VALUES (115, 'user115');
INSERT INTO public."user" (id, name) VALUES (116, 'user116');
INSERT INTO public."user" (id, name) VALUES (117, 'user117');
INSERT INTO public."user" (id, name) VALUES (118, 'user118');
INSERT INTO public."user" (id, name) VALUES (119, 'user119');
INSERT INTO public."user" (id, name) VALUES (120, 'user120');
INSERT INTO public."user" (id, name) VALUES (121, 'user121');
INSERT INTO public."user" (id, name) VALUES (122, 'user122');
INSERT INTO public."user" (id, name) VALUES (123, 'user123');
INSERT INTO public."user" (id, name) VALUES (124, 'user124');
INSERT INTO public."user" (id, name) VALUES (125, 'user125');
INSERT INTO public."user" (id, name) VALUES (126, 'user126');
INSERT INTO public."user" (id, name) VALUES (127, 'user127');
INSERT INTO public."user" (id, name) VALUES (128, 'user128');
INSERT INTO public."user" (id, name) VALUES (129, 'user129');
INSERT INTO public."user" (id, name) VALUES (130, 'user130');
INSERT INTO public."user" (id, name) VALUES (131, 'user131');
INSERT INTO public."user" (id, name) VALUES (132, 'user132');
INSERT INTO public."user" (id, name) VALUES (133, 'user133');
INSERT INTO public."user" (id, name) VALUES (134, 'user134');
INSERT INTO public."user" (id, name) VALUES (135, 'user135');
INSERT INTO public."user" (id, name) VALUES (136, 'user136');
INSERT INTO public."user" (id, name) VALUES (137, 'user137');
INSERT INTO public."user" (id, name) VALUES (138, 'user138');
INSERT INTO public."user" (id, name) VALUES (139, 'user139');
INSERT INTO public."user" (id, name) VALUES (140, 'user140');
INSERT INTO public."user" (id, name) VALUES (141, 'user141');
INSERT INTO public."user" (id, name) VALUES (142, 'user142');
INSERT INTO public."user" (id, name) VALUES (143, 'user143');
INSERT INTO public."user" (id, name) VALUES (144, 'user144');
INSERT INTO public."user" (id, name) VALUES (145, 'user145');
INSERT INTO public."user" (id, name) VALUES (146, 'user146');
INSERT INTO public."user" (id, name) VALUES (147, 'user147');
INSERT INTO public."user" (id, name) VALUES (148, 'user148');
INSERT INTO public."user" (id, name) VALUES (149, 'user149');
INSERT INTO public."user" (id, name) VALUES (150, 'user150');
INSERT INTO public."user" (id, name) VALUES (151, 'user151');
INSERT INTO public."user" (id, name) VALUES (152, 'user152');
INSERT INTO public."user" (id, name) VALUES (153, 'user153');
INSERT INTO public."user" (id, name) VALUES (154, 'user154');
INSERT INTO public."user" (id, name) VALUES (155, 'user155');
INSERT INTO public."user" (id, name) VALUES (156, 'user156');
INSERT INTO public."user" (id, name) VALUES (157, 'user157');
INSERT INTO public."user" (id, name) VALUES (158, 'user158');
INSERT INTO public."user" (id, name) VALUES (159, 'user159');
INSERT INTO public."user" (id, name) VALUES (160, 'user160');
INSERT INTO public."user" (id, name) VALUES (161, 'user161');
INSERT INTO public."user" (id, name) VALUES (162, 'user162');
INSERT INTO public."user" (id, name) VALUES (163, 'user163');
INSERT INTO public."user" (id, name) VALUES (164, 'user164');
INSERT INTO public."user" (id, name) VALUES (165, 'user165');
INSERT INTO public."user" (id, name) VALUES (166, 'user166');
INSERT INTO public."user" (id, name) VALUES (167, 'user167');
INSERT INTO public."user" (id, name) VALUES (168, 'user168');
INSERT INTO public."user" (id, name) VALUES (169, 'user169');
INSERT INTO public."user" (id, name) VALUES (170, 'user170');
INSERT INTO public."user" (id, name) VALUES (171, 'user171');
INSERT INTO public."user" (id, name) VALUES (172, 'user172');
INSERT INTO public."user" (id, name) VALUES (173, 'user173');
INSERT INTO public."user" (id, name) VALUES (174, 'user174');
INSERT INTO public."user" (id, name) VALUES (175, 'user175');
INSERT INTO public."user" (id, name) VALUES (176, 'user176');
INSERT INTO public."user" (id, name) VALUES (177, 'user177');
INSERT INTO public."user" (id, name) VALUES (178, 'user178');
INSERT INTO public."user" (id, name) VALUES (179, 'user179');
INSERT INTO public."user" (id, name) VALUES (180, 'user180');
INSERT INTO public."user" (id, name) VALUES (181, 'user181');
INSERT INTO public."user" (id, name) VALUES (182, 'user182');
INSERT INTO public."user" (id, name) VALUES (183, 'user183');
INSERT INTO public."user" (id, name) VALUES (184, 'user184');
INSERT INTO public."user" (id, name) VALUES (185, 'user185');
INSERT INTO public."user" (id, name) VALUES (186, 'user186');
INSERT INTO public."user" (id, name) VALUES (187, 'user187');
INSERT INTO public."user" (id, name) VALUES (188, 'user188');
INSERT INTO public."user" (id, name) VALUES (189, 'user189');
INSERT INTO public."user" (id, name) VALUES (190, 'user190');
INSERT INTO public."user" (id, name) VALUES (191, 'user191');
INSERT INTO public."user" (id, name) VALUES (192, 'user192');
INSERT INTO public."user" (id, name) VALUES (193, 'user193');
INSERT INTO public."user" (id, name) VALUES (194, 'user194');
INSERT INTO public."user" (id, name) VALUES (195, 'user195');
INSERT INTO public."user" (id, name) VALUES (196, 'user196');
INSERT INTO public."user" (id, name) VALUES (197, 'user197');
INSERT INTO public."user" (id, name) VALUES (198, 'user198');
INSERT INTO public."user" (id, name) VALUES (199, 'user199');
INSERT INTO public."user" (id, name) VALUES (200, 'user200');
INSERT INTO public."user" (id, name) VALUES (201, 'user201');
INSERT INTO public."user" (id, name) VALUES (202, 'user202');
INSERT INTO public."user" (id, name) VALUES (203, 'user203');
INSERT INTO public."user" (id, name) VALUES (204, 'user204');
INSERT INTO public."user" (id, name) VALUES (205, 'user205');
INSERT INTO public."user" (id, name) VALUES (206, 'user206');
INSERT INTO public."user" (id, name) VALUES (207, 'user207');
INSERT INTO public."user" (id, name) VALUES (208, 'user208');
INSERT INTO public."user" (id, name) VALUES (209, 'user209');
INSERT INTO public."user" (id, name) VALUES (210, 'user210');
INSERT INTO public."user" (id, name) VALUES (211, 'user211');
INSERT INTO public."user" (id, name) VALUES (212, 'user212');
INSERT INTO public."user" (id, name) VALUES (213, 'user213');
INSERT INTO public."user" (id, name) VALUES (214, 'user214');
INSERT INTO public."user" (id, name) VALUES (215, 'user215');
INSERT INTO public."user" (id, name) VALUES (216, 'user216');
INSERT INTO public."user" (id, name) VALUES (217, 'user217');
INSERT INTO public."user" (id, name) VALUES (218, 'user218');
INSERT INTO public."user" (id, name) VALUES (219, 'user219');
INSERT INTO public."user" (id, name) VALUES (220, 'user220');
INSERT INTO public."user" (id, name) VALUES (221, 'user221');
INSERT INTO public."user" (id, name) VALUES (222, 'user222');
INSERT INTO public."user" (id, name) VALUES (223, 'user223');
INSERT INTO public."user" (id, name) VALUES (224, 'user224');
INSERT INTO public."user" (id, name) VALUES (225, 'user225');
INSERT INTO public."user" (id, name) VALUES (226, 'user226');
INSERT INTO public."user" (id, name) VALUES (227, 'user227');
INSERT INTO public."user" (id, name) VALUES (228, 'user228');
INSERT INTO public."user" (id, name) VALUES (229, 'user229');
INSERT INTO public."user" (id, name) VALUES (230, 'user230');
INSERT INTO public."user" (id, name) VALUES (231, 'user231');
INSERT INTO public."user" (id, name) VALUES (232, 'user232');
INSERT INTO public."user" (id, name) VALUES (233, 'user233');
INSERT INTO public."user" (id, name) VALUES (234, 'user234');
INSERT INTO public."user" (id, name) VALUES (235, 'user235');
INSERT INTO public."user" (id, name) VALUES (236, 'user236');
INSERT INTO public."user" (id, name) VALUES (237, 'user237');
INSERT INTO public."user" (id, name) VALUES (238, 'user238');
INSERT INTO public."user" (id, name) VALUES (239, 'user239');
INSERT INTO public."user" (id, name) VALUES (240, 'user240');
INSERT INTO public."user" (id, name) VALUES (241, 'user241');
INSERT INTO public."user" (id, name) VALUES (242, 'user242');
INSERT INTO public."user" (id, name) VALUES (243, 'user243');
INSERT INTO public."user" (id, name) VALUES (244, 'user244');
INSERT INTO public."user" (id, name) VALUES (245, 'user245');
INSERT INTO public."user" (id, name) VALUES (246, 'user246');
INSERT INTO public."user" (id, name) VALUES (247, 'user247');
INSERT INTO public."user" (id, name) VALUES (248, 'user248');
INSERT INTO public."user" (id, name) VALUES (249, 'user249');
INSERT INTO public."user" (id, name) VALUES (250, 'user250');
INSERT INTO public."user" (id, name) VALUES (251, 'user251');
INSERT INTO public."user" (id, name) VALUES (252, 'user252');
INSERT INTO public."user" (id, name) VALUES (253, 'user253');
INSERT INTO public."user" (id, name) VALUES (254, 'user254');
INSERT INTO public."user" (id, name) VALUES (255, 'user255');
INSERT INTO public."user" (id, name) VALUES (256, 'user256');
INSERT INTO public."user" (id, name) VALUES (257, 'user257');
INSERT INTO public."user" (id, name) VALUES (258, 'user258');
INSERT INTO public."user" (id, name) VALUES (259, 'user259');
INSERT INTO public."user" (id, name) VALUES (260, 'user260');
INSERT INTO public."user" (id, name) VALUES (261, 'user261');
INSERT INTO public."user" (id, name) VALUES (262, 'user262');
INSERT INTO public."user" (id, name) VALUES (263, 'user263');
INSERT INTO public."user" (id, name) VALUES (264, 'user264');
INSERT INTO public."user" (id, name) VALUES (265, 'user265');
INSERT INTO public."user" (id, name) VALUES (266, 'user266');
INSERT INTO public."user" (id, name) VALUES (267, 'user267');
INSERT INTO public."user" (id, name) VALUES (268, 'user268');
INSERT INTO public."user" (id, name) VALUES (269, 'user269');
INSERT INTO public."user" (id, name) VALUES (270, 'user270');
INSERT INTO public."user" (id, name) VALUES (271, 'user271');
INSERT INTO public."user" (id, name) VALUES (272, 'user272');
INSERT INTO public."user" (id, name) VALUES (273, 'user273');
INSERT INTO public."user" (id, name) VALUES (274, 'user274');
INSERT INTO public."user" (id, name) VALUES (275, 'user275');
INSERT INTO public."user" (id, name) VALUES (276, 'user276');
INSERT INTO public."user" (id, name) VALUES (277, 'user277');
INSERT INTO public."user" (id, name) VALUES (278, 'user278');
INSERT INTO public."user" (id, name) VALUES (279, 'user279');
INSERT INTO public."user" (id, name) VALUES (280, 'user280');
INSERT INTO public."user" (id, name) VALUES (281, 'user281');
INSERT INTO public."user" (id, name) VALUES (282, 'user282');
INSERT INTO public."user" (id, name) VALUES (283, 'user283');
INSERT INTO public."user" (id, name) VALUES (284, 'user284');
INSERT INTO public."user" (id, name) VALUES (285, 'user285');
INSERT INTO public."user" (id, name) VALUES (286, 'user286');
INSERT INTO public."user" (id, name) VALUES (287, 'user287');
INSERT INTO public."user" (id, name) VALUES (288, 'user288');
INSERT INTO public."user" (id, name) VALUES (289, 'user289');
INSERT INTO public."user" (id, name) VALUES (290, 'user290');
INSERT INTO public."user" (id, name) VALUES (291, 'user291');
INSERT INTO public."user" (id, name) VALUES (292, 'user292');
INSERT INTO public."user" (id, name) VALUES (293, 'user293');
INSERT INTO public."user" (id, name) VALUES (294, 'user294');
INSERT INTO public."user" (id, name) VALUES (295, 'user295');
INSERT INTO public."user" (id, name) VALUES (296, 'user296');
INSERT INTO public."user" (id, name) VALUES (297, 'user297');
INSERT INTO public."user" (id, name) VALUES (298, 'user298');
INSERT INTO public."user" (id, name) VALUES (299, 'user299');
INSERT INTO public."user" (id, name) VALUES (300, 'user300');