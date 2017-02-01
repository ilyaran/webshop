CREATE TABLE items (
	  id bigserial ,
	  Title varchar (255) NOT NULL DEFAULT '',
	  Created  date  NOT NULL DEFAULT now(),
	  Seller  varchar (255) NOT NULL,
	  Price   numeric(20,2),
	  Image varchar (255) NOT NULL,
	  Description text
);
