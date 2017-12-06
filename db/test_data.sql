-- Requires there to be a user with id 1

INSERT INTO "claims" ("created_by_id","truth","description","title") VALUES ('1','0.5','This is a debate about the legacy of the recently-deceased Fidel Castro. Should he be praised or demonized for what he did in his 50 years as dictator of Cuba?','Fidel Castro deserves to be praised for his legacy');

-- Root Pro-truth arguments

INSERT INTO "claims" ("created_by_id","truth","description","title") VALUES ('1','0.5','It has been widely claimed that Fidel Castro implemented policies in Cuba that dramatically improved education and healthcare for Cuban citizens. Is this true?','Fidel Castro promoted healthcare and education in Cuba for all its citizens');

INSERT INTO "arguments" ("created_by_id", "target_claim_id","claim_id","type","strength","description","title") SELECT 1, d1.id, d2.id, 1, 0.75, 'Fidel Castro deserves praise for the improvements he made to healthcare and education in Cuba.', 'He promoted healthcare and education in Cuba for all' FROM claims d1 LEFT JOIN claims d2 ON d2.title = 'Fidel Castro promoted healthcare and education in Cuba for all its citizens' WHERE d1.title = 'Fidel Castro deserves to be praised for his legacy';

-- Edu and Healthcare Pro-truth arguments

INSERT INTO "claims" ("created_by_id","truth","description","title") VALUES ('1','0.5','','The Fidel Castro regime produced tens of thousands of doctors and teachers and achieved some of the lowest infant mortality and illiteracy rates in the Western hemisphere.');

INSERT INTO "arguments" ("created_by_id", "target_claim_id","claim_id","type","strength","description","title") SELECT 1, d1.id, d2.id, 1, 0.75, '', 'He produced tens of thousands of doctors and teachers and achieved some of the lowest infant mortality and illiteracy rates in the Western hemisphere.' FROM claims d1 LEFT JOIN claims d2 ON d2.title = 'The Fidel Castro regime produced tens of thousands of doctors and teachers and achieved some of the lowest infant mortality and illiteracy rates in the Western hemisphere.' WHERE d1.title = 'Fidel Castro promoted healthcare and education in Cuba for all its citizens';

INSERT INTO "links" ("created_by_id", "description", "title", "url", "claim_id") SELECT 1, 'An article which discusses his achievements in the area of healthcare', 'Cuba Leader Fidel Castro Dead at 90', 'http://www.aljazeera.com/news/americas/2016/11/cuba-leader-fidel-castro-dead-90-161126053354637.html', d.id FROM claims d WHERE title = 'The Fidel Castro regime produced tens of thousands of doctors and teachers and achieved some of the lowest infant mortality and illiteracy rates in the Western hemisphere.';

INSERT INTO "claims" ("created_by_id","truth","description","title") VALUES ('1','0.5','In 1961, two years after Castro''s revolution won power, the new Cuban government launched an ambitious campaign to stamp out illiteracy. Some 250,000 volunteer teachers, many of them young women, fanned out across the country, especially in rural areas where access to education was spotty and the need was greatest. In the space of a year, about 700,000 people learned to read and write, said "Maestra," a documentary that explores the initiative''s history. Today, Cuba reports a literacy rate of 99.8 percent, on par with the most developed nations in the world.','Fidel Castro''s policy was responsible for raising the literacy rate in Cuba to 99.8 percent.');

INSERT INTO "arguments" ("created_by_id", "target_claim_id","claim_id","type","strength","description","title") SELECT 1, d1.id, d2.id, 1, 0.75, '', 'He raised the literacy rate in Cuba to 99.8%' FROM claims d1 LEFT JOIN claims d2 ON d2.title = 'Fidel Castro''s policy was responsible for raising the literacy rate in Cuba to 99.8 percent.' WHERE d1.title = 'Fidel Castro promoted healthcare and education in Cuba for all its citizens';

INSERT INTO "links" ("created_by_id", "description", "title", "url", "claim_id") SELECT 1, 'An article from Fox News discussing the many impacts Fidel Castro had on life in Cuba', 'From Milk to Lightbulbs Fidel Castro Reshaped Life in Cuba', 'http://www.foxnews.com/world/2016/11/28/from-milk-to-lightbulbs-fidel-castro-reshaped-life-in-cuba.html', d.id FROM claims d WHERE title = 'The Fidel Castro regime produced tens of thousands of doctors and teachers and achieved some of the lowest infant mortality and illiteracy rates in the Western hemisphere.';

