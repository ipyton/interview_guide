CREATE TABLE HotPosts (
                          post_id bigint PRIMARY KEY,
                          resource_id bigint,
                          type text,
                          title text,
                          content text,
                          likes int,
                          extra1 text,
                          extra2 text,
                          extra3 text
) WITH default_time_to_live = 0;



CREATE TABLE questions (
                           question_id int PRIMARY KEY,
                           title text,
                           content text,
                           author_id int,
                           avatar text,
                           likes int,
                           views int,
                           tag1 text,
                           tag2 text,
                           tag3 text,
                           tag4 text,
                           tag5 text,
                           extra1 text,
                           extra2 text
) WITH default_time_to_live = 0;