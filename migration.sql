CREATE TABLE IF NOT EXISTS ranking (
  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  ranked_datetime DATETIME NOT NULL,
  rank INT NOT NULL,
  video_id VARCHAR(20) NOT NULL,
  title VARCHAR(100) NOT NULL,
  upload_date DATETIME NOT NULL,
  thumb_url TEXT NOT NULL,
  length VARCHAR(10) NOT NULL,
  view INT NOT NULL,
  comment INT NOT NULL,
  mylist INT NOT NULL,
  tags TEXT
)