-- Back to Root Pro-truth

INSERT INTO "claims" ("created_by_id","truth","description","title") VALUES ('1','0.5','','Fidel Castro was instrumental in ending Apartheid in South Africa');

INSERT INTO "arguments" ("created_by_id", "target_claim_id","claim_id","type","strength","description","title") SELECT 1, d1.id, d2.id, 1, 0.75, '', 'He was instrumental in ending Apartheid in South Africa' FROM claims d1 LEFT JOIN claims d2 ON d2.title = 'Fidel Castro was instrumental in ending Apartheid in South Africa' WHERE d1.title = 'Fidel Castro deserves to be praised for his legacy';

INSERT INTO "claims" ("created_by_id","truth","description","title") VALUES ('1','0.5','','During the reign of Fidel Castro, Cuba was instrumental in the liberation of many African nations');

INSERT INTO "arguments" ("created_by_id", "target_claim_id","claim_id","type","strength","description","title") SELECT 1, d1.id, d2.id, 1, 0.75, '', 'During his reign, Cuba was instrumental in the liberation of many African nations' FROM claims d1 LEFT JOIN claims d2 ON d2.title = 'During the reign of Fidel Castro, Cuba was instrumental in the liberation of many African nations' WHERE d1.title = 'Fidel Castro deserves to be praised for his legacy';

-- Liberation of African Nations Pro-truth arguments

INSERT INTO "claims" ("created_by_id","truth","description","title") VALUES ('1','0.5','','Under instruction from Fidel Castro, Cuba played an important role in support of liberation struggles in Angola and Mozambique');

INSERT INTO "arguments" ("created_by_id", "target_claim_id","claim_id","type","strength","description","title") SELECT 1, d1.id, d2.id, 1, 0.75, '', 'Cuba was instrumental in liberating Angola and Mozambique' FROM claims d1 LEFT JOIN claims d2 ON d2.title = 'Under instruction from Fidel Castro, Cuba played an important role in support of liberation struggles in Angola and Mozambique' WHERE d1.title = 'During the reign of Fidel Castro, Cuba was instrumental in the liberation of many African nations';

INSERT INTO "links" ("created_by_id", "description", "title", "url", "claim_id") SELECT 1, '', 'A Look at Fidel Castro''s Legacy from a Fair Perspective', 'https://www.theguardian.com/world/2016/nov/29/look-at-fidel-castro-legacy-from-a-fair-perspective', d.id FROM claims d WHERE title = 'Under instruction from Fidel Castro, Cuba played an important role in support of liberation struggles in Angola and Mozambique';

-- Liberation of African Nations Con-impact arguments

INSERT INTO "claims" ("created_by_id","truth","description","title") VALUES ('1','0.5','I''ve not seen comment on Cuba’s support for Mengistu Haile Mariam when Ethiopia was invaded by the Somali forces of Siad Barre. What is the judgment? Mengistu eventually ousted (to sanctuary in Zimbabwe), Siad forced out and dead; Somalia in a very sorry state and Ethiopia again in some turmoil after a period of relative contentment.','During the reign of Fidel Castro, Cuba''s attempts at helping Ethiopia ended in disaster');

INSERT INTO "arguments" ("created_by_id", "target_argument_id","claim_id","type","strength","description","title") SELECT 1, a1.id, d2.id, 6, 0.75, '', 'Cuba''s attempts at helping Ethiopia ended in disaster' FROM arguments a1 LEFT JOIN claims d2 ON d2.title = 'During the reign of Fidel Castro, Cuba''s attempts at helping Ethiopia ended in disaster' WHERE a1.title = 'During his reign, Cuba was instrumental in the liberation of many African nations';

-- Back to Root Pro-truth

INSERT INTO "claims" ("created_by_id","truth","description","title") VALUES ('1','0.5','','Fidel Castro was loved by the Cuban people');

INSERT INTO "arguments" ("created_by_id", "target_claim_id","claim_id","type","strength","description","title") SELECT 1, d1.id, d2.id, 1, 0.75, '', 'He was loved by the Cuban people' FROM claims d1 LEFT JOIN claims d2 ON d2.title = 'Fidel Castro was loved by the Cuban people' WHERE d1.title = 'Fidel Castro deserves to be praised for his legacy';

