CREATE TABLE collections (
                             collection_id int PRIMARY KEY,
                             resource_id int,
                             type text,
                             resource_title text,
                             content text
) WITH default_time_to_live = 0;



CREATE TABLE collection_items (
                                  user_id text,
                                  collection_id int,
                                  collection_name text,
                                  count int,
                                  extra1 text,
                                  extra2 text,
                                  extra3 text,
                                  PRIMARY KEY (user_id, collection_id)
) WITH default_time_to_live = 0;
