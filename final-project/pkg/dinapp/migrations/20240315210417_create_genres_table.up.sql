create table if not exists genres(
    genre_id serial primary key, 
    genre_title varchar(255) unique,
);

create table if not exists movies
(   movie_id serial primary key, 
    movie_title varchar(255),
    description varchar(1000),
    year_of_production int,
    genre_id int references genres(genre_id),
);

-- insert into genres (genre_id, genre_title)
-- values
--     (1, 'horor'),
--     (2, 'rom-com'),
--     (3, 'detective'),
--     (4, 'drama');


-- select * from genres;


-- insert into movies (movie_id, movie_title, description, year_of_production, genre_id)
-- values
--     (1, 'the conjuring', 'in the early 1970s, the Perron family - Roger, Carolyn, and their five daughters - move into a new home in the Rhode Island countryside. Before long, they start encountering strange noises and smells, stopped clocks, slamming doors, and figures lurking in dark corners.',                                                     2013, 1),
--     (2, 'it', 'seven young outcasts in Derry, Maine, are about to face their worst nightmare - an ancient, shape-shifting evil that emerges from the sewer every 27 years to prey on the town''s children.',                                                                                                                                                  2017, 1),
--     (3, 'as above, so below', 'when a team of explorers venture into the catacombs that lie beneath the streets of Paris, they uncover the dark secret that lies within this city of the dead.',                                                                                                                                                              2014, 1),
--     (4, 'hereditary', 'when the matriarch of the Graham family passes away, her daughter and grandchildren begin to unravel cryptic and increasingly terrifying secrets about their ancestry, trying to outrun the sinister fate they have inherited.',                                                                                                       2018, 1),
--     (5, 'annabelle', 'a couple begins to experience terrifying supernatural occurrences involving a vintage doll shortly after their home is invaded by satanic cultists.',                                                                                                                                                                                   2014, 1),
--     (6, 'split', 'the film follows a man with dissociative identity disorder who kidnaps and imprisons three teenage girls in an isolated underground facility.',                                                                                                                                                                                            2017, 1),

--     (7, 'legally blonde', 'the story follows Elle Woods (Witherspoon), a sorority girl who attempts to win back her ex-boyfriend Warner Huntington III (Davis) by getting a Juris Doctor degree at Harvard Law School, and in the process, overcomes stereotypes against blondes and triumphs as a successful lawyer.',                                       2001, 2),
--     (8, 'bride wars', 'two childhood best friends, who have made many plans together for their respective weddings, turn into sworn enemies in a race to get married first.',                                                                                                                                                                                 2009, 2),
--     (9, 'how to lose a guy in 10 days', 'a ladies man bet his friends that he can get a woman to fall in love with him in 10 days, but unbeknownst to him, the woman he''s dating is actually a magazine columnist working on a new column called How to Lose a Guy in 10 Days, and she''s doing everything she could to drive the guy she''s dating crazy.', 2003, 2),
--     (10, 'clueless', 'the plot centers on a beautiful, popular, and rich high school student who befriends a new student and decides to give her a makeover while playing a matchmaker for her teachers and examining her own existence.',                                                                                                                    1995, 2),
--     (11, 'easy A', 'when Olive lies to her best friend about losing her virginity to one of the college boys, a girl overhears their conversation. Soon, her story spreads across the entire school like wildfire.',                                                                                                                                         2010, 2),
--     (12, 'just my luck', 'it tells the story of Ashley Albright who works in public relations and is the luckiest person in Manhattan, while Jake Hardin is a janitor and would-be music producer who seems to have terrible luck until their good and bad luck is switched upon kissing each other at a masquerade ball which changes both their lives',     2006, 2),

--     (13, 'the silence of the lambs', 'a young F.B.I. cadet must receive the help of an incarcerated and manipulative cannibal killer to help catch another serial killer, a madman who skins his victims.',                                                                                                                                                   1991, 3),
--     (14, 'seven', 'set in an unnamed, crime-ridden city, Seven''s narrative follows disenchanted, nearly retired detective William Somerset (Freeman) and his newly transferred partner David Mills (Pitt) as they try to stop a serial killer from executing a series of murders based on the seven deadly sins.',                                           1995, 3),
--     (15, 'happy death day', 'movie follows college student Tree Gelbman, who is murdered on the night of her birthday but begins reliving the day repeatedly, at which point she sets out to find the killer and stop her death.',                                                                                                                            2017, 3),
--     (16, 'prisoners', 'the film follows the abduction of two young girls in Pennsylvania and the subsequent search for the perpetrator by the police. After police arrest a young suspect and release him, the father of one of the daughters takes matters into his own hands.',                                                                             2013, 3),
--     (17, 'knives out', 'daniel Craig leads an eleven-actor ensemble cast as Benoit Blanc, famed private detective summoned to investigate the death of bestselling author Harlan Thrombey.',                                                                                                                                                                  2019, 3),

--     (18, 'the lovely bones', 'the Lovely Bones is a 2002 novel by American writer Alice Sebold. It is the story of a teenage girl who, after being raped and murdered, watches from her personal Heaven as her family and friends struggle to move on with their lives while she comes to terms with her own death.',                                         2009, 4),
--     (19, 'the curious case of benjamin button', 'Born under unusual circumstances, Benjamin Button (Brad Pitt) springs into being as an elderly man in a New Orleans nursing home and ages in reverse.',                                                                                                                                                      2008, 4),
--     (20, 'seven pounds', 'a life-shattering secret torments Ben Thomas. In order to find redemption, he sets out to change the lives of seven strangers. Over the course of his journey, he meets and falls in love with a cardiac patient named Emily, and in so doing, complicates his mission.',                                                           2008, 4),
--     (21, 'the boy in the striped pyjamas', 'the Boy in the Striped Pajamas by John Boyne is a poignant and heart-wrenching story about a young boy named Bruno who befriends a boy on the other side of a concentration camp fence during World War II, unaware of the true nature of the camp.',                                                             2008, 4),
--     (22, 'the devil wears prada', 'the Devil Wears Prada is a 2003 novel by Lauren Weisberger about a young woman who is hired as a personal assistant to a powerful fashion magazine editor, a job that becomes nightmarish as she struggles to keep up with her boss''s grueling schedule and demeaning demands.',                                          2006, 4);


-- select * from movies;


-- update movies
-- set description = 'when Olive lies to her best friend about losing her virginity to one of the college boys, a girl overhears their conversation. Soon, her story spreads across the entire school like wildfire.'
-- where movie_id = 11;