-- Loved by Cubans Pro-truth

INSERT INTO "claims" ("created_by_id","truth","description","title") VALUES ('1','0.5','Castro walked the streets of Havana safely despite over 500 attempts on his life by the CIA. He is revered by his people – even the BBC in a report on Sunday acknowledged that the young people in Cuba think favourably of him.','Fidel Castro felt safe to walk the streets in Cuba even after 500 attempts of assassination by the CIA');

INSERT INTO "arguments" ("created_by_id", "target_claim_id","claim_id","type","strength","description","title") SELECT 1, d1.id, d2.id, 1, 0.75, '', 'He felt safe to to walk the streets in Cuba even after multiple attempts on his life by the CIA' FROM claims d1 LEFT JOIN claims d2 ON d2.title = 'Fidel Castro felt safe to walk the streets in Cuba even after 500 attempts of assassination by the CIA' WHERE d1.title = 'Fidel Castro was loved by the Cuban people';

INSERT INTO "links" ("created_by_id", "description", "title", "url", "claim_id") SELECT 1, '', 'A Look at Fidel Castro''s Legacy from a Fair Perspective', 'https://www.theguardian.com/world/2016/nov/29/look-at-fidel-castro-legacy-from-a-fair-perspective', d.id FROM claims d WHERE title = 'Fidel Castro felt safe to walk the streets in Cuba even after 500 attempts of assassination by the CIA';

INSERT INTO "claims" ("created_by_id","truth","description","title") VALUES ('1','0.5','His "struggle for good and equality for all" makes him "one of the great ones. Fidel is like a myth. He''s like Che," said Del Toro, referring to the late Argentinean revolutionary Ernesto "Che" Guevara.','Fidel Castro is as revered by the Cuban people as Che Guevara');

INSERT INTO "arguments" ("created_by_id", "target_claim_id","claim_id","type","strength","description","title") SELECT 1, d1.id, d2.id, 1, 0.75, '', 'He felt safe to to walk the streets in Cuba even after multiple attempts on his life by the CIA' FROM claims d1 LEFT JOIN claims d2 ON d2.title = 'Fidel Castro is as revered by the Cuban people as Che Guevara' WHERE d1.title = 'Fidel Castro was loved by the Cuban people';

INSERT INTO "links" ("created_by_id", "description", "title", "url", "claim_id") SELECT 1, '', 'Cuba Leader Fidel Castro Dead at 90', 'http://www.aljazeera.com/news/americas/2016/11/cuba-leader-fidel-castro-dead-90-161126053354637.html', d.id FROM claims d WHERE title = 'Fidel Castro is as revered by the Cuban people as Che Guevara';

-- Back to Root Pro-truth

INSERT INTO "claims" ("created_by_id","truth","description","title") VALUES ('1','0.5','','Fidel Castro created a more egalitarian society in Cuba');

INSERT INTO "arguments" ("created_by_id", "target_claim_id","claim_id","type","strength","description","title") SELECT 1, d1.id, d2.id, 1, 0.75, '', 'He created a more egalitarian society in Cuba' FROM claims d1 LEFT JOIN claims d2 ON d2.title = 'Fidel Castro created a more egalitarian society in Cuba' WHERE d1.title = 'Fidel Castro deserves to be praised for his legacy';

-- Egalitarian Pro-truth arguments

INSERT INTO "claims" ("created_by_id","truth","description","title") VALUES ('1','0.5','','The Batista regime in Cuba served only Batist''s own interests and those of the wealthy');

INSERT INTO "arguments" ("created_by_id", "target_claim_id","claim_id","type","strength","description","title") SELECT 1, d1.id, d2.id, 1, 0.75, '', 'His regime replaced the Batista regime, a dictator that served his own interests and those of the wealthy' FROM claims d1 LEFT JOIN claims d2 ON d2.title = 'The Batista regime in Cuba served only Batist''s own interests and those of the wealthy' WHERE d1.title = 'Fidel Castro created a more egalitarian society in Cuba';

INSERT INTO "claims" ("created_by_id","truth","description","title") VALUES ('1','0.5','','Fidel Castro implemented agrarian reform in Cuba that resulted in more economic equality');

