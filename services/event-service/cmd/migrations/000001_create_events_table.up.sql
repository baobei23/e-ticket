CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    location VARCHAR(255),
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    total_seats INTEGER NOT NULL,
    available_seats INTEGER NOT NULL,
    price DECIMAL(20, 2) NOT NULL
);

INSERT INTO events (
  name, description, location,
  start_time, end_time,
  total_seats, available_seats, price
) VALUES
('Konser Coldplay','Music of the Spheres','GBK',NOW()+INTERVAL '1 month',NOW()+INTERVAL '1 month 4 hours',50000,50000,3500000),
('Java Jazz Festival','International Jazz Music Festival','JIExpo Kemayoran',NOW()+INTERVAL '2 months',NOW()+INTERVAL '2 months 8 hours',30000,30000,1500000),
('We The Fest','Music, Arts & Fashion Festival','GBK Sports Complex',NOW()+INTERVAL '3 months',NOW()+INTERVAL '3 months 10 hours',40000,40000,1800000),
('DWP Jakarta','Electronic Dance Music Festival','JIExpo Kemayoran',NOW()+INTERVAL '4 months',NOW()+INTERVAL '4 months 12 hours',50000,50000,2200000),
('Soundrenaline','Largest Music Festival in Indonesia','Bali',NOW()+INTERVAL '5 months',NOW()+INTERVAL '5 months 9 hours',35000,35000,1700000),
('Prambanan Jazz','Jazz Festival with Cultural Heritage','Candi Prambanan',NOW()+INTERVAL '6 months',NOW()+INTERVAL '6 months 6 hours',20000,20000,1400000),
('Synchronize Fest','Indonesian Music Festival','Gambir Expo',NOW()+INTERVAL '7 months',NOW()+INTERVAL '7 months 8 hours',30000,30000,1200000),
('Head in the Clouds','Asian Music & Culture Festival','Jakarta',NOW()+INTERVAL '8 months',NOW()+INTERVAL '8 months 7 hours',25000,25000,2000000),
('Jakarta Fair Concert','Annual Jakarta Fair Music Event','PRJ Kemayoran',NOW()+INTERVAL '9 months',NOW()+INTERVAL '9 months 5 hours',40000,40000,900000),
('Rock in Solo','Rock Music Festival','Solo',NOW()+INTERVAL '10 months',NOW()+INTERVAL '10 months 6 hours',15000,15000,750000),
('Indie Fest Bandung','Independent Music Festival','Bandung',NOW()+INTERVAL '11 months',NOW()+INTERVAL '11 months 6 hours',12000,12000,650000),
('Pop Nation','Pop Music Concert','ICE BSD',NOW()+INTERVAL '12 months',NOW()+INTERVAL '12 months 5 hours',20000,20000,1300000),
('Hip Hop Nation','Hip Hop Music Showcase','Jakarta Convention Center',NOW()+INTERVAL '13 months',NOW()+INTERVAL '13 months 4 hours',18000,18000,1100000),
('K-Pop Wave','Korean Pop Music Concert','ICE BSD',NOW()+INTERVAL '14 months',NOW()+INTERVAL '14 months 4 hours',25000,25000,2500000),
('Metal Storm','Metal Music Festival','Bandung',NOW()+INTERVAL '15 months',NOW()+INTERVAL '15 months 7 hours',10000,10000,850000),
('Acoustic Night','Acoustic & Indie Live Session','Jakarta',NOW()+INTERVAL '16 months',NOW()+INTERVAL '16 months 3 hours',5000,5000,450000),
('Orchestra Live','Symphony Orchestra Performance','Teater Jakarta',NOW()+INTERVAL '17 months',NOW()+INTERVAL '17 months 2 hours',3000,3000,900000),
('EDM Beach Party','Beachside EDM Festival','Bali',NOW()+INTERVAL '18 months',NOW()+INTERVAL '18 months 10 hours',20000,20000,1900000),
('City Pop Night','Japanese City Pop Live','Jakarta',NOW()+INTERVAL '19 months',NOW()+INTERVAL '19 months 4 hours',8000,8000,700000),
('Indonesian Legends','Legendary Indonesian Musicians','GBK',NOW()+INTERVAL '20 months',NOW()+INTERVAL '20 months 6 hours',45000,45000,1600000);
