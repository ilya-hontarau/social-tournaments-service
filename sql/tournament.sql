 CREATE TABLE tournaments (
     id INT NOT NULL AUTO_INCREMENT,
     name VARCHAR(20) NOT NULL,
     deposit INT(10) UNSIGNED NOT NULL,
     prize INT(10) UNSIGNED NOT NULL,
     finished BOOL DEFAULT false,
     winner INT,
     FOREIGN KEY(winner) REFERENCES users(id),
     PRIMARY KEY(id)
);
 
CREATE TABLE participants (
    user_id INT NOT NULL,
    tournament_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (tournament_id) REFERENCES tournaments(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id,tournament_id)
);
