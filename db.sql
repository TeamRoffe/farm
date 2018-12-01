-- Create syntax for TABLE 'categories'
CREATE TABLE `categories` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4;

-- Create syntax for TABLE 'drink_ingredients'
CREATE TABLE `drink_ingredients` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `drink_id` int(11) DEFAULT NULL,
  `liquid_id` int(11) DEFAULT NULL,
  `volume` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4;

-- Create syntax for TABLE 'drinks'
CREATE TABLE `drinks` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `drink_name` varchar(100) NOT NULL,
  `url` varchar(255) DEFAULT NULL,
  `category` int(11) NOT NULL,
  `description` text,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4;

-- Create syntax for TABLE 'liquids'
CREATE TABLE `liquids` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `liquid_name` varchar(100) DEFAULT NULL,
  `url` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4;