INSERT INTO "arguments" ("created_by_id", "target_claim_id","claim_id","type","strength","description","title") SELECT 1, d1.id, d2.id, 1, 0.75, '', 'He implemented agrarian reform in Cuba that resulted in more economic equality' FROM claims d1 LEFT JOIN claims d2 ON d2.title = 'Fidel Castro implemented agrarian reform in Cuba that resulted in more economic equality' WHERE d1.title = 'Fidel Castro created a more egalitarian society in Cuba';

INSERT INTO "claims" ("created_by_id","truth","description","title") VALUES ('1','0.5','','Fidel Castro ended homelessness in Cuba');

INSERT INTO "arguments" ("created_by_id", "target_claim_id","claim_id","type","strength","description","title") SELECT 1, d1.id, d2.id, 1, 0.75, '', 'He ended homelessness in Cuba' FROM claims d1 LEFT JOIN claims d2 ON d2.title = 'Fidel Castro ended homelessness in Cuba' WHERE d1.title = 'Fidel Castro created a more egalitarian society in Cuba';

-- Ended homelessness Pro-truth arguments

INSERT INTO "claims" ("created_by_id","truth","description","title") VALUES ('1','0.5','','you don''t see anyone living on the street. Everyone has a house');

INSERT INTO "arguments" ("created_by_id", "target_claim_id","claim_id","type","strength","description","title") SELECT 1, d1.id, d2.id, 1, 0.75, '', 'you don''t see anyone living on the street. Everyone has a house' FROM claims d1 LEFT JOIN claims d2 ON d2.title = 'you don''t see anyone living on the street. Everyone has a house' WHERE d1.title = 'Fidel Castro ended homelessness in Cuba';

INSERT INTO "links" ("created_by_id", "description", "title", "url", "claim_id") SELECT 1, '', 'Cuba Leader Fidel Castro Dead at 90', 'http://www.aljazeera.com/news/americas/2016/11/cuba-leader-fidel-castro-dead-90-161126053354637.html', d.id FROM claims d WHERE title = 'you don''t see anyone living on the street. Everyone has a house';

-- Back to Egalitarian Pro-truth arguments

INSERT INTO "claims" ("created_by_id","truth","description","title") VALUES ('1','0.5','','In spite of hardships, the population is content due to the egalitarian nature of society in Cuba');

INSERT INTO "arguments" ("created_by_id", "target_claim_id","claim_id","type","strength","description","title") SELECT 1, d1.id, d2.id, 1, 0.75, '', 'In spite of hardships, the population is content due to the egalitarian nature of society' FROM claims d1 LEFT JOIN claims d2 ON d2.title = 'In spite of hardships, the population is content due to the egalitarian nature of society in Cuba' WHERE d1.title = 'Fidel Castro created a more egalitarian society in Cuba';

-- Contented with society Pro-truth arguments

INSERT INTO "claims" ("created_by_id","truth","description","title") VALUES ('1','0.5','"Our visit this month confirmed our belief that in Cuba, despite the hardships that the people have to endure, the majority of the population are cheerful, friendly and content with the egalitarian system that Castro introduced."','The Guardian reports that the population is generally content with the egalitarian nature of Cuba');

INSERT INTO "arguments" ("created_by_id", "target_claim_id","claim_id","type","strength","description","title") SELECT 1, d1.id, d2.id, 1, 0.75, '', 'According to The Guardian, the population of Cuba is content' FROM claims d1 LEFT JOIN claims d2 ON d2.title = 'The Guardian reports that the population is generally content with the egalitarian nature of Cuba' WHERE d1.title = 'In spite of hardships, the population is content due to the egalitarian nature of society in Cuba';

INSERT INTO "links" ("created_by_id", "description", "title", "url", "claim_id") SELECT 1, '', 'Cuba Leader Fidel Castro Dead at 90', 'http://www.aljazeera.com/news/americas/2016/11/cuba-leader-fidel-castro-dead-90-161126053354637.html', d.id FROM claims d WHERE title = 'The Guardian reports that the population is generally content with the egalitarian nature of Cuba';

-- Egalitarian Con-truth arguments

INSERT INTO "claims" ("created_by_id","truth","description","title") VALUES ('1','0.5','"It became clear very quickly after the revolution that while all Cubans were equal, some were more equal than others. And that is why, while ordinary Cubans scrape for crumbs, the privileged communists in Castro’s nomenklatura live like millionaires."','Even in Fidel Castro''s Cuba, not all people are treated equally');

INSERT INTO "arguments" ("created_by_id", "target_claim_id","claim_id","type","strength","description","title") SELECT 1, d1.id, d2.id, 2, 0.75, '', 'Even in Cuba, the government plays with favoritism' FROM claims d1 LEFT JOIN claims d2 ON d2.title = 'Even in Fidel Castro''s Cuba, not all people are treated equally' WHERE d1.title = 'Fidel Castro created a more egalitarian society in Cuba';

INSERT INTO "links" ("created_by_id", "description", "title", "url", "claim_id") SELECT 1, '', 'Monster Fidel Castro Leaves Blood of Innocents in his Wake', 'http://www.breitbart.com/national-security/2016/11/27/glazov-monster-fidel-castro-leaves-blood-innocents-wake/', d.id FROM claims d WHERE title = 'Even in Fidel Castro''s Cuba, not all people are treated equally';

-- Back to Root Pro-truth

INSERT INTO "claims" ("created_by_id","truth","description","title") VALUES ('1','0.5','','Fidel Castro serves as inspiration to all nations that feel a need to stand up to bullying nations like the United States');

INSERT INTO "arguments" ("created_by_id", "target_claim_id","claim_id","type","strength","description","title") SELECT 1, d1.id, d2.id, 1, 0.75, '', 'He serves as inspiration to all nations that feel a need to stand up to bullying nations like the United States' FROM claims d1 LEFT JOIN claims d2 ON d2.title = 'Fidel Castro serves as inspiration to all nations that feel a need to stand up to bullying nations like the United States' WHERE d1.title = 'Fidel Castro deserves to be praised for his legacy';

-- Inspiration Pro-truth arguments

INSERT INTO "claims" ("created_by_id","truth","description","title") VALUES ('1','0.5','He had the ability to face - together with his people - the most powerful enemy on earth. And I think [his] ideas will live on forever.','Fidel Castro had the ability to face - together with his people - the most powerful enemy on earth.');

INSERT INTO "arguments" ("created_by_id", "target_claim_id","claim_id","type","strength","description","title") SELECT 1, d1.id, d2.id, 1, 0.75, '', 'He had the ability to face - together with his people - the most powerful enemy on earth.' FROM claims d1 LEFT JOIN claims d2 ON d2.title = 'Fidel Castro had the ability to face - together with his people - the most powerful enemy on earth.' WHERE d1.title = 'Fidel Castro serves as inspiration to all nations that feel a need to stand up to bullying nations like the United States';

INSERT INTO "links" ("created_by_id", "description", "title", "url", "claim_id") SELECT 1, '', 'Cuba Leader Fidel Castro Dead at 90', 'http://www.aljazeera.com/news/americas/2016/11/cuba-leader-fidel-castro-dead-90-161126053354637.html', d.id FROM claims d WHERE title = 'Fidel Castro had the ability to face - together with his people - the most powerful enemy on earth.';

-- Back to Root Pro-truth

INSERT INTO "claims" ("created_by_id","truth","description","title") VALUES ('1','0.5','','Fidel Castro sent Cuban doctors to third world countries to help improve health conditions');

INSERT INTO "arguments" ("created_by_id", "target_claim_id","claim_id","type","strength","description","title") SELECT 1, d1.id, d2.id, 1, 0.75, '', 'He sent Cuban doctors to third world countries to help improve health conditions' FROM claims d1 LEFT JOIN claims d2 ON d2.title = 'Fidel Castro sent Cuban doctors to third world countries to help improve health conditions' WHERE d1.title = 'Fidel Castro deserves to be praised for his legacy';

-- Inspiration Pro-truth arguments

INSERT INTO "claims" ("created_by_id","truth","description","title") VALUES ('1','0.5','"...sending Cuba''s doctors to heal the Third World''s sick"','A Fox News article mentions that Fidel Castro sent Cuban doctors to other countries to help with their health system');

INSERT INTO "arguments" ("created_by_id", "target_claim_id","claim_id","type","strength","description","title") SELECT 1, d1.id, d2.id, 1, 0.75, '', 'He did, according to at least on Fox News article' FROM claims d1 LEFT JOIN claims d2 ON d2.title = 'A Fox News article mentions that Fidel Castro sent Cuban doctors to other countries to help with their health system' WHERE d1.title = 'Fidel Castro sent Cuban doctors to third world countries to help improve health conditions';

INSERT INTO "links" ("created_by_id", "description", "title", "url", "claim_id") SELECT 1, 'An article from Fox News discussing the many impacts Fidel Castro had on life in Cuba', 'From Milk to Lightbulbs Fidel Castro Reshaped Life in Cuba', 'http://www.foxnews.com/world/2016/11/28/from-milk-to-lightbulbs-fidel-castro-reshaped-life-in-cuba.html', d.id FROM claims d WHERE title = 'A Fox News article mentions that Fidel Castro sent Cuban doctors to other countries to help with their health system';

-- Root Con-truth Arguments

INSERT INTO "claims" ("created_by_id","truth","description","title") VALUES ('1','0.5','','Fidel Castro executed thousands of people for being political dissidents');

INSERT INTO "arguments" ("created_by_id", "target_claim_id","claim_id","type","strength","description","title") SELECT 1, d1.id, d2.id, 2, 0.75, '', 'He executed thousands of people for being political dissidents' FROM claims d1 LEFT JOIN claims d2 ON d2.title = 'Fidel Castro executed thousands of people for being political dissidents' WHERE d1.title = 'Fidel Castro deserves to be praised for his legacy';

-- Execution Pro-truth Arguments

INSERT INTO "claims" ("created_by_id","truth","description","title") VALUES ('1','0.5','','The Victims of Communism Memorial Foundation estimates 73,000 were killed since Fidel Castro came to power in 1959, until the time of his death');

INSERT INTO "arguments" ("created_by_id", "target_claim_id","claim_id","type","strength","description","title") SELECT 1, d1.id, d2.id, 1, 0.75, '', 'The Victims of Communism Memorial Foundation estimates 73,000 were killed since Castro came to power in 1959' FROM claims d1 LEFT JOIN claims d2 ON d2.title = 'The Victims of Communism Memorial Foundation estimates 73,000 were killed since Fidel Castro came to power in 1959, until the time of his death' WHERE d1.title = 'Fidel Castro executed thousands of people for being political dissidents';

INSERT INTO "links" ("created_by_id", "description", "title", "url", "claim_id") SELECT 1, '', 'Fidel Castro''s Legacy of Murder and Repression Whitewashed by the Left', 'http://www.theaustralian.com.au/opinion/columnists/janet-albrechtsen/fidel-castros-legacy-of-murder-and-repression-whitewashed-by-the-left/news-story/8e12657fc5a8fa70fdae9e5ba6f5daff', d.id FROM claims d WHERE title = 'The Victims of Communism Memorial Foundation estimates 73,000 were killed since Fidel Castro came to power in 1959, until the time of his death';

-- Execution Con-strength Arguments

INSERT INTO "claims" ("created_by_id","truth","description","title") VALUES ('1','0.5','While Castro may be rightly criticised for executing Batista supporters, even those guilty of torture and multiple murder, it may be salutary to remember that back then, in 1959, Britain executed people accused of a single murder. It was also a time when British forces were imprisoning and torturing Kenyans, and those of the French multiparty democracy were torturing and killing Algerians. Even those crimes pale before the horrors the US multiparty democracy was shortly to unleash on Vietnam.','The executions following Fidel Castro''s rise to power were necessary, and were minor in comparison to other atrocities committed by nations at that time.');

INSERT INTO "arguments" ("created_by_id", "target_argument_id","claim_id","type","strength","description","title") SELECT 1, a1.id, d2.id, 4, 0.75, '', 'The executions were minor in comparison to what was happening in the world at the time' FROM arguments a1 LEFT JOIN claims d2 ON d2.title = 'The executions following Fidel Castro''s rise to power were necessary, and were minor in comparison to other atrocities committed by nations at that time.' WHERE a1.title = 'He executed thousands of people for being political dissidents';

INSERT INTO "links" ("created_by_id", "description", "title", "url", "claim_id") SELECT 1, '', 'A Look at Fidel Castro''s Legacy from a Fair Perspective', 'https://www.theguardian.com/world/2016/nov/29/look-at-fidel-castro-legacy-from-a-fair-perspective', d.id FROM claims d WHERE title = 'The executions following Fidel Castro''s rise to power were necessary, and were minor in comparison to other atrocities committed by nations at that time.';

